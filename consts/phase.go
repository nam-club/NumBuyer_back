package consts

import "nam-club/NumBuyer_back/models/orgerrors"

type Phase string

const (
	PhaseBeforeStart = Phase("BEFORE_START")
	PhaseAuction     = Phase("AUCTION")
	PhaseCalculate   = Phase("CALCULATE")
	PhaseResult      = Phase("RESULT")
	PhaseEnd         = Phase("END")
)

func (v Phase) Valid() error {
	switch v {
	case PhaseBeforeStart, PhaseAuction, PhaseCalculate, PhaseResult, PhaseEnd:
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
