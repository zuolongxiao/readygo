package api

import (
	"readygo/pkg/errs"
	"readygo/pkg/settings"
	"readygo/pkg/store"
	"readygo/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// GenerateCaptcha
func GenerateCaptcha(c *gin.Context) {
	cw := utils.NewContextWrapper(c)

	if !settings.Captcha.Enabled {
		cw.Respond(errs.NotFoundError("captcha"), nil)
		return
	}

	driver := base64Captcha.DriverDigit{
		Height: settings.Captcha.Height,
		Width:  settings.Captcha.Width,
		Length: settings.Captcha.Length,
	}
	cap := base64Captcha.NewCaptcha(&driver, store.CaptchaStore)
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
