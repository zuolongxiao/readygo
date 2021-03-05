package routes

import (
	"github.com/gin-gonic/gin"
)

// Route Route
type Route struct {
	Method  string
	Pattern string
	Handler func(c *gin.Context)
	Desc    string
	Flag    string
}
