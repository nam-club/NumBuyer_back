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

	server.OnEvent("/", consts.ToServerJoinQuickMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinQuickMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
			return
		}

		roomId, e := logic.GetRandomRoomId()
		if e != nil {
			// 部屋が見つからなかった場合は新規作成
			switch errors.Unwrap(e).(type) {
			case *orgerrors.GameNotFoundError:
				resp, e := logic.CreateNewGame(req.PlayerName)
				if e != nil {
					s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
					return
				}

				s.LeaveAll()
				s.Join(resp.RoomID)

				server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.Response(resp))
				return
			default:
				s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
				return
			}
		} else {
			// 部屋が見つかった場合はその部屋に参加
			player, e := logic.CreateNewPlayer(req.PlayerName, roomId, false)
			if e != nil {
				s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
				return
			}

			resp := responses.JoinResponse{RoomID: roomId, PlayerID: player.PlayerID}
			server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.Response(resp))
		}
	})

	server.OnEvent("/", consts.ToServerJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinFriendMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
			return
		}

		player, e := logic.CreateNewPlayer(req.PlayerName, req.RoomID, false)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
			return
		}
		resp := responses.JoinResponse{RoomID: req.RoomID, PlayerID: player.PlayerID}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.Response(resp))
	})

	server.OnEvent("/", consts.ToServerCreateMatch, func(s socketio.Conn, msg string) {
		req := &requests.CreateMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
			return
		}

		resp, e := logic.CreateNewGame(req.PlayerName)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, utils.ResponseError(e))
			return
		}
		s.LeaveAll()
		s.Join(resp.RoomID)
		// フェーズのタイマーをスタート
		logic.NewPhaseScheduler(resp.RoomID, server).Start()
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.Response(resp))
	})

	server.OnEvent("/", consts.ToServerGamePlayersInfo, func(s socketio.Conn, msg string) {
		req := &requests.GamePlayerInfo{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGamePlayersInfo, utils.ResponseError(e))
			return
		}
		resp, e := logic.GetPlayersInfo(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FromServerGamePlayersInfo, utils.ResponseError(e))
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGamePlayersInfo, utils.Response(resp))
	})

	server.OnEvent("/", consts.ToServerGameNextTurn, func(s socketio.Conn, msg string) {
		req := &requests.GameNextTurn{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameNextTurn, utils.ResponseError(e))
			return
		}
		resp, e := logic.NextTurn(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FromServerGameNextTurn, utils.ResponseError(e))
			return
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameNextTurn, utils.Response(resp))
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
