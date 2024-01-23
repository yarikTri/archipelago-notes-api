package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	commonHttp "github.com/yarikTri/web-transport-cards/internal/common/http"
	"github.com/yarikTri/web-transport-cards/internal/pkg/auth"
)

const SESSION_DURATION_NANOSECONDS = 7889400000000000
const NANOSECONDS_IN_SECOND = 1000000000

type Handler struct {
	authServices auth.Usecase
	logger       logger.Logger
}

func NewHandler(au auth.Usecase, l logger.Logger) *Handler {
	return &Handler{
		authServices: au,
		logger:       l,
	}
}

// @Summary		Sign Up
// @Tags		Auth
// @Description	Create account
// @Accept		json
// @Produce     json
// @Param		req	body	SignUpRequest		true		"User info"
// @Success		200			{object}	models.UserTransfer	"Created user"
// @Failure		400			{object}	error				"Incorrect input"
// @Failure		500			{object}	error				"Server error"
// @Router		/auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var req SignUpRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid sign up request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid sign up request: %s", err.Error()))
		return
	}

	createdUser, err := h.authServices.SignUp(req.toUser())
	if err != nil {
		h.logger.Errorf("Error: %w", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, createdUser.ToTransfer())
}

// @Summary		Login
// @Tags		Auth
// @Description	Create session
// @Accept		json
// @Param		req	body	LoginRequest	true	"Username and password"
// @Success		200									"User logined"
// @Failure		400			{object}	error		"Incorrect input"
// @Failure		500			{object}	error		"Server error"
// @Router		/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	c.BindJSON(&req)

	if err := req.validate(); err != nil {
		h.logger.Infof("Invalid login request: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid login request: %s", err.Error()))
		return
	}

	sessionID, user, err := h.authServices.Login(req.Username, req.Password, time.Duration(SESSION_DURATION_NANOSECONDS))
	if err != nil {
		h.logger.Errorf("Error while login: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Error while login: %s", err.Error()))
		return
	}

	c.SetCookie(commonHttp.AUTH_COOKIE_NAME, sessionID, SESSION_DURATION_NANOSECONDS/NANOSECONDS_IN_SECOND, "", "", false, false)
	c.JSON(http.StatusOK, user.ToTransfer())
}

// @Summary		Logout
// @Tags		Auth
// @Description	Logout
// @Accept		json
// @Success		200								"User logined"
// @Failure		400			{object}	error	"Incorrect input"
// @Failure		500			{object}	error	"Server error"
// @Router		/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie(commonHttp.AUTH_COOKIE_NAME)
	if err != nil {
		h.logger.Infof("Error while logout: %w", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Error while logout: %s", err.Error()))
		return
	}

	err = h.authServices.Logout(sessionID)
	if err != nil {
		h.logger.Errorf("Error while logout: %w", err)
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error while logout: %s", err.Error()))
		return
	}

	c.SetCookie(commonHttp.AUTH_COOKIE_NAME, sessionID, -1, "", "", false, false)
	c.Status(http.StatusOK)
}
