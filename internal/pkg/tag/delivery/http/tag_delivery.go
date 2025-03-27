package http

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/tag/errors"
)

type Handler struct {
	tagUsecase tag.Usecase
	logger     logger.Logger
}

func NewHandler(tu tag.Usecase, l logger.Logger) *Handler {
	return &Handler{
		tagUsecase: tu,
		logger:     l,
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
	var req CreateAndLinkTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Failed to parse note ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid note ID format")
		return
	}

	tag, err := h.tagUsecase.CreateAndLinkTag(req.Name, noteID)
	if err != nil {
		h.logger.Errorf("Error while creating and linking tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, e.Error())
		case *errors.TagNameExistsError:
			c.JSON(http.StatusConflict, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
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
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Errorf("Failed to parse note ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid note ID format")
		return
	}

	if err := h.tagUsecase.UnlinkTagFromNote(tagID, noteID); err != nil {
		h.logger.Errorf("Error while unlinking tag from note: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		case *errors.TagLinkNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdateTag
// @Summary		Update tag
// @Tags		Tags
// @Description	Update tag by ID
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		UpdateTagRequest		true	"Tag info"
// @Success		200								"Tag updated"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		409			{object}	error				"Tag name already exists"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags [post]
func (h *Handler) UpdateTag(c *gin.Context) {
	var req UpdateTagRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid request data")
		return
	}

	id, err := uuid.FromString(req.ID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.ID)
		c.JSON(http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	if err := h.tagUsecase.UpdateTag(id, req.Name); err != nil {
		h.logger.Errorf("Error while updating tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		case *errors.TagNameExistsError:
			c.JSON(http.StatusConflict, e.Error())
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.Status(http.StatusOK)
}

// UpdateTagForNote
// @Summary		Update tag for note
// @Tags		Tags
// @Description	Update tag name for a specific note
// @Accept		json
// @Produce     json
// @Param		tagInfo	body		UpdateTagForNoteRequest		true	"Tag info"
// @Success		200								"Tag updated"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/note [post]
func (h *Handler) UpdateTagForNote(c *gin.Context) {
	var req UpdateTagForNoteRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update tag request: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid request data")
		return
	}

	tagID, err := uuid.FromString(req.TagID)
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", req.TagID)
		c.JSON(http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	noteID, err := uuid.FromString(req.NoteID)
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", req.NoteID)
		c.JSON(http.StatusBadRequest, "Invalid note ID format")
		return
	}

	if err := h.tagUsecase.UpdateTagForNote(tagID, noteID, req.Name); err != nil {
		h.logger.Errorf("Error while updating tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		case *errors.TagNameEmptyError:
			c.JSON(http.StatusBadRequest, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetNotesByTag
// @Summary		Get notes by tag
// @Tags		Tags
// @Description	Get all notes linked to a specific tag
// @Produce     json
// @Param		tagID path string true 						"Tag ID"
// @Success		200			{object}	[]models.Note		"Notes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/{tagID}/notes [get]
func (h *Handler) GetNotesByTag(c *gin.Context) {
	tagID, err := uuid.FromString(c.Param("tagID"))
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", c.Param("tagID"))
		c.JSON(http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	notes, err := h.tagUsecase.GetNotesByTag(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting notes by tag: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, notes)
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
	noteID, err := uuid.FromString(c.Param("noteID"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("noteID"))
		c.JSON(http.StatusBadRequest, "Invalid note ID format")
		return
	}

	tags, err := h.tagUsecase.GetTagsByNote(noteID)
	if err != nil {
		h.logger.Errorf("Error while getting tags by note: %w", err)
		c.JSON(http.StatusInternalServerError, "Internal server error")
		return
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
		c.JSON(http.StatusBadRequest, "Invalid request format")
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag1 ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid tag1 ID format")
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag2 ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid tag2 ID format")
		return
	}

	if err := h.tagUsecase.LinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Error while linking tags: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		case *errors.TagLinkExistsError:
			c.JSON(http.StatusConflict, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
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
// @Param		tagInfo	body		UnlinkTagsRequest		true	"Tag IDs"
// @Success		200								"Tags unlinked"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		404			{object}	error				"Tag not found"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/tags/unlink-tags [post]
func (h *Handler) UnlinkTags(c *gin.Context) {
	var req UnlinkTagsRequest
	if err := c.BindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid request format")
		return
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tag1ID, err := uuid.FromString(req.Tag1ID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag1 ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid tag1 ID format")
		return
	}

	tag2ID, err := uuid.FromString(req.Tag2ID)
	if err != nil {
		h.logger.Errorf("Failed to parse tag2 ID: %w", err)
		c.JSON(http.StatusBadRequest, "Invalid tag2 ID format")
		return
	}

	if err := h.tagUsecase.UnlinkTags(tag1ID, tag2ID); err != nil {
		h.logger.Errorf("Error while unlinking tags: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		case *errors.TagToTagLinkNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
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
	tagID, err := uuid.FromString(c.Param("tagID"))
	if err != nil {
		h.logger.Infof("Invalid tag id '%s'", c.Param("tagID"))
		c.JSON(http.StatusBadRequest, "Invalid tag ID format")
		return
	}

	tags, err := h.tagUsecase.GetLinkedTags(tagID)
	if err != nil {
		h.logger.Errorf("Error while getting linked tags: %w", err)

		// Handle specific error cases
		switch e := err.(type) {
		case *errors.TagNotFoundError:
			c.JSON(http.StatusNotFound, e.Error())
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, tags)
}
