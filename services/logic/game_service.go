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

	"go.uber.org/zap"
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
			CurrentTurn:      1,
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
		return nil, orgerrors.NewGameNotFoundError("game not found", map[string]string{"roomId": roomId})
	}
	return r, nil
}

// 次ターンで必要な情報を返却する
func FetchNextTurnInfo(roomId, playerId string) (*responses.NextTurnResponse, error) {
	game, err := GetGame(roomId)
	if err != nil {
		return nil, err
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
	game, err := GetGame(roomId)
	if err != nil {
		return nil, err
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
		player.ForceReady = false
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

func GenerateUpdateState(roomId string, firedAbilities map[string][]*db.Ability) (*responses.UpdateStateResponse, error) {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}
	game, e := db.GetGame(roomId)
	if e != nil {
		return nil, e
	}
	currentPhase, e := consts.ParsePhase(game.State.Phase)
	if e != nil {
		return nil, e
	}

	return responses.GenerateUpdateStateResponse(players, currentPhase, firedAbilities), nil
}

// ゲームを開始する
func StartGame(roomId string) error {
	// ゲーム開始のバリデーション
	game, err := GetGame(roomId)
	if err != nil {
		return err
	}
	if game.State.Phase != consts.PhaseWaiting.Value {
		return orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_GAME_NOT_WAITING_PHASE, "game status is not waiting", nil)
	}

	players, err := GetPlayers(roomId)
	if err != nil {
		return err
	}
	if len(players) < game.PlayersMin || game.PlayersMax < len(players) {
		return orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_GAME_NOT_MEET_PLAYER_NUMS,
			"players num is not meet.",
			map[string]string{
				"min":     strconv.Itoa(game.PlayersMin),
				"max":     strconv.Itoa(game.PlayersMax),
				"current": strconv.Itoa(len(players))})
	}

	// ゲーム開始処理
	if _, e := db.DeleteJoinableGame(roomId); e != nil {
		return e
	}
	if e := SetAllPlayersReady(roomId); e != nil {
		return orgerrors.NewInternalServerError("set players status ready failed.", nil)
	}

	if _, e := ShuffleAnswer(roomId); e != nil {
		return orgerrors.NewInternalServerError("failed to shuffle answer.", nil)
	}

	if _, e := ShuffleAuctionCard(roomId); e != nil {
		return orgerrors.NewInternalServerError("failed to shuffle auction.", nil)
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

func IncreaseTurn(roomId string) error {
	game, err := db.GetGame(roomId)
	if err != nil {
		return err
	}
	game.State.CurrentTurn += 1
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

func IsJoinable(roomId string) (bool, error) {
	game, err := GetGame(roomId)
	if err != nil {
		return false, err
	}
	return game.State.Phase == consts.PhaseWaiting.Value, nil
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

	game, err := GetGame(roomId)
	if err != nil {
		return nil, err
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
			return "", orgerrors.NewInternalServerError("", nil)
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
	return "", orgerrors.NewInternalServerError("create room id error", nil)
}

// ロビーから抜ける
func LeaveLobby(playerId, roomId string) (*responses.PlayersInfoResponse, error) {
	if !CheckPhase(roomId, consts.PhaseWaiting) {
		return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_GAME_NOT_WAITING_PHASE, "not waiting phase", nil)
	}

	player, err := db.GetPlayer(roomId, playerId)
	if err != nil {
		return nil, err
	}
	if player.IsOwner {
		// プレイヤー情報を削除した上でオーナー権限を付け替え
		if _, err := db.DeletePlayer(roomId, playerId); err != nil {
			return nil, err
		}

		players, err := db.GetPlayers(roomId)
		if err != nil {
			return nil, err
		}

		if len(players) > 0 {
			players[0].IsOwner = true
			db.SetPlayer(roomId, &players[0])
		}
	} else {
		if _, err := db.DeletePlayer(roomId, playerId); err != nil {
			return nil, err
		}
	}
	players, err := db.GetPlayers(roomId)
	if err != nil {
		return nil, err
	}

	return responses.GeneratePlayersInfoResponse(players, roomId), nil

}

func Clean(roomId string) {
	utils.Log.Debug("start delete...", zap.String("roomId", roomId))
	db.DeleteJoinableGame(roomId)
	db.DeletePlayers(roomId)
	db.DeleteGame(roomId)
	utils.Log.Debug("complete delete", zap.String("roomId", roomId))
}
