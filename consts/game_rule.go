package consts

const (
	// プレイヤーの初期カード枚数
	InitialCardsNum = 10

	// コイン数のクリア条件
	CoinClearNum = 100

	// 初期コイン数
	InitialCoin = 30

	// オークンションカードで符号が出る確率（パーセント）
	AuctionCodeProbability = 35

	// 数字の最小値(以上)、最小値(未満)
	TermMin = 1
	TermMax = 2

	// ターゲットカードの最小値(以上)、最小値(未満)
	TargetMin = 1
	TargetMax = 2

	// ランダムに生成されるオークションカードの枚数
	AuctionCardsNumMin = 1
	AuctionCardsNumMax = 3

	// 符号
	CodePlus   = "+"
	CodeMinus  = "-"
	CodeTimes  = "*"
	CodeDivide = "/"

	// オークションの時間をリセットした時に残す時間
	AuctionResetTimeRemains = 10

	// プレイヤーの最大入札回数
	AuctionMaxBidCount = 10

	// プレイヤー数の最小・最大
	QuickMatchPlayersMin  = 4
	QuickMatchPlayersMax  = 4
	FriendMatchPlayersMin = 1
	FriendMatchPlayersMax = 6
)

var (
	// 符号のslice
	// 乗算、除算カードを生成したい場合はコメントアウトを代わりに使う
	//	Codes = []string{CodePlus, CodeMinus, CodeTimes, CodeDivide}
	Codes = []string{CodePlus, CodeMinus}
)
