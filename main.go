// ※サンプルコード

package main

import (
	"log"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/services/db"
	"nam-club/NumBuyer_back/services/logic"
	"nam-club/NumBuyer_back/services/utils"
	"net/http"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func main() {
	// redis接続
	c := db.Connection()
	defer c.Close()

	// サーバーセットアップ
	router := gin.New()
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.LeaveAll()
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", consts.ToServerJoinQuickMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		u := logic.CreateNewPlayer("QUICK", "ITO")
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.ToResponseFormat(u))
	})

	server.OnEvent("/", consts.ToServerJoinFriendMatch, func(s socketio.Conn, msg string) {
		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(msg)

		u := logic.CreateNewPlayer("FRIEND", "JUNPEI")
		server.BroadcastToRoom("/", s.Rooms()[0], consts.FromServerGameJoin, utils.ToResponseFormat(u))
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Use(GinMiddleware("http://127.0.0.1:8080"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))

	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
