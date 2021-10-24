// redisへの接続部分
package db

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type RedisHandler struct {
	DBIndex int
	conn    redis.Conn
}

func NewRedisHandler(dbIndex int) *RedisHandler {

	ret := &RedisHandler{DBIndex: dbIndex}
	// redis接続
	c, err := redis.Dial("tcp", consts.Env.RedisUrl, redis.DialDatabase(dbIndex))
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	ret.conn = c
	// TODO Close()呼ばなくて問題ないか要確認
	// defer c.Close()

	return ret
}

// データの登録
func (o *RedisHandler) Set(key, value string) (string, error) {
	res, err := redis.String(o.conn.Do("SET", key, value))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの登録
func (o *RedisHandler) HSet(key, field, value string) (int64, error) {
	res, err := redis.Int64(o.conn.Do("HSET", key, field, value))
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// データの取得
func (o *RedisHandler) Get(key string) (string, error) {
	res, err := redis.String(o.conn.Do("GET", key))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// データの削除
func (o *RedisHandler) Delete(key string) (int, error) {
	res, err := redis.Int(o.conn.Do("DEL", key))
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの取得
func (o *RedisHandler) HGet(key, field string) (string, error) {
	res, err := redis.String(o.conn.Do("HGET", key, field))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの取得
func (o *RedisHandler) HVals(key string) ([][]byte, error) {
	res, err := redis.ByteSlices(o.conn.Do("HVALS", key))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// データベースに存在するキーの数を取得
func (o *RedisHandler) DBSize() (int64, error) {
	res, err := redis.Int64(o.conn.Do("DBSIZE"))
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// ランダムでキーをひとつ選択
func (o *RedisHandler) RandomKey() (string, error) {
	res, err := redis.String(o.conn.Do("RANDOMKEY"))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// データの取得(Redis: GET key)
func (o *RedisHandler) Exists(key string) (bool, error) {
	res, err := redis.Bool(o.conn.Do("EXISTS", key))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return res, nil
}
