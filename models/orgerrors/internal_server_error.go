package orgerrors

import "github.com/pkg/errors"

type InternalServerError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *InternalServerError) Error() string { return e.Message }

func NewInternalServerError(message string, params map[string]string) error {
	status := 500
	code := "error.game.notFound"

	if message == "" {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: "internal server error"})
	} else {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
