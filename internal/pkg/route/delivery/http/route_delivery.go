package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/auth"
	"github.com/yarikTri/web-transport-cards/internal/pkg/image"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"

	commonHttp "github.com/yarikTri/web-transport-cards/internal/common/http"
)

type Handler struct {
	routeServices  route.Usecase
	imageServices  image.Usecase
	ticketServices ticket.Usecase
	authServices   auth.Usecase
	logger         logger.Logger
}

func NewHandler(ru route.Usecase, iu image.Usecase, tu ticket.Usecase, au auth.Usecase, l logger.Logger) *Handler {
	return &Handler{
		routeServices:  ru,
		imageServices:  iu,
		ticketServices: tu,
		authServices:   au,
		logger:         l,
	}
}

// @Summary		Get route
// @Tags		Routes
// @Description	Get route by ID
// @Produce     json
// @Param		routeID path int true 							"Route ID"
// @Success		200			{object}	models.RouteTransfer	"Got route"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/routes/{routeID} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
		return
	}

	route, err := h.routeServices.GetByID(int(id))
	if err != nil {
		h.logger.Errorf("Error while getting route with id %d: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, route.ToTransfer())
}

// @Summary		List routes
// @Tags		Routes
// @Description	Get all routes
// @Produce     json
// @Success		200			{object}	ListRoutesResponse	"Got routes"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/routes [get]
func (h *Handler) List(c *gin.Context) {
	searchQuery := c.Query("route")
	routes, _ := h.routeServices.Search(searchQuery)

	routesTransfers := make([]models.RouteTransfer, 0)
	for _, route := range routes {
		routesTransfers = append(routesTransfers, route.ToTransfer())
	}

	var draftTicketID *int = nil
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err == nil { // РАВНО!
		user, errr := h.authServices.GetUserBySessionID(sessionID)
		if errr == nil { // РАВНО!
			h.logger.Infof("User not found")
			userID := int(user.ID)
			draftTicketID = h.getDraftTicketID(userID)
		}
	}

	c.JSON(http.StatusOK, ListRoutesResponse{draftTicketID, routesTransfers})
}

// @Summary		Create route
// @Tags		Routes
// @Description	Create route
// @Accept		json
// @Produce     json
// @Param		routeInfo	body		CreateRouteRequest		true	"Route info"
// @Success		200			{object}	models.RouteTransfer			"Route created"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/routes [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRouteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid create route request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid create route request: %s", err.Error()))
		return
	}

	createdRoute, err := h.routeServices.Create(req.ToRoute())
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, createdRoute.ToTransfer())
}

// @Summary		Update route
// @Tags		Routes
// @Description	Update route by ID
// @Accept		json
// @Produce     json
// @Param		routeID path int true 							"Route ID"
// @Param		routeInfo	body		UpdateRouteRequest		true	"Route info"
// @Success		200			{object}	models.RouteTransfer			"Updated route"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/routes/{routeID} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
		return
	}

	var req UpdateRouteRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update route request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid update route request: %s", err.Error()))
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

// @Summary		Delete route
// @Tags		Routes
// @Description	Delete route by ID
// @Produce     json
// @Param		routeID path int true 			"Route ID"
// @Success		200								"Route deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/routes/{routeID} [delete]
func (h *Handler) DeleteByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
		return
	}

	if err := h.routeServices.DeleteByID(int(id)); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Put route image
// @Tags		Routes
// @Description	Put image of route by ID
// @Accept 		multipart/form-data
// @Produce     json
// @Param		image formData file true 					"Route image"
// @Success		200											"Route image updated"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/routes/{routeID}/image [put]
func (h *Handler) PutImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Infof("Invalid route id '%s'", c.Param("id"))
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid route id '%s'", c.Param("id")))
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
