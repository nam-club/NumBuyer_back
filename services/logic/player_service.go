package logic

import (
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"

	"github.com/google/uuid"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName, roomId string, isOwner bool) (*db.Player, error) {

	if b, _ := db.ExistsGame(roomId); b == false {
		return nil, orgerrors.NewGameNotFoundError("game not found")
	}

	p := &db.Player{
		PlayerID:   generatePlayerId(roomId),
		PlayerName: playerName,
		Coin:       100,
	}

	var ret *db.Player
	var e error
	ret, e = db.AddPlayer(roomId, p)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

// プレイヤー情報を取得する
func GetPlayersInfo(roomId, playerId string) (*responses.PlayersInfoResponse, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}

	return responses.GeneratePlayersInfoResponse(players, roomId), nil
}

// プレイヤーIDを生成する
func generatePlayerId(roomId string) string {
	return uuid.Must(uuid.NewUUID()).String()
}
