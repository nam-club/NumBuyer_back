// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"unsafe"

	"github.com/pkg/errors"
)

type Game struct {
	RoomID  string   `json:"roomId"`
	State   State    `json:"state"`
	Players []Player `json:"players"`
}
type State struct {
	Phase   string `json:"phase"`
	Auction string `json:"auction"`
	Answer  int    `json:"answer"`
}

var rg *RedisHandler

func init() {
	rg = NewRedisHandler(0)
}

// ゲーム情報をセット
func SetGame(id string, game Game) (string, error) {
	j, e := json.Marshal(game)
	if e != nil {
		return "", errors.WithStack(e)
	}
	// byteからstringに変換
	str := *(*string)(unsafe.Pointer(&j))
	ret, e := rg.Set(id, str)
	if e != nil {
		return "", e
	}
	return ret, nil
}

// ゲーム情報を取得
func GetGame(id string) (Game, error) {
	r, e := rg.Get(id)
	if e != nil {
		return Game{}, e
	}

	var ret Game
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return Game{}, errors.WithStack(e)
	}
	return ret, nil
}

// ゲーム情報を取得
func GetRandomGameId() (string, error) {
	r, e := rg.RandomKey()
	if e != nil {
		return "", e
	}

	return r, nil
}

// ゲームが存在するかをチェック
func ExistsGame(id string) (bool, error) {
	r, e := rg.Exists(id)
	if e != nil {
		return false, e
	}
	return r, nil
}
