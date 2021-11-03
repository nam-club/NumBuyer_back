package orgerrors

import "github.com/pkg/errors"

type PlayerNotFoundError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *PlayerNotFoundError) Error() string { return e.Message }

func NewPlayerNotFoundError(message string) error {
	if message == "" {
		return errors.WithStack(&PlayerNotFoundError{Code: "error.player.notFound", Message: "not found game"})
	} else {
		return errors.WithStack(&PlayerNotFoundError{Code: "error.player.notFound", Message: message})
	}
}
