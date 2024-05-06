package http

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/summary"
)

type Handler struct {
	sumUsecase summary.Usecase
	logger     logger.Logger
}

func NewHandler(nu summary.Usecase, l logger.Logger) *Handler {
	return &Handler{
		sumUsecase: nu,
		logger:     l,
	}
}

// TODO: add token for microservice communication
func (h *Handler) SaveSummaryText(c *gin.Context) {
	type SaveSummaryRequest struct {
		ID           string `json:"id" valid:"required"`
		Text         string `json:"text"`
		Active       bool   `json:"active" valid:"required"`
		Platform     string `json:"platform" valid:"required"`
		Detalization string `json:"detalization" valid:"required"`
	}

	var req SaveSummaryRequest
	req.Active = true // default value
	c.BindJSON(&req)

	if _, err := valid.ValidateStruct(req); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}

	id, err := uuid.FromString(req.ID)
	if err != nil {
		h.logger.Errorf("Failed to cast id to uuid %s: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	summ, err := h.sumUsecase.SaveSummaryText(id, req.Text, req.Active, models.DetalizationFromString(req.Detalization), req.Platform)
	if err != nil {
		h.logger.Errorf("Error while saving summary: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, summ.ToTransfer())
}

// TODO: add token for microservice communication
func (h *Handler) UpdateSummaryTextRole(c *gin.Context) {
	type UpdateSummaryTextRoleRequest struct {
		ID           string `json:"id" valid:"required"`
		TextWithRole string `json:"text_with_role" valid:"required"`
		Role         string `json:"role" valid:"required"`
	}

	var req UpdateSummaryTextRoleRequest
	c.BindJSON(&req)

	id, err := uuid.FromString(req.ID)
	if err != nil {
		h.logger.Errorf("Failed to cast id to uuid %s: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
	}

	err = h.sumUsecase.UpdateSummaryTextRole(id, req.TextWithRole, req.Role)
	if err != nil {
		h.logger.Errorf("Error while saving summary: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// TODO: add token for microservice communication
func (h *Handler) FinishSummary(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Failed to cast id to uuid %s: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
	}

	err = h.sumUsecase.FinishSummary(id)
	if err != nil {
		h.logger.Errorf("Error while finishing summary: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// TODO: check permissions by user
func (h *Handler) GetSummary(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Failed to cast id to uuid %s: %w", id, err)
		c.JSON(http.StatusBadRequest, err)
	}

	summ, err := h.sumUsecase.GetSummary(id)
	if err != nil {
		h.logger.Errorf("Error while getting summary: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, summ.ToTransfer())
}

// TODO: add token for microservice communication
func (h *Handler) GetActiveSummaries(c *gin.Context) {
	summaries, err := h.sumUsecase.GetActiveSummaries()
	if err != nil {
		h.logger.Errorf("Error while getting summary: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	summaryTransferList := make([]models.SummaryTransfer, len(summaries))
	for i, summary := range summaries {
		summaryTransferList[i] = *summary.ToTransfer()
	}

	c.JSON(http.StatusOK, summaryTransferList)
}
