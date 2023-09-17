package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/pkg/route"
)

type Handler struct {
	services route.Usecase
	logger   logger.Logger
}

func NewHandler(ru route.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: ru,
		logger:   l,
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "get"})
}

func (h *Handler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "list"})
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "create"})
}

func (h *Handler) DeleteByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "delete"})
}
