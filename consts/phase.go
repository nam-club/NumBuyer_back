package consts

import "nam-club/NumBuyer_back/models/orgerrors"

type Phase struct {
	Value    string
	Duration int
}

var (
	PhaseBeforeStart     = Phase{"BEFORE_START", PhaseTimeValueInfinite}
	PhaseWating          = Phase{"WAITING", PhaseTimeValueInfinite}
	PhaseBeforeAuction   = Phase{"BEFORE_AUCTION", 14}
	PhaseAuction         = Phase{"AUCTION", 30}
	PhaseAuctionResult   = Phase{"AUCTION_RESULT", 5}
	PhaseCalculate       = Phase{"CALCULATE", 20}
	PhaseCalculateResult = Phase{"CALCULATE_RESULT", 5}
	PhaseEnd             = Phase{"END", PhaseTimeValueInfinite}
)

const (
	// 自動で部屋を閉じるまでの時間（秒）
	TimeAutoEnd = 120

	// -1なら無制限
	PhaseTimeValueInfinite = -1
)

func ParsePhase(s string) (v Phase, err error) {
	switch s {
	case PhaseBeforeStart.Value:
		return PhaseBeforeStart, nil
	case PhaseWating.Value:
		return PhaseWating, nil
	case PhaseBeforeAuction.Value:
		return PhaseBeforeAuction, nil
	case PhaseAuction.Value:
		return PhaseAuction, nil
	case PhaseAuctionResult.Value:
		return PhaseAuctionResult, nil
	case PhaseCalculate.Value:
		return PhaseCalculate, nil
	case PhaseCalculateResult.Value:
		return PhaseCalculateResult, nil
	case PhaseEnd.Value:
		return PhaseEnd, nil
	default:
		return PhaseBeforeStart, orgerrors.NewInternalServerError("invalid phase type")
	}
}
