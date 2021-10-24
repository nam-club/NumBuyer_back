package consts

const (
	// 各フェーズの持ち時間（秒）
	PhaseTimeBeforeStart     = PhaseTimeValueInfinite
	PhaseTimeBeforeAuction   = 14
	PhaseTimeAuction         = 30
	PhaseTimeAuctionResult   = 5
	PhaseTimeCalculate       = 20
	PhaseTimeCalculateResult = 5
	PhaseTimeEnd             = PhaseTimeValueInfinite

	// -1なら無制限
	PhaseTimeValueInfinite = -1
)
