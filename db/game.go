// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"nam-club/NumBuyer_back/models/orgerrors"
	"unsafe"

	"github.com/pkg/errors"
)

type Game struct {
	RoomID  string   `json:"roomId"`
	State   State    `json:"state"`
	Players []Player `json:"players"`
}
type State struct {
	Phase       string `json:"phase"`
	Auction     string `json:"auction"`
	Answer      string `json:"answer"`
	ChangedTime string `json:"changedTime"`
}

var rg *RedisHandler

func init() {
	rg = NewRedisHandler(0)
}

// ゲーム情報をセット
func SetGame(id string, game *Game) (*Game, error) {
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

//ランダムな部屋IDを取得
func GetRandomRoomId() (string, error) {
	l, e := rg.DBSize()
	if e != nil {
		return "", e
	}
	if l < 1 {
		return "", orgerrors.NewGameNotFoundError("")
	}

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
