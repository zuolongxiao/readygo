package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zuolongxiao/readygo/pkg/errs"
)

const (
	usernameKey    = "username"
	permissionsKey = "permissions"
	successCode    = 0
	successMessage = "OK"
)

// ContextWrapper gin context wrapper
type ContextWrapper struct {
	ctx *gin.Context
}

// NewContextWrapper new context wrapper
func NewContextWrapper(c *gin.Context) *ContextWrapper {
	return &ContextWrapper{ctx: c}
}

type response struct {
	Code    int         `json:"errcode"`
	Message string      `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"`
}

func (w *ContextWrapper) respond(err error, data interface{}, abort bool) {
	if err == nil {
		w.ctx.JSON(http.StatusOK, response{
			Code:    successCode,
			Message: successMessage,
			Data:    data,
		})

		if abort {
			w.ctx.Abort()
		}

		return
	}

	errorCode := errs.BadRequestErrorCode
	statusCode := http.StatusBadRequest
	if e, ok := err.(errs.AppError); ok {
		errorCode = e.ErrorCode()
		statusCode = e.StatusCode()
	}

	w.ctx.JSON(statusCode, response{
		Code:    errorCode,
		Message: err.Error(),
		Data:    data,
	})

	if abort {
		w.ctx.Abort()
	}
}

// Respond send response
func (w *ContextWrapper) Respond(err error, data interface{}) {
	w.respond(err, data, false)
}

// RespondAndAbort send response and abort
func (w *ContextWrapper) RespondAndAbort(err error, data interface{}) {
	w.respond(err, data, true)
}

// SetUsername set username to context
func (w *ContextWrapper) SetUsername(username string) {
	w.ctx.Set(usernameKey, username)
}

// GetUsername get username from context
func (w *ContextWrapper) GetUsername() string {
	return w.ctx.GetString(usernameKey)
}

// SetPermissions set permissions to context
func (w *ContextWrapper) SetPermissions(permissions []string) {
	w.ctx.Set(permissionsKey, permissions)
}

// GetPermissions get permissions from context
func (w *ContextWrapper) GetPermissions() []string {
	return w.ctx.GetStringSlice(permissionsKey)
}

// Bind bind and validate
func (w *ContextWrapper) Bind(o interface{}) error {
	formats := map[string]string{
		"required": "%s is required",
		"oneof":    "%s must be one of %s",
		"alphanum": "%s must be alphanum",
		"min":      "%s must be greater than %s",
		"max":      "%s must be less than %s",
		"eqfield":  "%s is not matched",
	}

	if err := w.ctx.ShouldBindJSON(o); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			for _, f := range verr {
				var errmsg string
				tag := f.ActualTag()
				field := f.Field()
				param := f.Param()

				switch tag {
				case "required":
					errmsg = fmt.Sprintf(formats[tag], field)
				case "oneof":
					errmsg = fmt.Sprintf(formats[tag], field, param)
				case "alphanum":
					errmsg = fmt.Sprintf(formats[tag], field)
				case "min":
					errmsg = fmt.Sprintf(formats[tag], field, param)
				case "max":
					errmsg = fmt.Sprintf(formats[tag], field, param)
				case "eqfield":
					errmsg = fmt.Sprintf(formats[tag], field)
				default:
					errmsg = fmt.Sprintf("%s is invalid", field)
				}

				return errs.ValidationError(errmsg)
			}
		} else {
			return errs.BadRequestError("")
		}
	}

	return nil
}
