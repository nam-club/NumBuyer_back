// 参加可能ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"nam-club/NumBuyer_back/models/orgerrors"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

type JoinableGame struct {
	CreateAt string `json:"createAt"`
}

var rj *RedisHandler

func init() {
	rj = NewRedisHandler( /*index=*/ 2)
}

// ゲーム情報をセット
func SetJoinableGame(roomId string) error {

	newRecord := &JoinableGame{CreateAt: time.Now().Format(time.RFC3339)}
	j, e := json.Marshal(newRecord)
	if e != nil {
		return errors.WithStack(e)
	}
	// byteからstringに変換
	str := *(*string)(unsafe.Pointer(&j))
	_, e = rj.Set(roomId, str)
	if e != nil {
		return e
	}

	return nil
}

// ゲーム情報を削除
func DeleteJoinableGame(id string) (int, error) {
	return rj.Delete(id)
}

//ランダムな部屋IDを取得
func GetRandomRoomId() (string, error) {
	l, e := rj.DBSize()
	if e != nil {
		return "", e
	}
	if l < 1 {
		return "", orgerrors.NewGameNotFoundError("")
	}

	r, e := rj.RandomKey()
	if e != nil {
		return "", e
	}

	return r, nil
}
