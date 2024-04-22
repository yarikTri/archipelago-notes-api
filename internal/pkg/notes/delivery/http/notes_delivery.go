package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/common/http/auth"
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

func (h *Handler) checkAccess(c *gin.Context, noteID uuid.UUID, method methodName) *models.NoteAccess {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Infof("Unathorized request for note %s, method %s", noteID.String(), method)
		c.JSON(http.StatusUnauthorized, "")
		return nil
	}

	access, err := h.notesUsecase.GetUserAccess(noteID, userID)
	if err != nil {
		h.logger.Errorf("Error while check access for user with id %s: %w", userID.String(), err)
		c.JSON(http.StatusInternalServerError, "Can't check access")
		return nil
	}

	for _, a := range methodsAccessMap[method] {
		if a == access {
			return &access
		}
	}

	h.logger.Infof("Access forbidden for user %s, note %s, method %s", userID.String(), noteID.String(), method)
	c.JSON(http.StatusForbidden, "")
	return nil
}

// Get
// @Summary		Get note
// @Tags		Notes
// @Description	Get note by ID
// @Produce     json
// @Param		noteID path string true 						"Note ID"
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

	access := h.checkAccess(c, id, getMethodName)
	if access == nil {
		return
	}

	note, err := h.notesUsecase.GetByID(id)
	if err != nil {
		h.logger.Errorf("Error while getting note with id %d: %w", id, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, note.ToTransfer(getAllowedMethods(*access)))
}

// List
// @Summary		List notes
// @Tags		Notes
// @Description	Get all notes user has access to
// @Accept 		json
// @Produce     json
// @Success		200			{object}	ListNotesResponse	"Notes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/notes [get]
func (h *Handler) List(c *gin.Context) {
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Infof("Unathorized request for listing notes")
		c.JSON(http.StatusUnauthorized, "")
		return
	}

	notes, err := h.notesUsecase.List(userID)
	if err != nil {
		h.logger.Errorf("Error while listing notes: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	notesTransfers := make([]*models.NoteTransfer, 0)
	for _, note := range notes {
		access := models.NoteAccessFromString(*note.Access)
		notesTransfers = append(notesTransfers, note.ToTransfer(getAllowedMethods(access)))
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
	userID, err := auth.GetUserId(c)
	if err != nil {
		h.logger.Infof("Unathorized request for listing notes")
		c.JSON(http.StatusUnauthorized, "")
		return
	}

	var req CreateNoteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid create notes request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	createdNote, err := h.notesUsecase.Create(req.DirID, req.AutomergeURL, req.Title, userID)
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, createdNote.ToTransfer(getAllowedMethods(models.ManageAccessNoteAccess)))
}

// Update
// @Summary		Update note
// @Tags		Notes
// @Description	Update note by ID
// @Accept		json
// @Produce     json
// @Param		noteID path string true 						"Note ID"
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

	access := h.checkAccess(c, id, updateMethodName)
	if access == nil {
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

	c.JSON(http.StatusOK, updatedNote.ToTransfer(getAllowedMethods(*access)))
}

// Delete
// @Summary		Delete note
// @Tags		Notes
// @Description	Delete note by ID
// @Produce     json
// @Param		noteID path string true 		"Note ID"
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

// SetAccess
// @Summary		Set Access
// @Tags		Notes
// @Description	Set access to note to user
// @Accept		json
// @Produce     json
// @Param		noteID path string true 				"Note ID"
// @Param		userID path string true 				"User to set access ID"
// @Param		access	body	SetAccessRequest true	"Note info"
// @Success		200										"Note deleted"
// @Failure		400			{object}	error			"Incorrect input"
// @Failure		500			{object}	error			"Server error"
// @Router		/api/notes/{noteID}/access/{userID} [post]
func (h *Handler) SetAccess(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if access := h.checkAccess(c, id, setAccessMethodName); access == nil {
		return
	}

	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		h.logger.Infof("Invalid note id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var req SetAccessRequest
	c.BindJSON(&req)
	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid set access request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.notesUsecase.SetUserAccess(id, userID, models.NoteAccessFromString(req.Access), req.WithInvitation); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
