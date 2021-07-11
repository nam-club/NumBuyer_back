// プレイヤー関連のDB操作
package db

import (
	"nam-club/NumBuyer_back/consts"

	"github.com/gomodule/redigo/redis"
)

// Connection
func Connection() redis.Conn {
	c, err := redis.Dial("tcp", consts.Env.RedisUrl)
	if err != nil {
		panic(err)
	}
	return c
}

// データの登録(Redis: SET key value)
func Set(key, value string, c redis.Conn) string {
	res, err := redis.String(c.Do("SET", key, value))
	if err != nil {
		panic(err)
	}
	return res
}

// データの取得(Redis: GET key)
func Get(key string, c redis.Conn) string {
	res, err := redis.String(c.Do("GET", key))
	if err != nil {
		panic(err)
	}
	return res
}
