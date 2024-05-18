package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	commonHttp "github.com/yarikTri/archipelago-notes-api/internal/common/http/constants"
)

func SetUserIDMiddleware(resolveUserID func(sessionID string) (uuid.UUID, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, _ := c.Cookie(commonHttp.SessionIdCookieName)
		userID, _ := resolveUserID(sessionID)

		c.Request.Header.Set(commonHttp.UserIdHeader, userID.String())
		c.Next()
	}
}
