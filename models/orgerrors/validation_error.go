package orgerrors

import "github.com/pkg/errors"

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string { return e.Message }

func NewValidationError(message string) error {
	if message == "" {
		return errors.WithStack(&ValidationError{Code: "error.validation", Message: "validation error"})
	} else {
		return errors.WithStack(&ValidationError{Code: "error.validation", Message: message})
	}
}
