// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

type Game struct {
	RoomID string `json:"roomId"`
	State  State  `json:"state"`
}
type State struct {
	Phase       string `json:"phase"`
	Auction     string `json:"auction"`
	Answer      string `json:"answer"`
	ChangedTime string `json:"changedTime"`
}

var rg *RedisHandler

func init() {
	rg = NewRedisHandler( /*index=*/ 0)
}

// ゲーム情報をセット
func SetGame(id string, game *Game) (*Game, error) {
	// 変更時間を更新する
	game.State.ChangedTime = time.Now().Format(time.RFC3339)

	j, e := json.Marshal(game)
	if e != nil {
		return nil, errors.WithStack(e)
	}
	// byteからstringに変換
	str := *(*string)(unsafe.Pointer(&j))
	_, e = rg.Set(id, str)
	if e != nil {
		return nil, e
	}

	return game, nil
}

// ゲーム情報を取得
func GetGame(id string) (*Game, error) {
	r, e := rg.Get(id)
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
	return rg.Delete(id)
}

// ゲームが存在するかをチェック
func ExistsGame(id string) (bool, error) {
	r, e := rg.Exists(id)
	if e != nil {
		return false, e
	}
	return r, nil
}
