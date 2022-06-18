// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"

	"github.com/pkg/errors"
)

var rgl *RedisHandler

func init() {
	rgl = NewRedisHandler( /*index=*/ 3)
}

// ゲームのロック情報を取得
func SetGameLock(id string) (*Game, error) {
	r, e := rgl.Get(id)
	if e != nil {
		return nil, e
	}

	var ret *Game
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return nil, errors.WithStack(e)
	}
	return ret, nil
}

// ゲーム情報を削除
func DeleteGame(id string) (int, error) {
	return rgl.Delete(id)
}

// ゲームが存在するかをチェック
func ExistsGame(id string) (bool, error) {
	r, e := rgl.Exists(id)
	if e != nil {
		return false, e
	}
	return r, nil
}

func ScanGame(iter int) (int, []string, error) {
	return rgl.Scan(iter)
}
