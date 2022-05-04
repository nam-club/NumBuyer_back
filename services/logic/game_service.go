package logic

import (
	"crypto/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"sort"
	"strconv"
	"time"
)

// 新規ゲームを生成する
func CreateNewGame(playerName string, playersMin, playersMax int, gameMode consts.GameMode, abilities []consts.Ability) (*responses.JoinResponse, error) {

	var id string
	var e error
	if id, e = generateRoomId(); e != nil {
		return nil, e
	}

	g := &db.Game{
		RoomID:     id,
		PlayersMin: playersMin,
		PlayersMax: playersMax,
		CreatedAt:  time.Now().Format(time.RFC3339),
		State: db.State{
			Phase:            consts.PhaseWaiting.Value,
			Auction:          []string{},
			Answer:           "",
			PhaseChangedTime: time.Now().Format(time.RFC3339),
		},
	}

	if _, e = db.SetGame(id, g); e != nil {
		return nil, e
	}

	if gameMode == consts.GameModeQuickMatch {
		if e = db.SetJoinableGame(id); e != nil {
			return nil, e
		}
	}

	// フレンドマッチならゲーム作成者にオーナー権限をつける
	isOwner := gameMode == consts.GameModeFriendMatch
	player, e := CreateNewPlayer(playerName, id, isOwner, abilities)
	if e != nil {
		return nil, e
	}

	return responses.GenerateJoinResponse(id, player)
}

// ランダムなゲームIDを一つ取得する
func GetRandomRoomId() (string, error) {
	r, e := db.GetRandomRoomId()
	if e != nil {
		return "", e
	}
	return r, nil
}

func GetGame(roomId string) (*db.Game, error) {
	r, e := db.GetGame(roomId)
	if e != nil {
		return nil, e
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
	// ゲーム状態を更新する
	game, err := db.GetGame(roomId)
	if err != nil {
		return nil, orgerrors.NewGameNotFoundError("")
	}
	game.State.Phase = nextPhase.Value
	game.State.PhaseChangedTime = time.Now().Format(time.RFC3339)
	game, _ = db.SetGame(roomId, game)

	playerIds, _ := db.GetPlayerIds(roomId)
	firedAbilities := map[string][]*db.Ability{}
	for _, playerId := range playerIds {
		// プレイヤーの準備状態をリセット
		player, _ := db.GetPlayer(roomId, playerId)
		player.Ready = false
		db.SetPlayer(roomId, player)

		// アビリティを発動する
		firedAbs := []*db.Ability{}
		for _, ability := range player.Abilities {
			if ab, _ := consts.ParseAbility(ability.ID); ab.Timing != consts.AbilityTimingWait {
				continue
			}
			firedAb, e := FireAbility(game, player, ability.ID)
			if e != nil {
				return nil, e
			}
			if firedAb != nil {
				firedAbs = append(firedAbs, firedAb)
			}

		}
		firedAbilities[player.PlayerID] = firedAbs
	}
	players, _ := db.GetPlayers(roomId)
	return responses.GenerateUpdateStateResponse(players, nextPhase, firedAbilities), nil
}

func GenerateUpdateState(nextPhase consts.Phase, roomId string) (*responses.UpdateStateResponse, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}
	return responses.GenerateUpdateStateResponse(players, nextPhase, map[string][]*db.Ability{}), nil
}

// ゲームを開始する
func StartGame(roomId string) error {
	// ゲーム開始のバリデーション
	game, err := db.GetGame(roomId)
	if err != nil {
		return orgerrors.NewValidationError("get room failed: " + roomId)
	}
	if game.State.Phase != consts.PhaseWaiting.Value {
		return orgerrors.NewValidationError("game status is not waiting")
	}

	players, err := db.GetPlayers(roomId)
	if err != nil {
		return orgerrors.NewValidationError("get players failed: " + roomId)
	}
	if len(players) < game.PlayersMin || game.PlayersMax < len(players) {
		return orgerrors.NewValidationError("players num is not meet. min=" + strconv.Itoa(game.PlayersMin) + ", max=" + strconv.Itoa(game.PlayersMax) + ", current=" + strconv.Itoa(len(players)))
	}

	// ゲーム開始処理
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
	players, err = db.GetPlayers(roomId)
	if err != nil {
		return err
	}
	for _, player := range players {
		player.Cards = utils.GenerateRandomCard(consts.InitialCardsNum)
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

	// Coin数でsort (降順）
	sort.Slice(players, func(i, j int) bool { return players[i].Coin > players[j].Coin })

	resp := &responses.FinishGameResponse{Players: make([]responses.FinishGamePlayers, len(players))}
	for i, player := range players {
		resp.Players[i].PlayerName = player.PlayerName
		resp.Players[i].Coin = player.Coin
		isSameRank := false // 同立順位か
		if i > 0 {
			isSameRank = players[i-1].Coin == players[i].Coin
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
