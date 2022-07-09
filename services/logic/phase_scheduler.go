package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"strconv"

	"time"

	socketio "github.com/googollee/go-socket.io"
	"go.uber.org/zap"
)

type PhaseSheduler struct {
	roomId string
	server *socketio.Server
}

func CanCreateGameScheduler(roomId string) error {
	if exists, _ := db.ExistsGame(roomId); !exists {
		return orgerrors.NewGameNotFoundError("", nil)
	}
	game, e := GetGame(roomId)
	if e != nil {
		return e
	}
	if p, e := consts.ParsePhase(game.State.Phase); e != nil {
		return orgerrors.NewInternalServerError("", nil)
	} else if p != consts.PhaseWaiting {
		return orgerrors.NewValidationError("", "game already started", nil)
	}
	return nil
}

func NewPhaseScheduler(roomId string, server *socketio.Server) *PhaseSheduler {
	return &PhaseSheduler{roomId: roomId, server: server}
}

func (o *PhaseSheduler) Start() {
	go o.monitor()
}

func (o *PhaseSheduler) monitor() {
LOOP:
	for {
		time.Sleep(1 * time.Second)

		game, e := GetGame(o.roomId)
		if e != nil {
			Clean(o.roomId)
			break LOOP
		}
		startTime, e := time.Parse(time.RFC3339, game.State.PhaseChangedTime)
		if e != nil {
			Clean(o.roomId)
			break LOOP
		}

		phase, e := consts.ParsePhase(game.State.Phase)
		if phase.NextPhase == nil || *phase.NextPhase == consts.PhaseEnd || e != nil {
			Clean(o.roomId)
			break LOOP
		}
		nextPhase := *phase.NextPhase

		// フェーズの更新が指定時間以上ない場合強制終了
		if startTime.Add(time.Duration(consts.TimeAutoEnd) * time.Second).Before(time.Now()) {
			Clean(o.roomId)
			break LOOP
		}

		players, _ := GetPlayers(o.roomId)
		if len(players) == 0 {
			Clean(o.roomId)
			break LOOP
		}

		// 表示更新待ちのアビリティを取得し、ステータスを更新する
		firedAbilities := ProccessReadyUpdateAbilities(o.roomId, players)
		if len(firedAbilities) > 0 {
			updateStateResp, _ := GenerateUpdateState(o.roomId, firedAbilities)
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.Response(updateStateResp))
		}

		// 全プレイヤーが準備済み または 指定時間を経過した場合、次フェーズに移動する
		if IsAllPlayersReady(players) {
			o.phaseFinishAction(phase, nextPhase)
		} else if phase.Duration != consts.PhaseTimeValueInfinite &&
			startTime.Add(time.Duration(phase.Duration+phase.Grace)*time.Second).Before(time.Now()) {
			// 前ターンの計算フェーズで正答者がいなかった場合、ターゲットカード更新フェーズをスキップ
			if phase == consts.PhaseGiveCards && game.State.SkipShowTarget {
				nextPhase = consts.PhaseShowAuction
			}
			o.phaseFinishAction(phase, nextPhase)
		}
	}
}

func (o *PhaseSheduler) phaseFinishAction(current, next consts.Phase) {
	utils.Log.Debug("phase finishing",
		zap.String("roomId", o.roomId),
		zap.String("current", fmt.Sprintf("%v", current)),
		zap.String("next", fmt.Sprintf("%v", next)))

	playerIds, e := db.GetPlayerIds(o.roomId)
	if e != nil {
		utils.Log.Error("users info fetch failed.", zap.String("error", fmt.Sprintf("%v", e)))
	}
	// ゲームをロックする
	for retry := 0; retry < consts.MutexRetryCount; retry++ {
		lockFailed := false
		for _, playerId := range playerIds {
			if locked, e := db.SetLock(db.CreateLockKey(o.roomId, playerId), consts.MutexTTL); locked {
				if e != nil {
					utils.Log.Warn("lock failed.", zap.String("error", fmt.Sprintf("%v", e)))
				}
				lockFailed = true
				db.DeleteLock(db.CreateLockKey(o.roomId, playerId))
			}
		}
		if !lockFailed {
			break
		}
		utils.Log.Warn("game db is locked. try lock again...", zap.String("retry count", strconv.Itoa(retry)))
		time.Sleep(consts.MutexRetrySpan * time.Millisecond)
	}

	// ロック情報を最後に削除
	for _, playerId := range playerIds {
		defer db.DeleteLock(db.CreateLockKey(o.roomId, playerId))
	}

	switch current {
	case consts.PhaseWaiting:
		o.nextPhase(next)
	case consts.PhaseReady:
		o.nextPhase(next)
	case consts.PhaseGiveCards:
		o.nextPhase(next)
	case consts.PhaseShowTarget:
		o.nextPhase(next)
	case consts.PhaseShowAuction:
		o.nextPhase(next)
	case consts.PhaseAuction:
		o.auctionFinishAction(next)
	case consts.PhaseAuctionResult:
		o.nextPhase(next)
	case consts.PhaseCalculate:
		o.calculateFinishAction(next)
	case consts.PhaseCalculateResult:
		o.setUpNextTurn(next)
	case consts.PhaseNextTurn:
		o.nextPhase(next)
	}
}
func (o *PhaseSheduler) nextPhase(next consts.Phase) error {

	resp, e := NextPhase(next, o.roomId)
	if e != nil {
		utils.Log.Debug("error broadcast", zap.String("roomId", o.roomId), zap.String("next", fmt.Sprintf("%v", next)))
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
		return e
	}
	utils.Log.Debug("broadcast", zap.String("roomId", o.roomId), zap.String("next", fmt.Sprintf("%v", next)))
	o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.Response(resp))
	return nil
}

