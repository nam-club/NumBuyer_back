package orgerrors

import "github.com/pkg/errors"

type MutexError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *MutexError) Error() string { return e.Message }

func NewMutexError(message string, params map[string]string) error {
	status := 500
	code := "error.mutex"

	if message == "" {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: "mutex error"})
	} else {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
