package responses

import "nam-club/NumBuyer_back/consts"

type GameStartResponse struct {
	RoomID     string     `json:"roomId"`
	GoalCoin   int        `json:"goalCoin"`
	PhaseTimes PhaseTimes `json:"phaseTimes"`
}

type PhaseTimes struct {
	Ready           int `json:"ready"`
	GiveCards       int `json:"giveCards"`
	ShowTarget      int `json:"showTarget"`
	ShowAuction     int `json:"showAuction"`
	Auction         int `json:"auction"`
	AuctionResult   int `json:"auctionResult"`
	Calculate       int `json:"calculate"`
	CalculateResult int `json:"calculateResult"`
	NextTurn        int `json:"nextTurn"`
}

func GenerateGameStartResponse(roomId string, coinClearNum int) *GameStartResponse {
	return &GameStartResponse{RoomID: roomId, GoalCoin: coinClearNum, PhaseTimes: PhaseTimes{
		Ready:           consts.PhaseReady.Duration,
		GiveCards:       consts.PhaseGiveCards.Duration,
		ShowTarget:      consts.PhaseShowTarget.Duration,
		ShowAuction:     consts.PhaseShowAuction.Duration,
		Auction:         consts.PhaseAuction.Duration,
		AuctionResult:   consts.PhaseAuctionResult.Duration,
		Calculate:       consts.PhaseCalculate.Duration,
		CalculateResult: consts.PhaseCalculateResult.Duration,
		NextTurn:        consts.PhaseNextTurn.Duration,
	}}
}
