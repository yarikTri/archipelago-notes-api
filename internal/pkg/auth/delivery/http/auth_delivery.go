package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	commonAuth "github.com/yarikTri/archipelago-notes-api/internal/common/http/auth"
	commonHttp "github.com/yarikTri/archipelago-notes-api/internal/common/http/constants"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/auth"
	"net/http"
)

type Handler struct {
	authUsecase auth.Usecase
	logger      logger.Logger
}

func NewHandler(au auth.Usecase, l logger.Logger) *Handler {
	return &Handler{
		authUsecase: au,
		logger:      l,
	}
}

// CheckSession ..
func (h *Handler) CheckSession(c *gin.Context) {
	sessionID, err := commonAuth.GetSessionID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Forbidden")
		return
	}

	sessionUserID, err := h.authUsecase.GetUserIDBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Forbidden")
		return
	}

	if headerUserID, err := commonAuth.GetUserId(c); err != nil || sessionUserID != headerUserID {
		c.JSON(http.StatusUnauthorized, "Forbidden")
		return
	}

	c.JSON(http.StatusOK, "")
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserID string `json:"user_id"`
}

// SignUp ..
func (h *Handler) SignUp(c *gin.Context) {
	var signUpInfo SignUpRequest
	if err := c.BindJSON(&signUpInfo); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid sign up data")
		return
	}

	if len(signUpInfo.Password) > 72 {
		c.JSON(http.StatusBadRequest, "Too long password")
		return
	}

	userID, err := h.authUsecase.SignUp(signUpInfo.Email, signUpInfo.Email, signUpInfo.Password)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "Error while sign up")
		return
	}

	c.JSON(http.StatusOK, SignUpResponse{UserID: userID.String()})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
}

// Login ..
func (h *Handler) Login(c *gin.Context) {
	var credentials LoginRequest
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid sign up data")
		return
	}

	sessionID, userID, expiration, err := h.authUsecase.Login(credentials.Email, credentials.Password)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "Error while login")
		return
	}

	c.SetCookie(commonHttp.SessionIdCookieName, sessionID, int(expiration.Seconds()), "", "", true, true)
	c.JSON(http.StatusOK, LoginResponse{UserID: userID.String()})
}

// Logout ..
func (h *Handler) Logout(c *gin.Context) {
	sessionID, err := commonAuth.GetSessionID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Forbidden")
		return
	}

	if err := h.authUsecase.Logout(sessionID); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "Error while logout")
		return
	}

	c.SetCookie(commonHttp.SessionIdCookieName, "", -1, "", "", true, true)
}
