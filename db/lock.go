// ゲーム情報に関するDB操作
package db

var rl *RedisHandler

func init() {
	rl = NewRedisHandler( /*index=*/ 3)
}

// ゲームのロック情報を取得
func SetLock(lockKey string, ttl int) (bool, error) {
	r, e := rl.SetNXEX(lockKey, "LOCK", ttl)
	if e != nil {
		return false, e
	}

	return r == nil, nil
}

func ExistsLock(lockKey string) (bool, error) {
	return rl.Exists(lockKey)
}

func DeleteLock(lockKey string) error {
	_, e := rl.Delete(lockKey)
	if e != nil {
		return e
	}
	return nil
}