func (o *PhaseSheduler) auctionFinishAction(next consts.Phase) {

	resp, e := DetermineBuyer(o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.ResponseError(e))
		return
	}

	game, e := db.GetGame(o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.ResponseError(e))
		return
	}
	currentAuction := game.State.Auction
	if resp != nil {
		// 落札者が決まった時

		// 落札者にカードを追加
		buyer, e := AppendCard(o.roomId, resp.PlayerID, game.State.Auction)
		if e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.ResponseError(e))
			return
		}

		// プレイヤーのコインを減らす
		subtract, _ := strconv.Atoi(buyer.BuyAction.Value)
		buyer, e = SubtractCoin(o.roomId, resp.PlayerID, subtract)
		if e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.ResponseError(e))
			return
		}

		resp := &responses.BuyNotifyResponse{
			PlayerName:   buyer.PlayerName,
			Coin:         subtract,
			AuctionCards: currentAuction,
			IsPassAll:    false,
		}
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.Response(resp))
	} else {
		// 全員passした時

		resp := &responses.BuyNotifyResponse{
			AuctionCards: currentAuction,
			IsPassAll:    true,
		}
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.Response(resp))
	}

	// 次フェーズへ移動
	o.nextPhase(next)
}

func (o *PhaseSheduler) calculateFinishAction(next consts.Phase) {

	// アビリティ発動準備
	game, e := db.GetGame(o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
		return
	}
	if e = TryActivateAbilitiesIfHave(game, consts.AbilityIdNumViolence); e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
		return
	}

	// ゲーム終了条件を満たしているか
	if finished, _ := IsMeetClearCondition(o.roomId); finished {
		// 最新の状態を返却
		if state, e := GenerateUpdateState(o.roomId, nil); e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
			return
		} else {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.Response(state))
		}

		// ゲーム終了処理
		if resp, e := FinishGame(o.roomId); e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameFinishGame, utils.ResponseError(e))
			return
		} else {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameFinishGame, utils.Response(resp))
		}

	} else {
		// 正解者一覧を抽出し返却
		correctors, e := PickAllCorrector(o.roomId)
		if e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameCorrectPlayers, utils.ResponseError(e))
			return
		}

		resp := &responses.CorrectPlayersResponse{}
		for _, corrector := range correctors {
			resp.AnsPlayers = append(resp.AnsPlayers,
				responses.CorrectPlayer{
					PlayerName: corrector.PlayerName,
					AddedCoin: responses.AddedCoin{
						Total:        corrector.AnswerAction.PlusCoin,
						CardNumBonus: len(corrector.AnswerAction.Cards),
					}})
		}
		existsCorrector := len(resp.AnsPlayers) > 0
		resp.ExistsCorrect = existsCorrector

		// 正答者が一人でもいれば解答をシャッフル
		if existsCorrector {
			_, e := ShuffleAnswer(o.roomId)
			if e != nil {
				o.server.BroadcastToRoom("/", o.roomId, consts.FSGameCorrectPlayers, utils.ResponseError(e))
				return
			}
		}

		// ターゲット表示フェーズをスキップするフラグをセット
		// 正答者が一人もいなければスキップ
		if e = SetSkipShowTarget(o.roomId, !existsCorrector); e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.ResponseError(e))
			return
		}

		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameCorrectPlayers, utils.Response(resp))
		o.nextPhase(next)
	}
}

// 次ターンに移るときの準備をする
func (o *PhaseSheduler) setUpNextTurn(next consts.Phase) {
	ClearAndResetAuction(o.roomId)
	ClearCalculateAction(o.roomId)
	AddCardToAllPlayers(o.roomId)
	o.nextPhase(next)
}

// タイマーを指定した時間を残してリセットする
func ResetTimer(roomId string, remainSeconds int) (bool, error) {
	game, e := GetGame(roomId)
	if e != nil {
		return false, e
	}

	phase, e := consts.ParsePhase(game.State.Phase)
	if e != nil {
		return false, orgerrors.NewInternalServerError("", nil)
	}

	passed := phase.Duration - remainSeconds
	if passed <= 0 {
		return false, orgerrors.NewInternalServerError("", nil)
	}

	beforeTime, e := time.Parse(time.RFC3339, game.State.PhaseChangedTime)
	if e != nil {
		return false, orgerrors.NewInternalServerError("", nil)
	}

	// 現在時刻が前回のフェーズ変更から指定時間していない場合、何もしない
	passedThreashold := beforeTime.Add(time.Second * time.Duration(passed))
	if !time.Now().After(passedThreashold) {
		return false, nil
	}

	// 現在自国から指定時間を戻し、残り時間を引数に指定された時間にする
	game.State.PhaseChangedTime = time.Now().Add(time.Second * -time.Duration(passed)).Format(time.RFC3339)
	_, e = db.SetGame(roomId, game)
	if e != nil {
		return false, orgerrors.NewInternalServerError("", nil)
	}
	return true, nil
}
