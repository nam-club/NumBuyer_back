package orgerrors

import "github.com/pkg/errors"

type InternalServerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *InternalServerError) Error() string { return e.Message }

func NewInternalServerError(message string) error {
	if message == "" {
		return errors.WithStack(&ValidationError{Code: "error.internal", Message: "internal server error"})
	} else {
		return errors.WithStack(&ValidationError{Code: "error.internal", Message: message})
	}
}
