package cron

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/robfig/cron/v3"
)

func LaunchCron(server *socketio.Server) {
	c := cron.New()
	c.AddFunc("@every 3600s", func() {
		clean()
	})
	c.Start()
}

// DBのお掃除バッチ
func clean() {
	fmt.Printf("call2\n")
}
