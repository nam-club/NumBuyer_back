package orgerrors

import "github.com/pkg/errors"

type AbilityNotFoundError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *AbilityNotFoundError) Error() string { return e.Message }

func NewAbilityNotFoundError(message string, params map[string]string) error {
	status := 404
	code := "error.ability.notFound"

	if message == "" {
		return errors.WithStack(&AbilityNotFoundError{Status: status, Code: code, Message: "ability not found"})
	} else {
		return errors.WithStack(&AbilityNotFoundError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
