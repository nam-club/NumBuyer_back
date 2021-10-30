package consts

// エンドポイント
// 下記シートの "Golang const生成マン" を使って生成
// https://docs.google.com/spreadsheets/d/1hp18YKpLHtalW98dBsYN1FZiobg-8EeES_U8othOZsU/edit#gid=0
// TODO 仕様書に書いた定義は全部コードに自動で落とし込みたい。websocket版swagger的なツールはないか?
const (
	// サーバへのリクエストエンドポイント
	ToServerJoinQuickMatch  = "join/quick_match"
	ToServerJoinFriendMatch = "join/friend_match"
	ToServerCreateMatch     = "create/match"
	ToServerGamePlayersInfo = "game/players_info"
	ToServerGameStart       = "game/start"
	ToServerGameNextTurn    = "game/next_turn"
	ToServerGameBuy         = "game/buy"
	ToServerGameAnswer      = "game/answer"

	// サーバからのレスポンスエンドポイント
	FromServerGameJoin                = "game/join"
	FromServerGamePlayersInfo         = "game/players_info"
	FromServerGameNextTurn            = "game/next_turn"
	FromServerGameStart               = "game/start"
	FromServerGameStartToClients      = "game/start_to_clients"
	FromServerGameJoinMember          = "game/join_member"
	FromServerGameBuyToClient         = "game/buy_to_client"
	FromServerGameAnswerToClient      = "game/answer_to_client"
	FromServerGameAnswerToClients     = "game/answer_to_clients"
	FromServerGameTargetCardToClients = "game/targetCard_to_clients"
	FromServerGameFinishGame          = "game/finish_game"
	FromServerGameUpdatePlayerInfo    = "game/update_playerInfo"
	FromServerGameUpdateState         = "game/update_state"
)
