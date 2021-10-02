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
	_ = conn.Send("MULTI")
	// コマンドの実行
	defer redis.Values(conn.Do("EXEC"))

	f()
}

// データの登録(Redis: SET key value)
func Set(key, value string) (string, error) {
	res, err := redis.String(conn.Do("SET", key, value))
	if err != nil {
		return "", err
	}
	return res, nil
}

// データの取得(Redis: GET key)
func Get(key string) (string, error) {
	res, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}
	return res, nil
}

// ランダムでキーをひとつ選択
func RandomKey() (string, error) {
	res, err := redis.String(conn.Do("RANDOMKEY"))
	if err != nil {
		return "", err
	}
	return res, nil
}

// データの取得(Redis: GET key)
func Exists(key string) (bool, error) {
	res, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return res, nil
}
