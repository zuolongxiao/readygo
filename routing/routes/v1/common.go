package v1

import (
	"readygo/middlewares"
	"readygo/routing/routes"

	"github.com/gin-gonic/gin"
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
