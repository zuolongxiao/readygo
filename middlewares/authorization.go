package middlewares

import (
	"strings"

	"readygo/pkg/errs"
	"readygo/pkg/jobs"
	"readygo/utils"

	"github.com/gin-gonic/gin"
)

// Authorize check permissions
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		cw := utils.NewContextWrapper(c)

		handler := c.HandlerName()
		tmp := strings.Split(handler, ".")
		name := tmp[len(tmp)-1:][0]

		perms := jobs.GetPermissions()
		if !utils.StrInSlice(name, perms) {
			c.Next()
			return
		}

		ps := cw.GetPermissions()
		if utils.StrInSlice("*", ps) {
			c.Next()
			return
		}

		if utils.StrInSlice(name, ps) {
			c.Next()
			return
		}

		cw.RespondAndAbort(errs.ForbiddenError(name), nil)
	}
}
