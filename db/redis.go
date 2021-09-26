// redisへの接続部分
package db

import (
	"nam-club/NumBuyer_back/consts"

	"github.com/gomodule/redigo/redis"
)

var conn redis.Conn

func init() {
	// redis接続
	c, err := redis.Dial("tcp", consts.Env.RedisUrl)
	if err != nil {
		panic(err)
	}
	conn = c

	// TODO Close()呼ばなくて問題ないか要確認
	// defer c.Close()
}

func Atomic(f func()) {
	// トランザクションの開始
	var err error
	err = conn.Send("MULTI")
	if err != nil {
		panic(err)
	}
	f()
	// コマンドの実行
	_, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		panic(err)
	}
}

// データの登録(Redis: SET key value)
func Set(key, value string) string {
	res, err := redis.String(conn.Do("SET", key, value))
	if err != nil {
		panic(err)
	}
	return res
}

// データの取得(Redis: GET key)
func Get(key string) string {
	res, err := redis.String(conn.Do("GET", key))
	if err != nil {
		panic(err)
	}
	return res
}
