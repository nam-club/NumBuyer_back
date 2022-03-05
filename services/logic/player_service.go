package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"

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
		IsOwner:    isOwner,
		Coin:       consts.InitialCoin,
	}

	var ret *db.Player
	var e error
	ret, e = db.SetPlayer(roomId, p)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

// プレイヤー情報を取得する
func GetPlayersInfo(roomId string) (*responses.PlayersInfoResponse, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}

	return responses.GeneratePlayersInfoResponse(players, roomId), nil
}

// 全プレイヤーが次フェーズに移行する準備ができている状態にする
func SetAllPlayersReady(roomId string) error {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return e
	}

	for _, p := range players {
		p.Ready = true
		db.SetPlayer(roomId, &p)
	}

	return nil
}

// 全プレイヤーにランダムにカードを一枚付与する
func AddCardToAllPlayers(roomId string) error {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return e
	}

	for _, p := range players {
		p.Cards = append(p.Cards, utils.GenerateRandomCard(1)[0])
		db.SetPlayer(roomId, &p)
	}

	return nil
}

// 全プレイヤーが次フェーズに移行する準備ができているか
func IsAllPlayersReady(roomId string) (bool, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return false, e
	}

	ready := true
	for _, p := range players {
		if !p.Ready {
			ready = false
			break
		}
	}

	return ready, nil
}

// プレイヤーにカードを追加する
func AppendCard(roomId, playerId string, appendCards []string) (*db.Player, error) {
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	player.Cards = append(player.Cards, appendCards...)
	player, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}

	return player, nil
}

// プレイヤーのコインを減らす
func SubtractCoin(roomId, playerId string, subtract int) (*db.Player, error) {
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	subtracted := player.Coin - subtract
	if subtracted < 0 {
		return nil, orgerrors.NewValidationError("coin shortage")
	}

	player.Coin = subtracted
	player, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}

	return player, nil
}

// プレイヤーIDを生成する
func generatePlayerId(roomId string) string {
	return uuid.Must(uuid.NewUUID()).String()
}
