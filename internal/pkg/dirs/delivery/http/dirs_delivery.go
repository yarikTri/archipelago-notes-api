package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/dirs"
	"net/http"
	"strconv"
)

type Handler struct {
	dirsUsecase dirs.Usecase
	logger      logger.Logger
}

func NewHandler(du dirs.Usecase, l logger.Logger) *Handler {
	return &Handler{
		dirsUsecase: du,
		logger:      l,
	}
}

// Get
// @Summary		Get dir
// @Tags		Dirs
// @Description	Get dir by ID
// @Produce     json
// @Param		dirID path int true 					"Dir ID"
// @Success		200			{object}	models.Dir		"Dir"
// @Failure		400			{object}	error			"Incorrect input"
// @Failure		500			{object}	error			"Server error"
// @Router		/api/dirs/{dirID} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid dir id '%d'", id)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	dir, err := h.dirsUsecase.Get(id)
	if err != nil {
		h.logger.Errorf("Error while getting dir with id %d: %w", id, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, dir)
}

// GetTree
// @Summary		Get dir Tree
// @Tags		Dirs
// @Description	Get subtree of dir with id {dirID}
// @Produce     json
// @Success		200			{object}	models.DirTree	"Dir tree"
// @Failure		400			{object}	error			"Incorrect input"
// @Failure		500			{object}	error			"Server error"
// @Router		/api/dirs/{dirID}/tree [get]
func (h *Handler) GetTree(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid dir id '%d'", id)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	dirTree, err := h.dirsUsecase.GetTree(id)
	if err != nil {
		h.logger.Errorf("Error while get tree for dir: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, dirTree)
}

// Create
// @Summary		Create dir
// @Tags		Dirs
// @Description	Create dir
// @Accept		json
// @Produce     json
// @Param		dirInfo	body			CreateDirRequest	true	"Dir info"
// @Success		200			{object}	models.Dir			"Dir created"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/dirs [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateDirRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid create dir request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	createdDir, err := h.dirsUsecase.Create(req.Name, req.ParentDirID)
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, createdDir)
}

// Update
// @Summary		Update dir
// @Tags		Dirs
// @Description	Update dir by ID
// @Accept		json
// @Produce     json
// @Param		dirID path int true 						"Dir ID"
// @Param		dirInfo	body			UpdateDirRequest	true	"Dir info"
// @Success		200			{object}	models.Dir			"Updated dir"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/api/dirs/{dirID} [post]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid dir id '%d'", id)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var req UpdateDirRequest
	c.BindJSON(&req)
	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid update dir request: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	reqDir := req.ToDir()
	if id != reqDir.ID {
		h.logger.Infof("Query id (%d) doesn't match dir id (%d)", id, reqDir.ID)
		c.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("Query id (%d) doesn't match dir id (%d)", id, reqDir.ID),
		)
		return
	}

	updatedDir, err := h.dirsUsecase.Update(&reqDir)
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, updatedDir)
}

// Delete
// @Summary		Delete dir
// @Tags		Dirs
// @Description	Delete dir by ID
// @Produce     json
// @Param		dirID path int true 			"Dir ID"
// @Success		200								"Dir deleted"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/api/Dirs/{dirID} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid dir id '%d'", id)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.dirsUsecase.Delete(id); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
