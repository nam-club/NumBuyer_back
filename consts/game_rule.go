package consts

const (
	// コイン数のクリア条件
	CoinClearNum = 3000

	// オークンションカードで符号が出る確率（パーセント）
	AuctionCodeProbability = 33

	// 数字の最小値(以上)、最小値(未満)
	TermMin = 1
	TermMax = 100

	// 符号
	CodePlus   = "+"
	CodeMinus  = "-"
	CodeTimes  = "*"
	CodeDivide = "/"
)

var (
	// 符号のslice
	Codes = []string{CodePlus, CodeMinus, CodeTimes, CodeDivide}
)
