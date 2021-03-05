package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zuolongxiao/readygo/middlewares"
	"github.com/zuolongxiao/readygo/routing/routes"
)

// Prefix route prefix
const Prefix = "/api/v1"

// Middlewares route middlewares
var Middlewares = []gin.HandlerFunc{
	middlewares.Authenticate(),
	middlewares.Authorize(),
}

// Routes routes
var Routes = make([]routes.Route, 0, 64)
