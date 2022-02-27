package consts

// エンドポイント
// 下記シートの "Golang const生成マン" を使って生成
// https://docs.google.com/spreadsheets/d/1hp18YKpLHtalW98dBsYN1FZiobg-8EeES_U8othOZsU/edit#gid=0
// TODO 仕様書に書いた定義は全部コードに自動で落とし込みたい。websocket版swagger的なツールはないか?
const (
	// サーバへのリクエストエンドポイント
	TSJoinQuickMatch  = "join/quick_match"
	TSJoinFriendMatch = "join/friend_match"
	TSJoinRevive      = "join/revive"
	TSCreateMatch     = "create/match"
	TSGamePlayersInfo = "game/players_info"
	TSGameStart       = "game/start"
	TSGameNextTurn    = "game/next_turn"
	TSGameBid         = "game/bid"
	TSGameBuy         = "game/buy"
	TSGameCalculate   = "game/calculate"

	// サーバからのレスポンスエンドポイント
	FSGameJoin            = "game/join"
	FSGamePlayersInfo     = "game/players_info"
	FSGameNextTurn        = "game/next_turn"
	FSGameStart           = "game/start"
	FSGameBid             = "game/bid"
	FSGameBuyUpdate       = "game/buy_update"
	FSGameBuyNotify       = "game/buy_notify"
	FSGameCalculateResult = "game/calculate_result"
	FSGameCorrectPlayers  = "game/correct_players"
	FSGameUpdateAnswer    = "game/update_answer"
	FSGameFinishGame      = "game/finish_game"
	FSGameUpdateState     = "game/update_state"
)
