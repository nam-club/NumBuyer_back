package orgerrors

import "github.com/pkg/errors"

type GameNotFoundError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *GameNotFoundError) Error() string { return e.Message }

func NewGameNotFoundError(message string) error {
	if message == "" {
		return errors.WithStack(&ValidationError{Code: "error.game.notFound", Message: "not found game"})
	} else {
		return errors.WithStack(&ValidationError{Code: "error.game.notFound", Message: message})
	}
}
