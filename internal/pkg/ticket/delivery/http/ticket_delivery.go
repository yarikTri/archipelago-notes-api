package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/web-transport-cards/internal/pkg/ticket"
)

type Handler struct {
	services ticket.Usecase
	logger   logger.Logger
}

func NewHandler(tu ticket.Usecase, l logger.Logger) *Handler {
	return &Handler{
		services: tu,
		logger:   l,
	}
}

func (h *Handler) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) GetByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ticket": "get"})
}

func (h *Handler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ticket": "list"})
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ticket": "create"})
}

func (h *Handler) DeleteByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ticket": "delete"})
}
