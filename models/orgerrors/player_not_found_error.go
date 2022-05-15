package orgerrors

import "github.com/pkg/errors"

type PlayerNotFoundError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *PlayerNotFoundError) Error() string { return e.Message }

func NewPlayerNotFoundError(message string, params map[string]string) error {
	status := 404
	code := "error.player.notFound"

	if message == "" {
		return errors.WithStack(&PlayerNotFoundError{Status: status, Code: code, Message: "not found game"})
	} else {
		return errors.WithStack(&PlayerNotFoundError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
