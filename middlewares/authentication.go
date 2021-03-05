package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zuolongxiao/readygo/pkg/errs"
	"github.com/zuolongxiao/readygo/pkg/utils"
)

// TokenType Authorization token type
const TokenType = "Bearer"

// Authenticate verify JWT
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
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
