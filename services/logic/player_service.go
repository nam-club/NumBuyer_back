package logic

import (
	"nam-club/NumBuyer_back/models/responses"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName string, roomName string) *responses.Player {
	return &responses.Player{
		PlayerID:   1,
		PlayerName: playerName,
		RoomName:   roomName,
		Money:      100,
		// Cards: new {}
	}
}
