package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zuolongxiao/readygo/pkg/settings"
)

// Index index
func Index(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = settings.AppSetting.Name
	data["version"] = settings.AppSetting.Version

	c.JSON(http.StatusOK, data)
}
