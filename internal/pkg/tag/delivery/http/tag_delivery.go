package http

import (
	"fmt"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/common/http/auth"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/notes"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
)

type Handler struct {
	tagUsecase  tag.Usecase
	noteUsecase notes.Usecase
	logger      logger.Logger
}

func NewHandler(tu tag.Usecase, nu notes.Usecase, l logger.Logger) *Handler {
	return &Handler{
		tagUsecase:  tu,
		noteUsecase: nu,
		logger:      l,
	}
}

// CreateAndLinkTag
// @Summary		Create and link tag to note
// @Tags		Tags
// @Description	Create a new tag and link it to a note, or link existing tag if it exists
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		CreateAndLinkTagRequest		true	"Tag info"
// @Success		201			{object}	models.Tag		"Tag created and linked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/create [post]
func (h *Handler) CreateAndLinkTag(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Name   string `json:"name" binding:"required"`
		NoteID string `json:"note_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Invalid note ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
		return
	}

	// Check if note exists
	_, err = h.noteUsecase.GetByID(noteID)
	if err != nil {
		h.logger.Errorf("Note not found: %w", err)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Note with ID '%s' not found", noteID)})
		return
	}

	// Check if user has access to the note
	access, err := h.noteUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		h.logger.Errorf("Failed to get user access: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if access == models.EmptyNoteAccess {
		h.logger.Infof("User %s does not have access to note %s", userID, noteID)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to note with ID '%s'", noteID)})
		return
	}

	tag, err := h.tagUsecase.CreateAndLinkTag(req.Name, noteID)
	if err != nil {
		h.logger.Errorf("Failed to create and link tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNameExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Tag with name '%s' already exists for this note", e.Name)})
		case *errors.NoteNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Note with ID '%s' not found", e.ID)})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// UnlinkTagFromNote
// @Summary		Unlink tag from note
// @Tags		Tags
// @Description	Remove the link between a tag and a note, delete tag if it has no more links
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		UnlinkTagRequest		true	"Tag and note IDs"
// @Success		200								"Tag unlinked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/unlink [post]
func (h *Handler) UnlinkTagFromNote(c *gin.Context) {
	var req UnlinkTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Errorf("Invalid request data: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Invalid note ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
		return
	}

	if err := h.tagUsecase.UnlinkTagFromNote(tagID, noteID); err != nil {
		h.logger.Errorf("Failed to unlink tag from note: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateTag(c *gin.Context) {
	var req UpdateTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Errorf("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTag(id, req.Name)
	if err != nil {
		h.logger.Errorf("Failed to update tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagNameExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

func (h *Handler) UpdateTagForNote(c *gin.Context) {
	var req UpdateTagForNoteRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Errorf("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Invalid note ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTagForNote(tagID, noteID, req.Name)
	if err != nil {
		h.logger.Errorf("Failed to update tag for note: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

// GetNotesByTag
// @Summary		Get notes by tag
// @Tags		Tags
// @Description	Get all notes linked to a specific tag
// @Produce     json
// @Param		tagID path string true 						"Tag ID"
// @Success		200			{object}	[]models.NoteTransfer		"Notes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/{tagID}/notes [get]
func (h *Handler) GetNotesByTag(c *gin.Context) {
	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	notes, err := h.tagUsecase.GetNotesByTag(tagID)
	if err != nil {
		h.logger.Errorf("Failed to get notes by tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if notes == nil {
		notes = []models.Note{} // Return empty array instead of null
	}

	// Convert notes to transfers
	noteTransfers := make([]*models.NoteTransfer, len(notes))
	for i, note := range notes {
		noteTransfers[i] = note.ToTransfer([]string{"r"}) // Default to read access
	}

	c.JSON(http.StatusOK, noteTransfers)
}

// GetTagsByNote
// @Summary		Get tags by note
// @Tags		Tags
// @Description	Get all tags linked to a specific note
// @Produce     json
// @Param		noteID path string true 						"Note ID"
// @Success		200			{object}	[]models.Tag		"Tags"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/notes/{noteID}/tags [get]
func (h *Handler) GetTagsByNote(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("note_id")
	if id == "" {
		h.logger.Infof("Note ID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
		return
	}

	noteID, err := uuid.FromString(id)
	if err != nil {
		h.logger.Errorf("Invalid note ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
		return
	}

	// Check if note exists
	_, err = h.noteUsecase.GetByID(noteID)
	if err != nil {
		h.logger.Errorf("Note not found: %w", err)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Note with ID '%s' not found", noteID)})
		return
	}

	// Check if user has access to the note
	access, err := h.noteUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		h.logger.Errorf("Failed to get user access: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if access == models.EmptyNoteAccess {
		h.logger.Infof("User %s does not have access to note %s", userID, noteID)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to note with ID '%s'", noteID)})
		return
	}

	tags, err := h.tagUsecase.GetTagsByNote(noteID)
	if err != nil {
		h.logger.Errorf("Failed to get tags by note: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	c.JSON(http.StatusOK, tags)
}

// LinkTags
// @Summary		Link two tags together
// @Tags		Tags
// @Description	Create a link between two tags
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		LinkTagsRequest		true	"Tag IDs"
// @Success		200								"Tags linked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		409			{object}	error				"Tags already linked"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/link [post]
func (h *Handler) LinkTags(c *gin.Context) {
	var req LinkTagsRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Errorf("Invalid link tags request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Errorf("Invalid tag1 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag1 ID format: %v", err)})
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Invalid tag2 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag2 ID format: %v", err)})
		return
	}

	if err := h.tagUsecase.LinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Failed to link tags: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

// UnlinkTags
// @Summary		Unlink two tags
// @Tags		Tags
// @Description	Remove the link between two tags
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		LinkTagsRequest		true	"Tag IDs"
// @Success		200								"Tags unlinked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/unlink-tags [post]
func (h *Handler) UnlinkTags(c *gin.Context) {
	var req LinkTagsRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Errorf("Invalid unlink tags request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Errorf("Invalid tag1 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag1 ID format: %v", err)})
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Invalid tag2 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag2 ID format: %v", err)})
		return
	}

	if err := h.tagUsecase.UnlinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Failed to unlink tags: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetLinkedTags
// @Summary		Get linked tags
// @Tags		Tags
// @Description	Get all tags linked to a specific tag
// @Produce     json
// @Param		tagID path string true 						"Tag ID"
// @Success		200			{object}	[]models.Tag		"Linked tags"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/{tagID}/linked [get]
func (h *Handler) GetLinkedTags(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", c.Param("tag_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	// Get the first note associated with this tag to check access
	notes, err := h.tagUsecase.GetNotesByTag(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting notes for tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if len(notes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Tag with ID '%s' not found", tagID)})
		return
	}

	// Check if user has access to at least one note with this tag
	hasAccess := false
	for _, note := range notes {
		access, err := h.noteUsecase.GetUserAccess(note.ID, userID)
		if err != nil {
			continue
		}
		if access != models.EmptyNoteAccess {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to tag with ID '%s'", tagID)})
		return
	}

	tags, err := h.tagUsecase.GetLinkedTags(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting linked tags: %v", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if tags == nil {
		tags = []models.Tag{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, tags)
}

// DeleteTag
// @Summary		Delete tag
// @Tags		Tags
// @Description	Delete a tag and all its relations (notes and linked tags)
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		DeleteTagRequest		true	"Tag ID"
// @Success		200								"Tag deleted"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/delete [post]
func (h *Handler) DeleteTag(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req DeleteTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid delete tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.TagID)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	// Get the first note associated with this tag to check access
	notes, err := h.tagUsecase.GetNotesByTag(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting notes for tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if len(notes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Tag with ID '%s' not found", tagID)})
		return
	}

	// Check if user has access to all note with this tag
	hasAccess := true
	for _, note := range notes {
		access, err := h.noteUsecase.GetUserAccess(note.ID, userID)
		if err != nil {
			continue
		}
		if access != models.ManageAccessNoteAccess && access != models.WriteNoteAccess && access != models.ModifyNoteAccess {
			hasAccess = false
			break
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have sufficient access to delete tag with ID '%s'", tagID)})
		return
	}

	if err := h.tagUsecase.DeleteTag(tagID); err != nil {
		h.logger.Errorf("Error while deleting tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

// LinkExistingTag
// @Summary		Link existing tag to note
// @Tags		Tags
// @Description	Link an existing tag to a note
// @Accept		json
// @Produce     json
// @Param		note_id	path		string		true	"Note ID"
// @Param		tag_id	path		string		true	"Tag ID"
// @Success		201								"Tag linked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		403			{object}	error				"Access denied"
// @Failure		404			{object}	error				"Tag or note not found"
// @Failure		409			{object}	error				"Tag already linked"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/notes/{note_id}/tags/{tag_id} [post]
func (h *Handler) LinkExistingTag(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID, err := uuid.FromString(c.Param("note_id"))
	if err != nil {
		h.logger.Errorf("Invalid note ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
		return
	}

	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	// Check if user has access to the note
	access, err := h.noteUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		h.logger.Errorf("Failed to get user access: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if access == models.EmptyNoteAccess {
		h.logger.Infof("User %s does not have access to note %s", userID, noteID)
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("User does not have access to note with ID '%s'", noteID)})
		return
	}

	err = h.tagUsecase.LinkExistingTag(tagID, noteID)
	if err != nil {
		h.logger.Errorf("Failed to link existing tag: %w", err)
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Tag with ID '%s' not found", e.ID)})
		case *errors.NoteNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Note with ID '%s' not found", e.ID)})
		case *errors.TagLinkExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": "Tag is already linked to this note"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusCreated)
}

type suggestTagsRequest struct {
	Text    string `json:"text" binding:"required"`
	TagsNum *int   `json:"tags_num"`
}

type suggestTagsResponse struct {
	Tags []string `json:"tags"`
}

// SuggestTags godoc
// @Summary Suggest tags for given text
// @Description Generate tag suggestions using LLM
// @Tags tags
// @Accept json
// @Produce json
// @Param request body suggestTagsRequest true "Text to generate tags for"
// @Success 200 {object} suggestTagsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/tags/suggest [post]
func (h *Handler) SuggestTags(c *gin.Context) {
	var req suggestTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tags, err := h.tagUsecase.SuggestTags(req.Text, req.TagsNum)
	if err != nil {
		h.logger.Errorf("Failed to generate tags: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tags"})
		return
	}

	c.JSON(http.StatusOK, suggestTagsResponse{Tags: tags})
}
