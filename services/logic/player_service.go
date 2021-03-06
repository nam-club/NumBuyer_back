package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"

	"github.com/rs/xid"
)

// 新規プレイヤー情報を生成する
func CreateNewPlayer(playerName, roomId string, isOwner bool, abilities []consts.Ability) (*db.Player, error) {

	if b, _ := db.ExistsGame(roomId); !b {
		return nil, orgerrors.NewPlayerNotFoundError("player not found", nil)
	}

	dbAbilities := []db.Ability{}
	for _, v := range abilities {
		dbAbilities = append(dbAbilities, db.Ability{ID: v.ID, Remaining: v.UsableNum, Status: string(v.InitialStatus)})
	}

	p := &db.Player{
		PlayerID:   generatePlayerId(roomId),
		PlayerName: playerName,
		IsOwner:    isOwner,
		Abilities:  dbAbilities,
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

// プレイヤー情報を取得する
func GetPlayer(roomId, playerId string) (*db.Player, error) {
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, orgerrors.NewPlayerNotFoundError("player not found", map[string]string{"roomId": roomId, "playerId": playerId})
	}

	return player, nil
}

// ゲームの全プレイヤー情報を取得する
func GetPlayers(roomId string) ([]db.Player, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, orgerrors.NewPlayerNotFoundError("players not found", map[string]string{"roomId": roomId})
	}

	return players, nil
}

// プレイヤー情報を取得する
func GetPlayerInfo(roomId, playerId string) (*responses.PlayerInfoResponse, error) {
	player, e := GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	return responses.GeneratePlayerInfoResponse(*player), nil
}

// 全プレイヤーが次フェーズに移行する準備ができている状態にする
func SetAllPlayersReady(roomId string) error {
	players, e := GetPlayers(roomId)
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
	players, e := GetPlayers(roomId)
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
func IsAllPlayersReadyByRoomId(roomId string) (bool, error) {
	players, e := GetPlayers(roomId)
	if e != nil {
		return false, e
	}

	return IsAllPlayersReady(players), nil
}

// 全プレイヤーが次フェーズに移行する準備ができているか
func IsAllPlayersReady(players []db.Player) bool {
	ready := true
	for _, p := range players {
		if p.ForceReady {
			return true
		}

		if !p.Ready {
			ready = false
		}
	}

	return ready
}

// プレイヤーにカードを追加する
func AppendCard(roomId, playerId string, appendCards []string) (*db.Player, error) {
	player, e := GetPlayer(roomId, playerId)
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
		return nil, orgerrors.NewValidationError("", "coin shortage", nil)
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
	return xid.New().String()
}
