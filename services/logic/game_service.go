package logic

import (
	"crypto/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"time"
)

// 新規ゲームを生成する
func CreateNewGame(owner string) (*responses.JoinResponse, error) {

	var id string
	var e error
	if id, e = generateRoomId(); e != nil {
		return nil, e
	}

	g := &db.Game{
		RoomID: id,
		State: db.State{
			Phase:       consts.PhaseBeforeStart.Value,
			Auction:     "",
			Answer:      "",
			ChangedTime: time.Now().Format(time.RFC3339),
		},
	}

	if _, e = db.SetGame(id, g); e != nil {
		return nil, e
	}
	player, e := CreateNewPlayer(owner, id, true)
	if e != nil {
		return nil, e
	}

	ret := &responses.JoinResponse{RoomID: id, PlayerID: player.PlayerID, IsOwner: player.IsOwner}

	return ret, nil
}

// ランダムなゲームIDを一つ取得する
func GetRandomRoomId() (string, error) {
	r, e := db.GetRandomRoomId()
	if e != nil {
		return "", e
	}
	return r, nil
}

// 次ターンで必要な情報を返却する
func NextTurn(roomId, playerId string) (*responses.NextTurnResponse, error) {
	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	player.Ready = true
	player, e = db.AddPlayer(roomId, player)
	if e != nil {
		return nil, e
	}

	return responses.GenerateNextTurnResponse(*player, *game), nil
}

// 次フェーズに移行する
func NextPhase(nextPhase consts.Phase, roomId string) (*responses.NextPhaseResponse, error) {
	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	game.State.Phase = nextPhase.Value
	db.SetGame(roomId, game)

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}
	for _, p := range players {
		p.Ready = false
		db.AddPlayer(roomId, &p)
	}

	return responses.GenerateNextPhaseResponse(players, nextPhase), nil
}

// ゲームを開始する
func StartGame(roomId string) error {
	game, err := db.GetGame(roomId)
	if err != nil {
		return orgerrors.NewGameNotFoundError("")
	}
	game.State.Phase = consts.PhaseWating.Value
	db.SetGame(roomId, game)

	return nil
}

// ゲームIDを生成する
func generateRoomId() (string, error) {
	const letters = "0123456789"

	for i := 0; i < 3; i++ {

		// 乱数を生成
		b := make([]byte, 10)
		if _, err := rand.Read(b); err != nil {
			return "", orgerrors.NewInternalServerError("")
		}

		var result string
		for _, v := range b {
			// index が letters の長さに収まるように調整
			result += string(letters[int(v)%len(letters)])
		}
		if b, _ := db.ExistsGame(result); b == false {
			return result, nil
		}
	}
	return "", orgerrors.NewInternalServerError("create room id error")
}
