package web

import (
	"fmt"
	"net/http"
	"strings"
)

type errorCode string

const (
	UnauthorizedRequest errorCode = "unauthorized"
	BadRequest          errorCode = "bad_request"
	InternalServerError errorCode = "internal_server_error"
)

var (
	ErrBadRequest = func(desc string) ErrorInterface { return NewError(BadRequest, desc, "", http.StatusBadRequest) }
)

type ErrorInterface interface {
	WithCause(cause string) ErrorInterface
	Code() string
	Description() string
	HTTPStatusCode() int
	Error() string
	Cause() string
}

type ErrorFunc func(desc string) ErrorInterface

func NewError(errCode errorCode, desc string, cause string, httpCode int) ErrorInterface {
	return &customError{code: errCode, description: desc, httpStatusCode: httpCode, cause: cause}
}

type customError struct {
	code           errorCode
	description    string
	cause          string
	httpStatusCode int
}

func (e *customError) WithCause(cause string) ErrorInterface {
	if e.cause != "" {
		e.cause = strings.Join([]string{e.cause, cause}, ":")
	}

	return e
}

func (e *customError) Code() string {
	return string(e.code)
}

func (e *customError) Description() string {
	return e.description
}

func (e *customError) HTTPStatusCode() int {
	return e.httpStatusCode
}

func (e *customError) Error() string {
	return fmt.Sprintf("code: %s description: %s httpStatusCode: %d cause: %s",
		e.code,
		e.description,
		e.httpStatusCode,
		e.cause,
	)
}

func (e *customError) Cause() string {
	return e.cause
}
