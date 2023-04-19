package middlewares

import (
	"readygo/pkg/errs"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// TokenType Authorization token type
const TokenType = "Bearer"

var Whitelist = []string{
	"/api/v1/auth",
}

// Authenticate verify JWT
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.StrInSlice(c.Request.URL.Path, Whitelist) {
			c.Next()
			return
		}
		w := utils.NewContextWrapper(c)

		header := c.GetHeader("Authorization")
		if len(header) < len(TokenType)+1 {
			w.RespondAndAbort(errs.UnauthorizedError("missing token"), nil)
			return
		}

		token := header[len(TokenType)+1:]
		if token == "" {
			w.RespondAndAbort(errs.UnauthorizedError("empty token"), nil)
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			w.RespondAndAbort(errs.UnauthorizedError(err.Error()), nil)
			return
		}

		w.SetUsername(claims.Username)
		w.SetPermissions(claims.Permissions)
		c.Next()
	}
}
