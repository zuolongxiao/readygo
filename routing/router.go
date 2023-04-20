package routing

import (
	"reflect"
	"strings"

	"readygo/api"
	"readygo/middlewares"
	"readygo/pkg/settings"
	v1 "readygo/routing/routes/v1"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// NewRouters
func NewRouter() *gin.Engine {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	gin.SetMode(settings.Gin.Mode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	r.GET("/api", api.Index)

	v1Group := r.Group(v1.Prefix)
	v1Group.Use(v1.Middlewares...)
	for _, v := range v1.Routes {
		v1Group.Handle(v.Method, v.Pattern, v.Handler)
	}

	return r
}
