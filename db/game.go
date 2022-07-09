// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"unsafe"

	"github.com/pkg/errors"
)

type Game struct {
	RoomID     string `json:"roomId"`
	PlayersMin int    `json:"playersMin"`
	PlayersMax int    `json:"playersMax"`
	State      State  `json:"state"`
	CreatedAt  string `json:"createdAt"`
}
type State struct {
	CurrentTurn            int      `json:"currentTurn"`
	Phase                  string   `json:"phase"`
	Auction                []string `json:"auction"`
	AuctionMaxBid          string   `json:"auctionMaxBid"`
	AuctionLastBidPlayerId string   `json:"auctionLastBidPlayer"`
	SkipShowTarget         bool     `json:"skipShowTarget"`
	Answer                 string   `json:"answer"`
	PhaseChangedTime       string   `json:"phaseChangedTime"`
}

var rg *RedisHandler

func init() {
	rg = NewRedisHandler( /*index=*/ 0)
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

func ScanGame(iter int) (int, []string, error) {
	return rg.Scan(iter)
}
