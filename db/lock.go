// ゲーム情報に関するDB操作
package db

var rl *RedisHandler

type LockKey string

func init() {
	rl = NewRedisHandler( /*index=*/ 3)
}

func CreateLockKey(roomId, playerId string) LockKey {
	return LockKey(roomId + playerId)
}

// ゲームのロック情報を取得
func SetLock(lockKey LockKey, ttl int) (bool, error) {
	r, e := rl.SetNXEX(string(lockKey), "LOCK", ttl)
	if e != nil {
		return false, e
	}

	return r == nil, nil
}

func ExistsLock(lockKey LockKey) (bool, error) {
	return rl.Exists(string(lockKey))
}

func DeleteLock(lockKey LockKey) error {
	_, e := rl.Delete(string(lockKey))
	if e != nil {
		return e
	}
	return nil
}
