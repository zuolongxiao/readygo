package api

import (
	"net/http"

	"readygo/pkg/settings"

	"github.com/gin-gonic/gin"
)

// Index index
func Index(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = settings.AppSetting.Name
	data["version"] = settings.AppSetting.Version

	c.JSON(http.StatusOK, data)
}
