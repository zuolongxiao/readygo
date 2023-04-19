package api

import (
	"net/http"
	"readygo/pkg/settings"

	"github.com/gin-gonic/gin"
)

// Index index
func Index(c *gin.Context) {
	data := map[string]string{
		"name":    settings.App.Name,
		"version": settings.Version,
	}

	c.JSON(http.StatusOK, data)
}
