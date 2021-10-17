package logic

import (
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"

	"github.com/google/uuid"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName, roomId string, isOwner bool) (*responses.Player, error) {

	if b, _ := db.ExistsGame(roomId); b == false {
		return nil, orgerrors.NewGameNotFoundError("game not found")
	}

	p := db.Player{
		PlayerID:   generatePlayerId(roomId),
		PlayerName: playerName,
		Coin:       100,
	}

	var regist db.Player
	var e error
	regist, e = db.AddPlayer(roomId, p)
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

// プレイヤーIDを生成する
func generatePlayerId(roomId string) string {
	return uuid.Must(uuid.NewUUID()).String()

}
