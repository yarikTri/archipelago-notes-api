package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/models"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
	"github.com/yarikTri/web-transport-cards/internal/pkg/station"
)

type Handler struct {
	routeServices   route.Usecase
	stationServices station.Usecase
	logger          logger.Logger
}

func NewHandler(ru route.Usecase, su station.Usecase, l logger.Logger) *Handler {
	return &Handler{
		routeServices:   ru,
		stationServices: su,
		logger:          l,
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	_route, _ := h.routeServices.GetByID(uint32(id))

	stations, _ := h.stationServices.ListByRoute(_route.ID)

	c.HTML(http.StatusOK, "route.tmpl", gin.H{
		"routes": []models.RouteTransfer{_route.ToTransfer(stations)},
	})
}

func (h *Handler) List(c *gin.Context) {
	routes, _ := h.routeServices.List()

	routesTransfers := make([]models.RouteTransfer, 0)
	for _, route := range routes {
		stations, _ := h.stationServices.ListByRoute(route.ID)

		routesTransfers = append(routesTransfers, route.ToTransfer(stations))
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"routes": routesTransfers,
	})
}

func (h *Handler) Search(c *gin.Context) {
	searchQuery := c.Query("route")
	routes, _ := h.routeServices.Search(searchQuery)

	routesTransfers := make([]models.RouteTransfer, 0)
	for _, route := range routes {
		stations, _ := h.stationServices.ListByRoute(route.ID)

		routesTransfers = append(routesTransfers, route.ToTransfer(stations))
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"routes": routesTransfers,
		"filter": searchQuery,
	})
}

func (h *Handler) DeleteByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "delete"})
}
