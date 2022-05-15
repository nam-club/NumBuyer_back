package orgerrors

import "github.com/pkg/errors"

type ValidationError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (e *ValidationError) Error() string { return e.Message }

// バリデーションエラーはエラーコードを指定できるようにする。
// codeSuffix: error.validationに続くエラーコード。指定しなくてもOK
func NewValidationError(codeSuffix string, message string, params map[string]string) error {
	status := 400
	var code string
	if codeSuffix != "" {
		code = "error.validation." + codeSuffix
	} else {
		code = "error.validation"
	}

	if message == "" {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: "validation error"})
	} else {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
