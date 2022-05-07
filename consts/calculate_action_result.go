package consts

import (
	"nam-club/NumBuyer_back/models/orgerrors"
	"strings"
)

type CalculateActionResult string

const (
	CalculateActionResultCorrect   CalculateActionResult = "correct"
	CalculateActionResultIncorrect CalculateActionResult = "incorrect"
	CalculateActionResultPass      CalculateActionResult = "pass"
)

func (v CalculateActionResult) String() string {
	return string(v)
}

func (v CalculateActionResult) Valid() error {
	switch v {
	case CalculateActionResultCorrect, CalculateActionResultIncorrect, CalculateActionResultPass:
		return nil
	default:
		return orgerrors.NewValidationError("invalid bid action type")
	}
}

func (v *CalculateActionResult) UnmarshalJSON(b []byte) error {
	*v = CalculateActionResult(strings.Trim(string(b), `"`))
	return v.Valid()
}

func ParseCalculateActionResult(s string) (v CalculateActionResult, err error) {
	v = CalculateActionResult(s)
	err = v.Valid()
	return
}
