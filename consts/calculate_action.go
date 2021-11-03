package consts

import (
	"nam-club/NumBuyer_back/models/orgerrors"
	"strings"
)

type CalculateAction string

const (
	CalculateActionBid  CalculateAction = "answer"
	CalculateActionPass                 = "pass"
)

func (v CalculateAction) String() string {
	return string(v)
}

func (v CalculateAction) Valid() error {
	switch v {
	case CalculateActionBid, CalculateActionPass:
		return nil
	default:
		return orgerrors.NewValidationError("invalid bid action type")
	}
}

func (v *CalculateAction) UnmarshalJSON(b []byte) error {
	*v = CalculateAction(strings.Trim(string(b), `"`))
	return v.Valid()
}

func ParseCalculateAction(s string) (v CalculateAction, err error) {
	v = CalculateAction(s)
	err = v.Valid()
	return
}
