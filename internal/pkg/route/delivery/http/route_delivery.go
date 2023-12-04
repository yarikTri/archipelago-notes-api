package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/image"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
)

type Handler struct {
	routeServices route.Usecase
	imageServices image.Usecase
	logger        logger.Logger
}

func NewHandler(ru route.Usecase, iu image.Usecase, l logger.Logger) *Handler {
	return &Handler{
		routeServices: ru,
		imageServices: iu,
		logger:        l,
	}
}

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

func (h *Handler) List(c *gin.Context) {
	routes, err := h.routeServices.List()
	if err != nil {
		h.logger.Errorf("Error while listing routes: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	routesTransfers := make([]models.RouteTransfer, 0)
	for _, route := range routes {
		routesTransfers = append(routesTransfers, route.ToTransfer())
	}

	c.JSON(http.StatusOK, routesTransfers)
}

func (h *Handler) Search(c *gin.Context) {
	searchQuery := c.Query("route")
	routes, _ := h.routeServices.Search(searchQuery)

	routesTransfers := make([]models.RouteTransfer, 0)
	for _, route := range routes {
		routesTransfers = append(routesTransfers, route.ToTransfer())
	}

	c.JSON(http.StatusOK, routesTransfers)
}

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

	c.JSON(http.StatusOK, nil)
}

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

	c.JSON(http.StatusOK, nil)
}
