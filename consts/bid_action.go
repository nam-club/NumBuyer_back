package consts

import (
	"nam-club/NumBuyer_back/models/orgerrors"
	"strings"
)

type BidAction string

const (
	BidActionBid  BidAction = "bid"
	BidActionPass BidAction = "pass"
)

func (v BidAction) String() string {
	return string(v)
}

func (v BidAction) Valid() error {
	switch v {
	case BidActionBid, BidActionPass:
		return nil
	default:
		return orgerrors.NewValidationError("bid,action", "invalid bid action type", map[string]string{"action": v.String()})
	}
}

func (v *BidAction) UnmarshalJSON(b []byte) error {
	*v = BidAction(strings.Trim(string(b), `"`))
	return v.Valid()
}

func ParseBidAction(s string) (v BidAction, err error) {
	v = BidAction(s)
	err = v.Valid()
	return
}
