package routes

import (
	"encoding/json"
	"log"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/requests"
	"nam-club/NumBuyer_back/services/logic"

	socketio "github.com/googollee/go-socket.io"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func RoutesGame(server *socketio.Server) {

	server.OnEvent("/", consts.ToServerJoinQuickMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinQuickMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		roomId, e := logic.GetRandomRoomId()
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		if roomId == "" {
			s.Emit(consts.FromServerGameJoin, responseError(orgerrors.NewGameNotFoundError("")))
			return
		}

		_, e = logic.CreateNewPlayer(req.PlayerName, roomId, false)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		resp, e := logic.GeneratePlayersResponse(roomId)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(resp))
	})

	server.OnEvent("/", consts.ToServerJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinFriendMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		_, e := logic.CreateNewPlayer(req.PlayerName, req.RoomID, false)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		resp, e := logic.GeneratePlayersResponse(req.RoomID)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(resp))
	})

	server.OnEvent("/", consts.ToServerCreateMatch, func(s socketio.Conn, msg string) {
		req := &requests.CreateMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}

		ret, e := logic.CreateNewGame(req.PlayerName)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(e))
			return
		}
		s.LeaveAll()
		s.Join(ret.Players[0].RoomID)

		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(ret))
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

// インスタンスをレスポンス形式（JSON文字列）に変換する
func response(val interface{}) string {
	retJson, _ := json.Marshal(val)
	return string(retJson)
}

// インスタンスをレスポンス形式（JSON文字列）に変換する
func responseError(err error) string {
	errUnwrap := errors.Unwrap(err)
	var retJson []byte
	switch e := errUnwrap.(type) {
	case *orgerrors.ValidationError, *orgerrors.GameNotFoundError:
		retJson, _ = json.Marshal(e)
		break
	case *orgerrors.InternalServerError:
		log.Printf("[ERROR] %+v\n", e)
		retJson, _ = json.Marshal(e)
		break
	default:
		log.Printf("[ERROR] %+v\n", err)
		retJson, _ = json.Marshal(errors.Unwrap(orgerrors.NewInternalServerError("")))
		break
	}

	return string(retJson)
}
