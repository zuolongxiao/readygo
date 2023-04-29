package utils

import (
	"errors"
	"fmt"
	"net/http"

	"readygo/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func (cw *ContextWrapper) respond(err error, data interface{}, abort bool) {
	if err == nil {
		cw.ctx.JSON(http.StatusOK, response{
			Code:    successCode,
			Message: successMessage,
			Data:    data,
		})

		if abort {
			cw.ctx.Abort()
		}

		return
	}

	errorCode := errs.BadRequestErrorCode
	statusCode := http.StatusBadRequest
	if e, ok := err.(errs.AppError); ok {
		errorCode = e.ErrorCode()
		statusCode = e.StatusCode()
	}

	cw.ctx.JSON(statusCode, response{
		Code:    errorCode,
		Message: err.Error(),
		Data:    data,
	})

	if abort {
		cw.ctx.Abort()
	}
}

// Respond send response
func (cw *ContextWrapper) Respond(err error, data interface{}) {
	cw.respond(err, data, false)
}

// RespondAndAbort send response and abort
func (cw *ContextWrapper) RespondAndAbort(err error, data interface{}) {
	cw.respond(err, data, true)
}

// SetUsername set username to context
func (cw *ContextWrapper) SetUsername(username string) {
	cw.ctx.Set(usernameKey, username)
}

// GetUsername get username from context
func (cw *ContextWrapper) GetUsername() string {
	return cw.ctx.GetString(usernameKey)
}

// SetPermissions set permissions to context
func (cw *ContextWrapper) SetPermissions(permissions []string) {
	cw.ctx.Set(permissionsKey, permissions)
}

// GetPermissions get permissions from context
func (cw *ContextWrapper) GetPermissions() []string {
	return cw.ctx.GetStringSlice(permissionsKey)
}

// Bind bind and validate
func (cw *ContextWrapper) Bind(o interface{}) error {
	formats := map[string]string{
		"required": "%s is required",
		"oneof":    "%s must be one of %s",
		"alphanum": "%s must be alphanum",
		"min":      "%s must be greater than %s",
		"max":      "%s must be less than %s",
		"eqfield":  "%s is not matched",
	}

	if err := cw.ctx.ShouldBindJSON(o); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			var errmsg string
			for i, f := range validationErr {
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
				if i == 0 {
					break
				}
			}
			return errs.ValidationError(errmsg)
		} else {
			return errs.BadRequestError("")
		}
	}

	return nil
}
