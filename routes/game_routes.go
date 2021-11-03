package routes

import (
	"encoding/json"
	"errors"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/requests"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/services/logic"
	"nam-club/NumBuyer_back/utils"

	socketio "github.com/googollee/go-socket.io"
	"gopkg.in/go-playground/validator.v9"
)

func RoutesGame(server *socketio.Server) {

	server.OnEvent("/", consts.TSJoinQuickMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinQuickMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		roomId, e := logic.GetRandomRoomId()
		if e != nil {
			// 部屋が見つからなかった場合は新規作成
			switch errors.Unwrap(e).(type) {
			case *orgerrors.GameNotFoundError:
				resp, e := logic.CreateNewGame(req.PlayerName)
				if e != nil {
					s.Emit(consts.FSGameJoin, utils.ResponseError(e))
					return
				}

				s.LeaveAll()
				s.Join(resp.RoomID)

				s.Emit(consts.FSGameJoin, utils.Response(resp))
				return
			default:
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}
		} else {
			// 部屋が見つかった場合はその部屋に参加
			player, e := logic.CreateNewPlayer(req.PlayerName, roomId, false)
			if e != nil {
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}

			resp := responses.JoinResponse{RoomID: roomId, PlayerID: player.PlayerID}
			s.Emit(consts.FSGameJoin, utils.Response(resp))
		}
	})

	server.OnEvent("/", consts.TSJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinFriendMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		player, e := logic.CreateNewPlayer(req.PlayerName, req.RoomID, false)
		if e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}
		resp := responses.JoinResponse{RoomID: req.RoomID, PlayerID: player.PlayerID}
		s.Emit(consts.FSGameJoin, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSCreateMatch, func(s socketio.Conn, msg string) {
		req := &requests.CreateMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		resp, e := logic.CreateNewGame(req.PlayerName)
		if e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		// フェーズのタイマーをスタート
		if e := logic.CanCreateGameScheduler(resp.RoomID); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}
		logic.NewPhaseScheduler(resp.RoomID, server).Start()

		s.LeaveAll()
		s.Join(resp.RoomID)
		s.Emit(consts.FSGameJoin, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGamePlayersInfo, func(s socketio.Conn, msg string) {
		req := &requests.GamePlayerInfo{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGamePlayersInfo, utils.ResponseError(e))
			return
		}
		resp, e := logic.GetPlayersInfo(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FSGamePlayersInfo, utils.ResponseError(e))
			return
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGamePlayersInfo, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGameStart, func(s socketio.Conn, msg string) {
		req := &requests.GameStart{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}
		if e := logic.CanCreateGameScheduler(req.RoomID); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}

		if e := logic.StartGame(req.RoomID); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}

		resp := &responses.GameStartResponse{StartFlag: true}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameStart, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGameNextTurn, func(s socketio.Conn, msg string) {
		req := &requests.GameNextTurn{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameNextTurn, utils.ResponseError(e))
			return
		}
		resp, e := logic.NextTurn(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FSGameNextTurn, utils.ResponseError(e))
			return
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameNextTurn, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGameBid, func(s socketio.Conn, msg string) {
		req := &requests.GameBid{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameBid, utils.ResponseError(e))
			return
		}
		resp := &responses.BidResponse{PlayerName: "JUNPEI", Coin: 99}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameBid, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGameBuy, func(s socketio.Conn, msg string) {
		req := &requests.GameBuy{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameBuyUpdate, utils.ResponseError(e))
			return
		}
		resp := &responses.BuyUpdateResponse{PlayerID: "ID_JUNPEI_UPDATED", Coin: 66, Cards: []string{"21", "87", "*", "/", "4"}}
		s.Emit(consts.FSGameBuyUpdate, utils.Response(resp))
	})

	server.OnEvent("/", consts.TSGameCalculate, func(s socketio.Conn, msg string) {
		req := &requests.GameCalculate{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FSGameCalculateResult, utils.ResponseError(e))
			return
		}
		resp := &responses.CalculateResponse{IsCorrectAnswer: true, PlayerID: "ID_JUNPEI", Coin: 77, Cards: []string{"12", "33", "9", "+", "-"}}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameCalculateResult, utils.Response(resp))
	})
}

// リクエストメッセージの構造体への変換 & バリデーション
func valid(reqBody string, result interface{}) error {
	if e := json.Unmarshal([]byte(reqBody), result); e != nil {
		return orgerrors.NewValidationError(e.Error())
	}

	v := validator.New()
	if e := v.Struct(result); e != nil {
		return orgerrors.NewValidationError(e.Error())
	}

	return nil
}
