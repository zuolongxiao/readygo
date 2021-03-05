package errs

import "net/http"

// Error Codes
const (
	BadRequestErrorCode        = 10000
	ValidationErrorCode        = 10010
	DuplicatedErrorCode        = 10020
	UnauthorizedErrorCode      = 10030
	ForbiddenErrorCode         = 10040
	LockedErrorCode            = 10050
	NotFoundErrorCode          = 10060
	InternalServerErrorCode    = 10070
	DBErrorCode                = 10080
	ReferenceRestrictErrorCode = 10090
)

// AppError app error
type AppError interface {
	StatusCode() int
	ErrorCode() int
}

// BadRequestError BadRequestError
type BadRequestError string

func (e BadRequestError) Error() string {
	return "Bad request"
}

// StatusCode StatusCode
func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

// ErrorCode ErrorCode
func (e BadRequestError) ErrorCode() int {
	return BadRequestErrorCode
}

// ValidationError ValidationError
type ValidationError string

func (e ValidationError) Error() string {
	return "Validation error: " + string(e)
}

// StatusCode StatusCode
func (e ValidationError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

// ErrorCode ErrorCode
func (e ValidationError) ErrorCode() int {
	return ValidationErrorCode
}

// DuplicatedError DuplicatedError
type DuplicatedError string

func (e DuplicatedError) Error() string {
	return "Resource is duplicated: " + string(e)
}

// StatusCode StatusCode
func (e DuplicatedError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

// ErrorCode ErrorCode
func (e DuplicatedError) ErrorCode() int {
	return DuplicatedErrorCode
}

// UnauthorizedError UnauthorizedError
type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return "Unauthorized: " + string(e)
}

// StatusCode implement StatusCode
func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

// ErrorCode ErrorCode
func (e UnauthorizedError) ErrorCode() int {
	return UnauthorizedErrorCode
}

// ForbiddenError ForbiddenError
type ForbiddenError string

func (e ForbiddenError) Error() string {
	return "Forbidden: " + string(e)
}

// StatusCode implement StatusCode
func (e ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}

// ErrorCode ErrorCode
func (e ForbiddenError) ErrorCode() int {
	return ForbiddenErrorCode
}

// LockedError LockedError
type LockedError string

func (e LockedError) Error() string {
	return "User locked: " + string(e)
}

// StatusCode implement StatusCode
func (e LockedError) StatusCode() int {
	return http.StatusLocked
}

// ErrorCode ErrorCode
func (e LockedError) ErrorCode() int {
	return LockedErrorCode
}

// NotFoundError NotFoundError
type NotFoundError string

func (e NotFoundError) Error() string {
	return "Resource does not exist: " + string(e)
}

// StatusCode implement StatusCode
func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// ErrorCode ErrorCode
func (e NotFoundError) ErrorCode() int {
	return NotFoundErrorCode
}

// InternalServerError InternalServerError
type InternalServerError string

func (e InternalServerError) Error() string {
	return "Server error: " + string(e)
}

// StatusCode StatusCode
func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

// ErrorCode ErrorCode
func (e InternalServerError) ErrorCode() int {
	return InternalServerErrorCode
}

// DBError DBError
type DBError string

func (e DBError) Error() string {
	return "DB error: " + string(e)
}

// StatusCode  StatusCode
func (e DBError) StatusCode() int {
	return http.StatusInternalServerError
}

// ErrorCode ErrorCode
func (e DBError) ErrorCode() int {
	return DBErrorCode
}

// ReferenceRestrictError ReferenceRestrictError
type ReferenceRestrictError string

func (e ReferenceRestrictError) Error() string {
	return "Reference restricted: " + string(e)
}

// StatusCode StatusCode
func (e ReferenceRestrictError) StatusCode() int {
	return http.StatusUnprocessableEntity
}

// ErrorCode ErrorCode
func (e ReferenceRestrictError) ErrorCode() int {
	return ReferenceRestrictErrorCode
}
