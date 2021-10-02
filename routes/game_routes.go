package routes

import (
	"encoding/json"
	"log"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/models/requests"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/services/logic"

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
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorValidation, e))
			return
		}

		roomId, e := logic.GetRandomGameId()
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorInternal, e))
			return
		}

		if roomId == "" {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorGameNotFound, nil))
			return
		}

		u, e := logic.CreateNewPlayer(req.PlayerName, roomId, false)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorInternal, nil))
			return
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(u))
	})

	server.OnEvent("/", consts.ToServerJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		req := &requests.JoinFriendMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorValidation, e))
			return
		}

		u, e := logic.CreateNewPlayer(req.PlayerName, req.RoomID, false)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorInternal, e))
			return
		}
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(u))
	})

	server.OnEvent("/", consts.ToServerCreateMatch, func(s socketio.Conn, msg string) {
		req := &requests.CreateMatch{}
		if e := valid(msg, req); e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorValidation, e))
			return
		}

		ret, e := logic.CreateNewGame(req.PlayerName)
		if e != nil {
			s.Emit(consts.FromServerGameJoin, responseError(*responses.ErrorInternal, e))
			return
		}
		s.LeaveAll()
		s.Join(ret.RoomID)

		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, response(ret))
	})
}

// リクエストメッセージの構造体への変換 & バリデーション
func valid(reqBody string, result interface{}) error {
	if e := json.Unmarshal([]byte(reqBody), result); e != nil {
		return e
	}

	v := validator.New()
	if e := v.Struct(result); e != nil {
		return e
	}

	return nil
}

// インスタンスをレスポンス形式（JSON文字列）に変換する
func response(val interface{}) string {
	retJson, _ := json.Marshal(val)
	return string(retJson)
}

// インスタンスをレスポンス形式（JSON文字列）に変換する
func responseError(respErr responses.Error, err error) string {
	if respErr.Code == responses.ErrorInternal.Code {
		log.Println("[ERROR]", err)
	}

	retJson, _ := json.Marshal(respErr)
	return string(retJson)
}
