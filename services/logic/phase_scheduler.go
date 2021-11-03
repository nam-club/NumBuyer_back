package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"

	"time"

	socketio "github.com/googollee/go-socket.io"
)

type PhaseSheduler struct {
	roomId string
	server *socketio.Server
}

func CanCreateGameScheduler(roomId string) error {
	if exists, _ := db.ExistsGame(roomId); !exists {
		return orgerrors.NewGameNotFoundError("")
	}
	game, e := db.GetGame(roomId)
	if e != nil {
		return orgerrors.NewInternalServerError("")
	}
	if p, e := consts.ParsePhase(game.State.Phase); e != nil {
		return orgerrors.NewInternalServerError("")
	} else if p != consts.PhaseBeforeStart {
		return orgerrors.NewValidationError("game already started")
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
	for {
		time.Sleep(2 * time.Second)

		game, e := db.GetGame(o.roomId)
		if e != nil {
			o.clean()
			break
		}
		startTime, e := time.Parse(time.RFC3339, game.State.ChangedTime)
		if e != nil {
			o.clean()
			break
		}

		phase, e := consts.ParsePhase(game.State.Phase)
		if e != nil {
			o.clean()
			break
		}

		var threshold int
		var next consts.Phase
		switch phase {
		case consts.PhaseBeforeStart:
			threshold = consts.PhaseBeforeStart.Duration
			next = consts.PhaseWating
		case consts.PhaseWating:
			threshold = consts.PhaseWating.Duration
			next = consts.PhaseBeforeAuction
		case consts.PhaseBeforeAuction:
			threshold = consts.PhaseBeforeAuction.Duration
			next = consts.PhaseAuction
		case consts.PhaseAuction:
			threshold = consts.PhaseAuction.Duration
			next = consts.PhaseAuctionResult
		case consts.PhaseAuctionResult:
			threshold = consts.PhaseAuctionResult.Duration
			next = consts.PhaseCalculate
		case consts.PhaseCalculate:
			threshold = consts.PhaseCalculate.Duration
			next = consts.PhaseCalculateResult
		case consts.PhaseCalculateResult:
			threshold = consts.PhaseCalculateResult.Duration
			next = consts.PhaseWating
		case consts.PhaseEnd:
			o.clean()
			break
		default:
			// 呼び出されないケース
			o.clean()
			break
		}

		fmt.Printf("thredhold: %v, current: %v, next: %v\n", threshold, phase, next)

		if startTime.Add(time.Duration(consts.TimeAutoEnd) * time.Second).Before(time.Now()) {
			o.clean()
			break
		}

		if threshold != consts.PhaseTimeValueInfinite {
			if startTime.Add(time.Duration(threshold) * time.Second).Before(time.Now()) {
				o.phaseFinishAction(phase, next)
			}
		} else if ready, _ := IsAllPlayersReady(o.roomId); ready {
			o.phaseFinishAction(phase, next)
		}
	}
}

func (o *PhaseSheduler) phaseFinishAction(current, next consts.Phase) {

	switch current {
	case consts.PhaseBeforeStart:
		o.nextPhase(next)
	case consts.PhaseWating:
		o.nextPhase(next)
	case consts.PhaseBeforeAuction:
		o.nextPhase(next)
	case consts.PhaseAuction:
		o.auctionFinishAction(next)
	case consts.PhaseAuctionResult:
		o.nextPhase(next)
	case consts.PhaseCalculate:
		o.calculateFinishAction(next)
	case consts.PhaseCalculateResult:
		o.calculateResultFinishAction(next)
	}
}
func (o *PhaseSheduler) nextPhase(next consts.Phase) {

	resp, e := NextPhase(next, o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
		return
	}
	o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.Response(resp))
}

func (o *PhaseSheduler) auctionFinishAction(next consts.Phase) {

	// TODO 落札者を決定しレスポンスに詰める処理をここに。落札者がいない場合はbroadcastしない
	resp := &responses.BuyNotifyResponse{}
	o.server.BroadcastToRoom("/", o.roomId, consts.FSGameBuyNotify, utils.Response(resp))

	o.nextPhase(next)
}

func (o *PhaseSheduler) calculateFinishAction(next consts.Phase) {

	if finished, _ := IsMeetClearCondition(o.roomId); finished {
		// ゲーム終了処理
		resp, e := FinishGame(o.roomId)
		if e != nil {
			o.server.BroadcastToRoom("/", o.roomId, consts.FSGameFinishGame, utils.ResponseError(e))
			return
		}
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameFinishGame, utils.Response(resp))
	} else {
		// TODO 正解者一覧の抽出処理
		resp := &responses.CorrectPlayersResponse{}
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameCorrectPlayers, utils.Response(resp))
		o.nextPhase(next)
	}

}

func (o *PhaseSheduler) calculateResultFinishAction(next consts.Phase) {
	// TODO 次のAnswerを生成する処理
	resp := &responses.UpdateAnswerResponse{}
	o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateAnswer, utils.Response(resp))
}

func (o *PhaseSheduler) finishGame() {
	FinishGame(o.roomId)
}

func (o *PhaseSheduler) clean() {
	db.DeleteGame(o.roomId)
}
