package orgerrors

import "github.com/pkg/errors"

type GameNotFoundError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *GameNotFoundError) Error() string { return e.Message }

func NewGameNotFoundError(message string, params map[string]string) error {
	status := 404
	code := "error.game.notFound"

	if message == "" {
		return errors.WithStack(&GameNotFoundError{Status: status, Code: code, Message: "not found game"})
	} else {
		return errors.WithStack(&GameNotFoundError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
