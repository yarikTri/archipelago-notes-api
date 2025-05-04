package http

import (
	"fmt"
	"net/http"
	"strings"

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

func (h *Handler) checkTagOwnership(c *gin.Context, tagID uuid.UUID, userID uuid.UUID) error {
	// Check if the tag belongs to the user
	isOwner, err := h.tagUsecase.IsTagUsers(userID, tagID)

	if err != nil {
		h.logger.Errorf("Failed to check tag ownership: %w", err)

		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if !isOwner {
		h.logger.Infof("User %s does not own tag %s", userID, tagID)
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not own this tag"})
		return fmt.Errorf("user does not own tag")
	}

	return nil
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

	tag, err := h.tagUsecase.CreateAndLinkTag(req.Name, noteID, userID)
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
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := h.checkTagOwnership(c, tagID, userID); err != nil {
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

// UpdateTag
// @Summary		Update tag
// @Tags		Tags
// @Description	Update the name of an existing tag
// @Accept		json
// @Produce     json
// @Param		request	body		UpdateTagRequest		true	"Tag update request"
// @Success		200			{object}	models.Tag		"Updated tag"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		403			{object}	error				"Access denied"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		409			{object}	error				"Tag name conflict"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/update [put]
func (h *Handler) UpdateTag(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := h.checkTagOwnership(c, id, userID); err != nil {
		return
	}

	updatedTag, err := h.tagUsecase.UpdateTag(id, req.Name, userID)
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

// func (h *Handler) UpdateTagForNote(c *gin.Context) {
// 	var req UpdateTagForNoteRequest
// 	if err := c.BindJSON(&req); err != nil {
// 		h.logger.Errorf("Failed to bind request: %w", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := req.validate(); err != nil {
// 		h.logger.Errorf("Invalid update tag request: %w", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	tagID, err := uuid.FromString(req.TagID)
// 	if err != nil {
// 		h.logger.Errorf("Invalid tag ID format: %w", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
// 		return
// 	}

// 	noteID, err := uuid.FromString(req.NoteID)
// 	if err != nil {
// 		h.logger.Errorf("Invalid note ID format: %w", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid note ID format: %v", err)})
// 		return
// 	}

// 	updatedTag, err := h.tagUsecase.UpdateTagForNote(tagID, noteID, req.Name)
// 	if err != nil {
// 		h.logger.Errorf("Failed to update tag for note: %w", err)
// 		switch e := err.(type) {
// 		case *errors.TagNotFoundError:
// 			c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
// 		case *errors.TagNameEmptyError:
// 			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
// 		default:
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, updatedTag)
// }

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
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := uuid.FromString(c.Param("tag_id"))
	if err != nil {
		h.logger.Errorf("Invalid tag ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag ID format: %v", err)})
		return
	}

	if err := h.checkTagOwnership(c, tagID, userID); err != nil {
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
// @Router		/api/tags/note/{noteID} [get]
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

	tags, err := h.tagUsecase.GetTagsByNoteForUser(noteID, userID)
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
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := h.checkTagOwnership(c, tag1ID, userID); err != nil {
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Invalid tag2 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag2 ID format: %v", err)})
		return
	}

	if err := h.checkTagOwnership(c, tag2ID, userID); err != nil {
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
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := h.checkTagOwnership(c, tag1ID, userID); err != nil {
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Invalid tag2 ID format: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid tag2 ID format: %v", err)})
		return
	}

	if err := h.checkTagOwnership(c, tag2ID, userID); err != nil {
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

	if err := h.checkTagOwnership(c, tagID, userID); err != nil {
		return
	}

	tags, err := h.tagUsecase.GetLinkedTagsForUser(tagID, userID)
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

	if err := h.checkTagOwnership(c, tagID, userID); err != nil {
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

// LinkTagToNote
// @Summary		Link existing tag to note
// @Tags		Tags
// @Description	Link an existing tag to a note
// @Produce     json
// @Param		note_id	path		string		true	"Note ID"
// @Param		tag_id	path		string		true	"Tag ID"
// @Success		201								"Tag linked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		403			{object}	error				"Access denied"
// @Failure		404			{object}	error				"Tag or note not found"
// @Failure		409			{object}	error				"Tag already linked"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/{tag_id}/link/{note_id} [post]
func (h *Handler) LinkTagToNote(c *gin.Context) {
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

	if err := h.checkTagOwnership(c, tagID, userID); err != nil {
		return
	}

	err = h.tagUsecase.LinkTagToNote(tagID, noteID)
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
// @Tags Tags
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

	// Explicitly check for empty text
	if strings.TrimSpace(req.Text) == "" {
		h.logger.Infof("Empty text provided for tag suggestion")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text cannot be empty"})
		return
	}

	tags, err := h.tagUsecase.SuggestTags(req.Text, req.TagsNum)
	if err != nil {
		h.logger.Errorf("Failed to generate tags: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate tags: %v", err)})
		return
	}

	// Ensure tags is never nil
	if tags == nil {
		tags = []string{}
	}

	c.JSON(http.StatusOK, suggestTagsResponse{Tags: tags})
}

var DEFAULT_LIMIT_LIST_CLOSEST uint32 = 3

type listClosestTagsRequest struct {
	Limit   *uint32 `json:"limit"`
	TagName string  `json:"name"`
}

// ListClosestTags
// @Summary		List closest tags
// @Tags		Tags
// @Description	Get a list of tags closest to the given tag
// @Accept json
// @Produce     json
// @Param		limit	body		listClosestTagsRequest		true	"Limit and name"				"Tag ID"
// @Success		200			{object}	[]models.Tag		"Closest tags"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/closest [post]
func (h *Handler) ListClosestTags(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Errorf("Failed to get user ID: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: maybe check that user exists.

	limit := DEFAULT_LIMIT_LIST_CLOSEST

	var req listClosestTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Limit != nil {
		limit = *req.Limit
	}

	if req.TagName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cant be empty"})
		return
	}

	tags, err := h.tagUsecase.ListClosestTags(req.TagName, userID, limit)
	if err != nil {
		h.logger.Errorf("Failed to list closest tags: %w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if tags == nil {
		tags = []models.Tag{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, tags)
}
