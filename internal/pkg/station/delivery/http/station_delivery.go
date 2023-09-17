package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/pkg/station"
)

type Handler struct {
	services station.Usecase
	logger   logger.Logger
}

func NewHandler(su station.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: su,
		logger:   l,
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"station": "get"})
}

func (h *Handler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"station": "list"})
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"station": "create"})
}

func (h *Handler) DeleteByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"station": "delete"})
}
