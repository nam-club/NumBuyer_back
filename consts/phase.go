package consts

type Phase string

const (
	PhaseBeforeStart = Phase("BEFORE_START")
	PhaseAuction     = Phase("AUCTION")
	PhaseCalculate   = Phase("CALCULATE")
	PhaseResult      = Phase("RESULT")
)
