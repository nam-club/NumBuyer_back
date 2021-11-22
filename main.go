package main

import (
	"log"
	"nam-club/NumBuyer_back/config"
	"nam-club/NumBuyer_back/routes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
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
	// 設定初期化
	config.InitConfig()

	// サーバーセットアップ
	router := gin.New()
	opt := engineio.Options{PingTimeout: time.Minute * 3, PingInterval: time.Second * 25}
	server := socketio.NewServer(&opt)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.LeaveAll()
		log.Println("connected:", s.ID())
		return nil
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

	routes.RoutesGame(server)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:3000"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("../asset"))

	if err := router.Run(":8001"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
