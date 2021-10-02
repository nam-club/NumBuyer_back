package logic

import (
	"errors"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/responses"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName, roomId string, isOwner bool) (*responses.Player, error) {

	if b, _ := db.ExistsGame(roomId); b == false {
		return nil, errors.New("invalid game id")
	}

	p := db.Player{
		PlayerID:   generatePlayerId(roomId),
		PlayerName: playerName,
		Coin:       100,
	}

	var regist db.Player
	var e error
	regist, e = db.SetPlayer(roomId, p)
	if e != nil {
		return nil, e
	}

	ret := &responses.Player{
		PlayerID:   regist.PlayerID,
		PlayerName: regist.PlayerName,
		RoomID:     roomId,
		Money:      regist.Coin,
		Cards:      regist.Cards,
	}

	return ret, nil
}

// ゲームIDを生成する
func generatePlayerId(roomId string) int {
	game, _ := db.GetGame(roomId)
	return len(game.Players) + 1
}
