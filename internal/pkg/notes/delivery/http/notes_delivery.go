package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/notes"
	"net/http"
)

type Handler struct {
	notesUsecase notes.Usecase
	logger       logger.Logger
}

func NewHandler(nu notes.Usecase, l logger.Logger) *Handler {
	return &Handler{
		notesUsecase: nu,
		logger:       l,
	}
}

// Get
// @Summary		Get note
// @Tags		Notes
// @Description	Get note by ID
// @Produce     json
// @Param		noteID path int true 							"Note ID"
// @Success		200			{object}	models.NoteTransfer		"Note"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/api/notes/{noteID} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	note, err := h.notesUsecase.GetByID(id)
	if err != nil {
		h.logger.Errorf("Error while getting note with id %d: %w", id, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, note.ToTransfer())
}

// List
// @Summary		List notes
// @Tags		Notes
// @Description	Get all notes
// @Produce     json
// @Success		200			{object}	ListNotesResponse	"Notes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/notes [get]
func (h *Handler) List(c *gin.Context) {
	notes, err := h.notesUsecase.List()
	if err != nil {
		h.logger.Errorf("Error while listing notes: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	notesTransfers := make([]*models.NoteTransfer, 0)
	for _, note := range notes {
		notesTransfers = append(notesTransfers, note.ToTransfer())
	}

	c.JSON(http.StatusOK, ListNotesResponse{notesTransfers})
}

// Create
// @Summary		Create note
// @Tags		Notes
// @Description	Create note
// @Accept		json
// @Produce     json
// @Param		noteInfo	body		CreateNoteRequest		true	"Note info"
// @Success		200			{object}	models.NoteTransfer				"Note created"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/api/notes [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateNoteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid create notes request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	createdNote, err := h.notesUsecase.Create(req.DirID, req.AutomergeURL, req.Title)
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, createdNote.ToTransfer())
}

// Update
// @Summary		Update note
// @Tags		Notes
// @Description	Update note by ID
// @Accept		json
// @Produce     json
// @Param		noteID path int true 							"Note ID"
// @Param		noteInfo	body		UpdateNoteRequest		true	"Note info"
// @Success		200			{object}	models.NoteTransfer		"Updated note"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/api/notes/{noteID} [post]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var req UpdateNoteRequest
	c.BindJSON(&req)
	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update note request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reqNote := req.ToNote()
	if id != reqNote.ID {
		h.logger.Infof("Query id (%s) doesn't match note id (%s)", id.String(), reqNote.ID.String())
		c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("Query id (%s) doesn't match note id (%s)", id.String(), reqNote.ID.String()),
		)
		return
	}

	updatedNote, err := h.notesUsecase.Update(reqNote)
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedNote.ToTransfer())
}

// Delete
// @Summary		Delete note
// @Tags		Notes
// @Description	Delete note by ID
// @Produce     json
// @Param		noteID path int true 			"Note ID"
// @Success		200								"Note deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/api/notes/{noteID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.notesUsecase.DeleteByID(id); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
