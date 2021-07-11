package consts

// エンドポイント
// 下記シートの "Golang const生成マン" を使って生成
// https://docs.google.com/spreadsheets/d/1hp18YKpLHtalW98dBsYN1FZiobg-8EeES_U8othOZsU/edit#gid=0
// TODO 仕様書に書いた定義は全部コードに自動で落とし込みたい。websocket版swagger的なツールはないか?
const (
	// サーバへのリクエストエンドポイント
	ToServerJoinQuickMatch     = "join/quick_match"
	ToServerJoinFriendMatch    = "join/friend_match"
	ToServerGameStartToServer  = "game/start_to_server"
	ToServerGameBuyToServer    = "game/buy_to_server"
	ToServerGameAnswerToServer = "game/answer_to_server"
	// サーバからのレスポンスエンドポイント
	FromServerGameJoin                = "game/join"
	FromServerGameStartToClients      = "game/start_to_clients"
	FromServerGameJoinMember          = "game/join_member"
	FromServerGameBuyToClient         = "game/buy_to_client"
	FromServerGameAnswerToClient      = "game/answer_to_client"
	FromServerGameAnswerToClients     = "game/answer_to_clients"
	FromServerGameTargetCardToClients = "game/targetCard_to_clients"
	FromServerGameFinishGame          = "game/finish_game"
	FromServerGameUpdatePlayerInfo    = "game/update_playerInfo"
)
