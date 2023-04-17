package api

import (
	"net/http"
	"readygo/pkg/settings"

	"github.com/gin-gonic/gin"
)

// Index index
func Index(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = settings.App.Name
	data["version"] = settings.Version

	c.JSON(http.StatusOK, data)
}
