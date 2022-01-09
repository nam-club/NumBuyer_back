package logic

import (
	"crypto/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"sort"
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
			Phase:            consts.PhaseWaiting.Value,
			Auction:          "",
			Answer:           "",
			PhaseChangedTime: time.Now().Format(time.RFC3339),
		},
	}

	if _, e = db.SetGame(id, g); e != nil {
		return nil, e
	}

	if e = db.SetJoinableGame(id); e != nil {
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
func FetchNextTurnInfo(roomId, playerId string) (*responses.NextTurnResponse, error) {
	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	player.Ready = true
	player, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}

	return responses.GenerateNextTurnResponse(*player, *game), nil
}

// 次フェーズに移行する
func NextPhase(nextPhase consts.Phase, roomId string) (*responses.UpdateStateResponse, error) {
	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	game.State.Phase = nextPhase.Value
	game.State.PhaseChangedTime = time.Now().Format(time.RFC3339)
	db.SetGame(roomId, game)

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}
	for _, p := range players {
		p.Ready = false
		db.SetPlayer(roomId, &p)
	}

	return responses.GenerateUpdateStateResponse(players, nextPhase), nil
}

func GenerateUpdateState(nextPhase consts.Phase, roomId string) (*responses.UpdateStateResponse, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}
	return responses.GenerateUpdateStateResponse(players, nextPhase), nil
}

// ゲームを開始する
func StartGame(roomId string) error {
	joinable := CheckPhase(roomId, consts.PhaseWaiting)
	if !joinable {
		return orgerrors.NewValidationError("game status is not waiting")
	}

	if _, e := db.DeleteJoinableGame(roomId); e != nil {
		return e
	}

	if e := SetAllPlayersReady(roomId); e != nil {
		return orgerrors.NewInternalServerError("set players status ready failed.")
	}

	if _, e := ShuffleAnswer(roomId); e != nil {
		return orgerrors.NewInternalServerError("failed to shuffle answer.")
	}

	if _, e := ShuffleAuctionCard(roomId); e != nil {
		return orgerrors.NewInternalServerError("failed to shuffle auction.")
	}

	// プレイヤーに初期カードを付与する
	players, err := db.GetPlayers(roomId)
	if err != nil {
		return err
	}
	for _, player := range players {
		for i := 0; i < consts.InitialCardsNum; i++ {
			player.Cards = append(player.Cards, utils.GenerateRandomCard())
		}
		db.SetPlayer(roomId, &player)
	}

	return nil
}

// ゲームのクリア条件を満たしているかチェックする
func IsMeetClearCondition(roomId string) (bool, error) {
	players, err := db.GetPlayers(roomId)
	if err != nil {
		return false, err
	}
	for _, player := range players {
		if player.Coin >= consts.CoinClearNum {
			return true, nil
		}
	}
	return false, nil
}

// ターゲット表示フェーズスキップフラグをセットする
func SetSkipShowTarget(roomId string, skip bool) error {
	game, err := db.GetGame(roomId)
	if err != nil {
		return err
	}
	game.State.SkipShowTarget = skip
	if _, err := db.SetGame(roomId, game); err != nil {
		return err
	}

	return nil
}

func CheckPhase(roomId string, phase consts.Phase) bool {
	game, err := db.GetGame(roomId)
	if err != nil {
		return false
	}
	return game.State.Phase == phase.Value
}

// ゲームを終了する（ゲーム終了条件のチェックは行わない）
func FinishGame(roomId string) (*responses.FinishGameResponse, error) {
	players, err := db.GetPlayers(roomId)
	if err != nil {
		return nil, err
	}

	// 回答時刻でsort (昇順）
	sort.Slice(players, func(i, j int) bool {
		tiStr := players[i].AnswerAction.AnswerTime
		tjStr := players[j].AnswerAction.AnswerTime
		if tiStr == "" {
			return false
		}
		if tjStr == "" {
			return true
		}
		ti, _ := time.Parse(time.RFC3339, tiStr)
		tj, _ := time.Parse(time.RFC3339, tjStr)
		return ti.Before(tj)
	})

	// Coin数でsort (降順）
	sort.Slice(players, func(i, j int) bool { return players[i].Coin > players[j].Coin })

	resp := &responses.FinishGameResponse{Players: make([]responses.FinishGamePlayers, len(players))}
	for i, player := range players {
		resp.Players[i].PlayerName = player.PlayerName
		resp.Players[i].Coin = player.Coin
		isSameRank := false // 同立順位か
		if i > 0 {
			if players[i-1].Coin == players[i].Coin &&
				players[i-1].AnswerAction.AnswerTime == players[i].AnswerAction.AnswerTime {
				isSameRank = true
			}
		}
		if isSameRank {
			resp.Players[i].Rank = resp.Players[i-1].Rank
		} else {
			resp.Players[i].Rank = i + 1
		}
	}

	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	game.State.Phase = consts.PhaseEnd.Value
	db.SetGame(roomId, game)

	return resp, nil
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
		if b, _ := db.ExistsGame(result); !b {
			return result, nil
		}
	}
	return "", orgerrors.NewInternalServerError("create room id error")
}
