package cron

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/utils"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func LaunchCron(server *socketio.Server) {
	c := cron.New()
	c.AddFunc("@every 7200s", func() {
		clean()
	})
	c.Start()
}

// DBのお掃除バッチ
func clean() {
	// 削除対象抽出
	iter := 0
	deleteKeys := []string{}
	for {
		var keys []string
		var err error
		if iter, keys, err = db.ScanGame(iter); err != nil {
			utils.Log.Error("batch: scan db", zap.String("error", err.Error()))
			return
		} else {
			for _, key := range keys {
				if g, err := db.GetGame(key); err != nil {
					utils.Log.Error("batch: get game failed", zap.String("error", err.Error()), zap.String("key", key))
				} else {
					createdAt, e := time.Parse(time.RFC3339, g.CreatedAt)
					if e != nil {
						deleteKeys = append(deleteKeys, key)
					} else {
						if createdAt.Add(time.Duration(consts.TimeAutoDelete) * time.Second).Before(time.Now()) {
							deleteKeys = append(deleteKeys, key)
						}

					}
				}
			}
		}
		if iter == 0 {
			break
		}
	}

	// 削除実施
	for _, del := range deleteKeys {
		if _, err := db.DeletePlayers(del); err != nil {
			utils.Log.Error("batch: failed to delete players", zap.String("error", err.Error()), zap.String("key", del))
		}
		if _, err := db.DeleteJoinableGame(del); err != nil {
			utils.Log.Error("batch: failed to delete joinablegame", zap.String("error", err.Error()), zap.String("key", del))
		}
		if _, err := db.DeleteGame(del); err != nil {
			utils.Log.Error("batch: failed to delete game", zap.String("error", err.Error()), zap.String("key", del))
		}
	}

}
