package routes

import (
	"log"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/services/logic"
	"nam-club/NumBuyer_back/services/utils"

	socketio "github.com/googollee/go-socket.io"
)

func RoutesGame(server *socketio.Server) {

	server.OnEvent("/", consts.ToServerJoinQuickMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		u := logic.CreateNewPlayer("JUNPEI", "rid", false)
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.ToResponseFormat(u))
	})

	server.OnEvent("/", consts.ToServerJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		u := logic.CreateNewPlayer("FRIEND", "JUNPEI", false)
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.ToResponseFormat(u))
	})

	server.OnEvent("/", consts.ToServerCreateMatch, func(s socketio.Conn, msg string) {
		log.Println("creatematch: ", msg)

		g := logic.CreateNewGame("OWENER")
		s.LeaveAll()
		s.Join(g.RoomID)

		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.ToResponseFormat(g))
	})
}
