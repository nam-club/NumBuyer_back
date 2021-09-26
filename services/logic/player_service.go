package logic

import (
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/responses"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName, roomId string, isOwner bool) *responses.Player {

	p := db.Player{
		PlayerID:   1,
		PlayerName: playerName,
		Coin:       100,
	}

	db.SetPlayer(roomId, p)

	// TODO 要定義
	return &responses.Player{}
}
