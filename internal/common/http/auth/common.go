package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	commonHttp "github.com/yarikTri/archipelago-notes-api/internal/common/http/constants"
)

func GetSessionID(c *gin.Context) (string, error) {
	return c.Cookie(commonHttp.SessionIdCookieName)
}

func GetUserId(c *gin.Context) (uuid.UUID, error) {
	return uuid.FromString(c.GetHeader(commonHttp.UserIdHeader))
}
