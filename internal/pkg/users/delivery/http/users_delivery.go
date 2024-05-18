package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	"github.com/gofrs/uuid/v5"
	"github.com/yarikTri/archipelago-notes-api/internal/models"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/users"
	"net/http"
	"strconv"
)

type Handler struct {
	usersUsecase users.Usecase
	logger       logger.Logger
}

func NewHandler(uu users.Usecase, l logger.Logger) *Handler {
	return &Handler{
		usersUsecase: uu,
		logger:       l,
	}
}

// Get
// @Summary		Get user
// @Tags		Users
// @Description	Get user by user id
// @Produce     json
// @Param		userID path string true 							"User ID"
// @Success		200			{object}	models.UserTransfer		"User"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/api/users/{userID} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		h.logger.Infof("Invalid user id '%s'", id)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.usersUsecase.GetByID(id)
	if err != nil {
		h.logger.Errorf("Error while getting user with id %d: %w", id, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user.ToTransfer())
}

// Search
// @Summary		Search users
// @Tags		Users
// @Description	Search users by query
// @Produce     json
// @Param		q query string true 							"Query of search"
// @Success		200			{object}	SearchUsersResponse		"Found users"
// @Failure		400			{object}	error					"Incorrect input"
// @Failure		500			{object}	error					"Server error"
// @Router		/api/users/ [get]
func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	users, err := h.usersUsecase.Search(query)
	if err != nil {
		h.logger.Errorf("Error while searching users by query %s: %w", query, err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	userTransfers := make([]*models.UserTransfer, 0)
	for _, u := range users {
		userTransfers = append(userTransfers, u.ToTransfer())
	}

	c.JSON(http.StatusOK, SearchUsersResponse{Users: userTransfers})
}

// SetRootDirID
// @Summary		Set root dir id
// @Tags		Users
// @Description	Set root dir id by user id
// @Param		userID path string true 								"User ID"
// @Param		rootDirID path int true 								"Root dir ID"
// @Success		200			{object}	string							"Root dir setted"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/api/users/{userID}/root_dir/{rootDirID} [post]
func (h *Handler) SetRootDirID(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		h.logger.Infof("Invalid user id '%s'", userID)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	rootDirID, err := strconv.Atoi(c.Param("rootDirID"))
	if err != nil {
		h.logger.Infof("Invalid root dir id '%s'", userID)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.usersUsecase.SetRootDirByID(userID, rootDirID); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

// SendEmailConfirmation
// @Summary		Send email confirmation
// @Tags		Users
// @Description	Send user's email confirmation
// @Param		userID path string true 								"User ID"
// @Success		200			{object}	string							"Mail sent"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/api/users/{userID}/send_email_confirmation [post]
func (h *Handler) SendEmailConfirmation(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		h.logger.Infof("Invalid user id '%s'", userID)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.usersUsecase.SendEmailConfirmation(userID); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, []byte{})
}

// ConfirmEmail
// @Summary		Confirm email
// @Tags		Users
// @Description	Confirm user's email
// @Param		userID path string true 								"User ID"
// @Success		200			{object}	string							"Email confirmed"
// @Failure		400			{object}	error							"Incorrect input"
// @Failure		500			{object}	error							"Server error"
// @Router		/api/users/{userID}/confirm_email [post]
func (h *Handler) ConfirmEmail(c *gin.Context) {
	userID, err := uuid.FromString(c.Param("userID"))
	if err != nil {
		h.logger.Infof("Invalid user id '%s'", userID)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.usersUsecase.ConfirmEmail(userID); err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "")
}
