package cron

import (
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/utils"

	socketio "github.com/googollee/go-socket.io"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func LaunchCron(server *socketio.Server) {
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		clean()
	})
	c.Start()
}

// DBのお掃除バッチ
func clean() {
	iter := 0
	for {
		if iter, keys, err := db.ScanGame(iter); err != nil {
			utils.Log.Error("start delete...", zap.String("error", err.Error()))
			return
		} else {
			println("from cron", iter, keys)
		}

		if iter == 0 {
			break
		}
	}
}
