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
	Grace     int // 計算フェーズなどでフロントより多く持たせる時間
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
	PhaseWaiting = Phase{"WAITING", PhaseTimeValueInfinite, 0, &PhaseReady}
	PhaseReady = Phase{"READY", 2, 0, &PhaseGiveCards}
	PhaseGiveCards = Phase{"GIVE_CARDS", 3, 0, &PhaseShowTarget}
	PhaseShowTarget = Phase{"SHOW_TARGET", 3, 0, &PhaseShowAuction}
	PhaseShowAuction = Phase{"SHOW_AUCTION", 3, 0, &PhaseAuction}
	PhaseAuction = Phase{"AUCTION", 15, 1, &PhaseAuctionResult}
	PhaseAuctionResult = Phase{"AUCTION_RESULT", 5, 0, &PhaseCalculate}
	PhaseCalculate = Phase{"CALCULATE", 20, 1, &PhaseCalculateResult}
	PhaseCalculateResult = Phase{"CALCULATE_RESULT", 5, 0, &PhaseNextTurn}
	PhaseNextTurn = Phase{"NEXT_TURN", 2, 0, &PhaseReady}
	PhaseEnd = Phase{"END", PhaseTimeValueInfinite, 0, nil}
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
