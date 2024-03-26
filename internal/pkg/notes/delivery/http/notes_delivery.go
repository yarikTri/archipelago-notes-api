package http

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-nodes-api/internal/models"
	"github.com/yarikTri/archipelago-nodes-api/internal/pkg/notes"
)

type Handler struct {
	nodeServices node.Usecase
	logger       logger.Logger
}

func NewHandler(nu node.Usecase, l logger.Logger) *Handler {
	return &Handler{
		nodeServices: nu,
		logger:       l,
	}
}

// GetByID
// @Summary		Get notes
// @Tags		Nodes
// @Description	Get notes by ID
// @Produce     json
// @Param		nodeID path int true 							"Note ID"
// @Success		200			{object}	models.NoteTransfer		"Note"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/nodes/{nodeID} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid notes id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid notes id '%s'", c.Param("id")))
		return
	}

	node, err := h.nodeServices.GetByID(id)
	if err != nil {
		h.logger.Errorf("Error while getting notes with id %d: %w", id, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, node.ToTransfer())
}

// List
// @Summary		List nodes
// @Tags		Nodes
// @Description	Get all nodes
// @Produce     json
// @Success		200			{object}	ListNodesResponse	"Nodes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/nodes [get]
func (h *Handler) List(c *gin.Context) {
	nodes, err := h.nodeServices.List()
	if err != nil {
		h.logger.Errorf("Error while listing nodes: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	nodeTransfers := make([]models.NoteTransfer, 0)
	for _, node := range nodes {
		nodeTransfers = append(nodeTransfers, node.ToTransfer())
	}

	c.JSON(http.StatusOK, ListNodesResponse{nodeTransfers})
}

// Create
// @Summary		Create notes
// @Tags		Nodes
// @Description	Create notes
// @Accept		json
// @Produce     json
// @Param		nodeInfo	body		CreateNoteRequest		true	"Note info"
// @Success		200			{object}	models.NoteTransfer				"Note created"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/nodes [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateNoteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid create notes request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid create notes request: %s", err.Error()))
		return
	}

	createdRoute, err := h.nodeServices.Create(req.ToRoute())
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, createdRoute.ToTransfer())
}

// Update
// @Summary		Update notes
// @Tags		Nodes
// @Description	Update notes by ID
// @Accept		json
// @Produce     json
// @Param		routeID path int true 							"Note ID"
// @Param		routeInfo	body		UpdateRouteRequest		true	"Note info"
// @Success		200			{object}	models.NoteTransfer			"Updated notes"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/routes/{routeID} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid notes id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid notes id '%s'", c.Param("id")))
		return
	}

	var req UpdateNoteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update notes request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid update notes request: %s", err.Error()))
		return
	}

	updatedRoute, err := h.routeServices.Update(req.ToRoute(id))
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, updatedRoute.ToTransfer())
}

// @Summary		Delete notes
// @Tags		Routes
// @Description	Delete notes by ID
// @Produce     json
// @Param		routeID path int true 			"Route ID"
// @Success		200								"Route deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/routes/{routeID} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusUnauthorized, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID))
	if !isModerator {
		h.logger.Infof("Forbidden to delete notes")
		c.JSON(http.StatusForbidden, fmt.Sprintf("Forbidden to delete notes"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid notes id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid notes id '%s'", c.Param("id")))
		return
	}

	if err := h.routeServices.DeleteByID(int(id)); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Put notes image
// @Tags		Routes
// @Description	Put image of notes by ID
// @Accept 		multipart/form-data
// @Produce     json
// @Param		image formData file true 					"Route image"
// @Success		200											"Route image updated"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/routes/{routeID}/image [put]
func (h *Handler) PutImage(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("No session cookie")
		c.JSON(http.StatusUnauthorized, "No session cookie")
		return
	}

	user, err := h.authServices.GetUserBySessionID(sessionID)
	if err != nil {
		h.logger.Infof("User not found")
		c.JSON(http.StatusBadRequest, "User not found")
		return
	}

	isModerator, _ := h.authServices.CheckUserIsModerator(int(user.ID))
	if !isModerator {
		h.logger.Infof("Forbidden to update notes's image")
		c.JSON(http.StatusForbidden, fmt.Sprintf("Forbidden to update notes's image"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid notes id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid notes id '%s'", c.Param("id")))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Infof("Can't parse multipart form")
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Can't parse multipart form"))
		return
	}
	fileHeader := form.File["image"][0]

	image, err := form.File["image"][0].Open()
	if err != nil {
		h.logger.Infof("Can't get image from request")
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Can't get image from request"))
		return
	}

	imageUUID, err := h.imageServices.Put(c, image, fileHeader.Size)
	if err != nil {
		h.logger.Errorf("Can't save image: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Can't save image: %s", err.Error()))
		return
	}

	if err := h.routeServices.UpdateImageUUID(int(id), imageUUID); err != nil {
		h.logger.Errorf("Can't save image: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Can't save image: %s", err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getDraftTicketID(userID int) *int {
	foundTicket, err := h.ticketServices.GetDraft(userID)
	if err != nil {
		return nil
	}

	ticketID := int(foundTicket.ID)
	return &ticketID
}
