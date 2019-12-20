package errors

import (
	errors_extra "errors"
	"fmt"
)

const (
	CodeInternalServerError = 500
	CodeBadParameter        = 600
	CodeDBError             = 501
	CodeItemAlreadyExist    = 401
)

type APIError struct {
	code int
	err  error
}

func (h *APIError) Error() string {
	return h.err.Error()
}

func (h *APIError) Code() int {
	return h.code
}

func NewAPIError(code int, message string, params ...interface{}) *APIError {
	return &APIError{
		code: code,
		err:  errors_extra.New(fmt.Sprintf(message, params...)),
	}
}

func NewAPIErrorWithError(code int, err error) *APIError {
	return &APIError{
		code: code,
		err:  err,
	}
}
