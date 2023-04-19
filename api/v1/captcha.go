package api

import (
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var captchaStore = base64Captcha.DefaultMemStore

// GenerateCaptcha
func GenerateCaptcha(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	driver := base64Captcha.DriverDigit{
		Height: 40,
		Width:  100,
		Length: 6,
	}
	cap := base64Captcha.NewCaptcha(&driver, captchaStore)
	id, b64s, err := cap.Generate()
	if err != nil {
		cw.Respond(err, nil)
		return
	}
	data := map[string]interface{}{
		"data": b64s,
		"id":   id,
	}
	cw.Respond(nil, data)
}
