package http

import (
	"errors"
	"net/http"
	"os"

	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
	commonAuth "github.com/yarikTri/archipelago-notes-api/internal/common/http/auth"
	commonHttp "github.com/yarikTri/archipelago-notes-api/internal/common/http/constants"
	"github.com/yarikTri/archipelago-notes-api/internal/pkg/auth"
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

	sessionID, userID, expiration, err := h.authUsecase.SignUp(signUpInfo.Email, signUpInfo.Name, signUpInfo.Password)
	if err != nil {
		h.logger.Error(err.Error())
		var consistentError *pq.Error
		if errors.As(err, &consistentError) && consistentError.Code == "23505" {
			c.JSON(http.StatusBadRequest, "User already exists")
			return
		}
		c.JSON(http.StatusInternalServerError, "Error while sign up")
		return
	}

	c.SetCookie(commonHttp.SessionIdCookieName, sessionID, int(expiration.Seconds()), "", "", true, true)
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
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Not found user"})
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
	c.JSON(http.StatusOK, "OK")
}

// ClearAllSessions ..
func (h *Handler) ClearAllSessions(c *gin.Context) {
	adminPassword := c.GetHeader(commonHttp.AdminPasswordHeader)
	if adminPassword == "" {
		h.logger.Error("Admin password header is missing")
		c.JSON(http.StatusUnauthorized, "Admin password is required")
		return
	}

	if adminPassword != os.Getenv("ADMIN_PASSWORD") {
		h.logger.Error("Invalid admin password provided")
		c.JSON(http.StatusForbidden, "Invalid admin password")
		return
	}

	if err := h.authUsecase.ClearAllSessions(); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, "Error while clearing sessions")
		return
	}

	c.JSON(http.StatusOK, "All sessions cleared successfully")
}
