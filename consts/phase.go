package consts

import "nam-club/NumBuyer_back/models/orgerrors"

type Phase string

const (
	PhaseBeforeStart     = Phase("BEFORE_START")
	PhaseBeforeAuction   = Phase("BEFORE_AUCTION")
	PhaseAuction         = Phase("AUCTION")
	PhaseAuctionResult   = Phase("AUCTION_RESULT")
	PhaseCalculate       = Phase("CALCULATE")
	PhaseCalculateResult = Phase("CALCULATE_RESULT")
	PhaseEnd             = Phase("END")
)

func (v Phase) Valid() error {
	switch v {
	case PhaseBeforeStart, PhaseBeforeAuction, PhaseAuction, PhaseCalculate, PhaseCalculateResult, PhaseEnd:
		return nil
	default:
		return orgerrors.NewInternalServerError("invalid phase type")
	}
}

func ParsePhase(s string) (v Phase, err error) {
	v = Phase(s)
	err = v.Valid()
	return
}
