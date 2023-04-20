package api

import (
	"net/http"
	"readygo/pkg/global"

	"github.com/gin-gonic/gin"
)

// Index index
func Index(c *gin.Context) {
	data := map[string]string{
		"name":    global.Name,
		"version": global.Version,
		"author":  global.Author,
		"email":   global.Email,
	}

	c.JSON(http.StatusOK, data)
}
