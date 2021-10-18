package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/utils"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type PhaseSheduler struct {
	roomId string
	server *socketio.Server
}

func NewPhaseScheduler(roomId string, server *socketio.Server) *PhaseSheduler {
	return &PhaseSheduler{roomId: roomId}
}

func (o *PhaseSheduler) Start() {
	go o.monitor()
}

func (o *PhaseSheduler) monitor() {
	for {
		time.Sleep(1 * time.Second)

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
			threshold = consts.PhaseTimeBeforeStart
			next = consts.PhaseAuction
		case consts.PhaseAuction:
			threshold = consts.PhaseTimeAuction
			next = consts.PhaseCalculate
		case consts.PhaseCalculate:
			threshold = consts.PhaseTimeCalculate
			next = consts.PhaseResult
		case consts.PhaseResult:
			threshold = consts.PhaseTimeResult
			next = consts.PhaseAuction
		default:
			o.finish()
			break
		}

		if threshold != consts.PhaseTimeValueInfinite {
			if startTime.Add(time.Duration(threshold) * time.Second).Before(time.Now()) {
				o.nextPhase(next)
			}
		}
	}
}

func (o *PhaseSheduler) nextPhase(phase consts.Phase) {
	resp, e := NextPhase(phase, o.roomId)
	if e != nil {
		o.server.BroadcastToRoom("/", o.roomId, consts.FromServerGameUpdateState, utils.ResponseError(e))
	}
	o.server.BroadcastToRoom("/", o.roomId, consts.FromServerGameUpdateState, utils.Response(resp))
}

func (o *PhaseSheduler) finish() {
	db.DeleteGame(o.roomId)
}
