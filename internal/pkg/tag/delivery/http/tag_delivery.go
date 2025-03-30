package http

import (
	"fmt"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
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
	userID, err := uuid.FromString(c.GetHeader("X-User-Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Name   string `json:"name" binding:"required"`
		NoteID string `json:"note_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	// Check if note exists
	_, err = h.noteUsecase.GetByID(noteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Check if user has access to the note
	access, err := h.noteUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if access == models.EmptyNoteAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have access to this note"})
		return
	}

	tag, err := h.tagUsecase.CreateAndLinkTag(req.Name, noteID)
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Failed to parse note ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	if err := h.tagUsecase.UnlinkTagFromNote(tagID, noteID); err != nil {
		h.logger.Errorf("Error while unlinking tag from note: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UpdateTag(c *gin.Context) {
	var req UpdateTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	id, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.TagID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTag(id, req.Name)
	if err != nil {
		h.logger.Errorf("Error while updating tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagNameExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedTag)
}

func (h *Handler) UpdateTagForNote(c *gin.Context) {
	var req UpdateTagForNoteRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.TagID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", req.NoteID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTagForNote(tagID, noteID, req.Name)
	if err != nil {
		h.logger.Errorf("Error while updating tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
		h.logger.Infof("Invalid tag id '%s'", c.Param("tag_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	notes, err := h.tagUsecase.GetNotesByTag(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting notes by tag: %v", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
	id := c.Param("note_id")
	if id == "" {
		h.logger.Infof("Note ID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
		return
	}

	noteID, err := uuid.FromString(id)
	if err != nil {
		h.logger.Infof("Invalid note id '%s': %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	// Check if note exists
	_, err = h.noteUsecase.GetByID(noteID)
	if err != nil {
		h.logger.Infof("Note not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	tags, err := h.tagUsecase.GetTagsByNote(noteID)
	if err != nil {
		h.logger.Errorf("Error while getting tags by note: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Infof("Invalid link tags request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Infof("Invalid tag1 id '%s'", req.Tag1ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag1 ID format"})
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Infof("Invalid tag2 id '%s'", req.Tag2ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag2 ID format"})
		return
	}

	if err := h.tagUsecase.LinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Error while linking tags: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkExistsError:
			c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Infof("Invalid unlink tags request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Infof("Invalid tag1 id '%s'", req.Tag1ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag1 ID format"})
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Infof("Invalid tag2 id '%s'", req.Tag2ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag2 ID format"})
		return
	}

	if err := h.tagUsecase.UnlinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Error while unlinking tags: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		case *errors.TagLinkNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", c.Param("tag_id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
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
	var req DeleteTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid delete tag request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.TagID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	if err := h.tagUsecase.DeleteTag(tagID); err != nil {
		h.logger.Errorf("Error while deleting tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) LinkExistingTag(c *gin.Context) {
	userID, err := uuid.FromString(c.GetHeader("X-User-Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	noteID, err := uuid.FromString(c.Param("note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// Check if user has access to the note
	access, err := h.noteUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if access == models.EmptyNoteAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have access to this note"})
		return
	}

	err = h.tagUsecase.LinkExistingTag(tagID, noteID)
	if err != nil {
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
