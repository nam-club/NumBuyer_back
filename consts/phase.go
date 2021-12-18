package consts

import "nam-club/NumBuyer_back/models/orgerrors"

const (
	// 自動で部屋を閉じるまでの時間（秒）
	TimeAutoEnd = 120

	// -1なら無制限
	PhaseTimeValueInfinite = -1
)

type Phase struct {
	Value     string
	Duration  int
	NextPhase *Phase
}

var (
	PhaseWaiting         Phase
	PhaseReady           Phase
	PhaseGiveCards       Phase
	PhaseShowTarget      Phase
	PhaseShowAuction     Phase
	PhaseAuction         Phase
	PhaseAuctionResult   Phase
	PhaseCalculate       Phase
	PhaseCalculateResult Phase
	PhaseNextTurn        Phase
	PhaseEnd             Phase
)

func init() {
	PhaseWaiting = Phase{"WAITING", PhaseTimeValueInfinite, &PhaseReady}
	PhaseReady = Phase{"READY", 2, &PhaseGiveCards}
	PhaseGiveCards = Phase{"GIVE_CARDS", 3, &PhaseShowTarget}
	PhaseShowTarget = Phase{"SHOW_TARGET", 3, &PhaseShowAuction}
	PhaseShowAuction = Phase{"SHOW_AUCTION", 3, &PhaseAuction}
	PhaseAuction = Phase{"AUCTION", 30, &PhaseAuctionResult}
	PhaseAuctionResult = Phase{"AUCTION_RESULT", 5, &PhaseCalculate}
	PhaseCalculate = Phase{"CALCULATE", 20, &PhaseCalculateResult}
	PhaseCalculateResult = Phase{"CALCULATE_RESULT", 5, &PhaseNextTurn}
	PhaseNextTurn = Phase{"NEXT_TURN", 2, &PhaseReady}
	PhaseEnd = Phase{"END", PhaseTimeValueInfinite, nil}
}
func ParsePhase(s string) (v Phase, err error) {
	switch s {
	case PhaseWaiting.Value:
		return PhaseWaiting, nil
	case PhaseReady.Value:
		return PhaseReady, nil
	case PhaseGiveCards.Value:
		return PhaseGiveCards, nil
	case PhaseShowTarget.Value:
		return PhaseShowTarget, nil
	case PhaseShowAuction.Value:
		return PhaseShowAuction, nil
	case PhaseAuction.Value:
		return PhaseAuction, nil
	case PhaseAuctionResult.Value:
		return PhaseAuctionResult, nil
	case PhaseCalculate.Value:
		return PhaseCalculate, nil
	case PhaseCalculateResult.Value:
		return PhaseCalculateResult, nil
	case PhaseNextTurn.Value:
		return PhaseNextTurn, nil
	case PhaseEnd.Value:
		return PhaseEnd, nil
	default:
		return PhaseWaiting, orgerrors.NewInternalServerError("invalid phase type")
	}
}
