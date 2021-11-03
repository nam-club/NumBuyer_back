package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
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
			o.finish()
			break
		}
		startTime, e := time.Parse(time.RFC3339, game.State.ChangedTime)
		if e != nil {
			o.finish()
			break
		}

		phase, e := consts.ParsePhase(game.State.Phase)
		if e != nil {
			o.finish()
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
		default:
			o.finish()
			break
		}

		fmt.Printf("thredhold: %v, current: %v, next: %v\n", threshold, phase, next)

		if startTime.Add(time.Duration(consts.TimeAutoEnd) * time.Second).Before(time.Now()) {
			o.finish()
			break
		}

		if threshold != consts.PhaseTimeValueInfinite {
			if startTime.Add(time.Duration(threshold) * time.Second).Before(time.Now()) {
				o.nextPhase(next)
			}
		} else if ready, _ := IsAllPlayersReady(o.roomId); ready {
			o.nextPhase(next)
		}
	}
}

func (o *PhaseSheduler) nextPhase(phase consts.Phase) {
	resp, e := NextPhase(phase, o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.ResponseError(e))
	}

	o.server.BroadcastToRoom("/", o.roomId, consts.FSGameUpdateState, utils.Response(resp))
}

func (o *PhaseSheduler) finish() {
	db.DeleteGame(o.roomId)
}
