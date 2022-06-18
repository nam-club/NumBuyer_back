// redisへの接続部分
package db

import (
	"nam-club/NumBuyer_back/config"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type MutexMode int

const (
	Safe MutexMode = iota
	Force
)

type RedisHandler struct {
	DBIndex int
	pool    *redis.Pool
}

func NewRedisHandler(dbIndex int) (newHandler *RedisHandler) {
	newHandler = &RedisHandler{DBIndex: dbIndex}
	newHandler.pool = newPool(config.Env.RedisUrl, dbIndex)
	return
}

func newPool(addr string, dbIndex int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr, redis.DialDatabase(dbIndex)) },
	}
}

// データの登録
func (o *RedisHandler) Set(key, value string) (string, error) {
	conn := o.pool.Get()
	res, err := redis.String(conn.Do("SET", key, value))
	defer conn.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// データの登録。TTL付きでロックする
// セットに成功したら"OK", 失敗したらnilが返される
func (o *RedisHandler) SetNXEX(key, value string, ttl int) (interface{}, error) {
	conn := o.pool.Get()
	res, err := conn.Do("SET", key, value, "NX", "EX", ttl)
	defer conn.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの登録
func (o *RedisHandler) HSet(key, field, value string) (int64, error) {
	conn := o.pool.Get()
	res, err := redis.Int64(conn.Do("HSET", key, field, value))
	defer conn.Close()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// データの取得
func (o *RedisHandler) Get(key string) (string, error) {
	conn := o.pool.Get()
	res, err := redis.String(conn.Do("GET", key))
	defer conn.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// データの削除
func (o *RedisHandler) Delete(key string) (int, error) {
	conn := o.pool.Get()
	res, err := redis.Int(conn.Do("DEL", key))
	defer conn.Close()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの取得
func (o *RedisHandler) HGet(key, field string) (string, error) {
	conn := o.pool.Get()
	res, err := redis.String(conn.Do("HGET", key, field))
	defer conn.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// Hashデータの取得
func (o *RedisHandler) HVals(key string) ([][]byte, error) {
	conn := o.pool.Get()
	res, err := redis.ByteSlices(conn.Do("HVALS", key))
	defer conn.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

// データベースに存在するキーの数を取得
func (o *RedisHandler) DBSize() (int64, error) {
	conn := o.pool.Get()
	res, err := redis.Int64(conn.Do("DBSIZE"))
	defer conn.Close()
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return res, nil
}

// ランダムでキーをひとつ選択
func (o *RedisHandler) RandomKey() (string, error) {
	conn := o.pool.Get()
	res, err := redis.String(conn.Do("RANDOMKEY"))
	defer conn.Close()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

// データの取得(Redis: GET key)
func (o *RedisHandler) Exists(key string) (bool, error) {
	conn := o.pool.Get()
	res, err := redis.Bool(conn.Do("EXISTS", key))
	defer conn.Close()
	if err != nil {
		return false, errors.WithStack(err)
	}

	return res, nil
}

func (o *RedisHandler) Scan(iter int) (int, []string, error) {
	conn := o.pool.Get()
	defer conn.Close()

	var keys []string
	if arr, err := redis.Values(conn.Do("SCAN", iter)); err != nil {
		return 0, nil, errors.WithStack(err)
	} else {
		iter, _ = redis.Int(arr[0], nil)
		keys, _ = redis.Strings(arr[1], nil)
	}

	return iter, keys, nil
}
