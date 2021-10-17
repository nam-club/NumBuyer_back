// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"nam-club/NumBuyer_back/models/orgerrors"
	"strings"
	"unsafe"

	"github.com/pkg/errors"
)

type Player struct {
	PlayerID     string       `json:"playerId"`
	PlayerName   string       `json:"playerName"`
	Coin         int          `json:"coin"`
	Cards        []string     `json:"cards"`
	BuyAction    BuyAction    `json:"buyAction"`
	AnswerAction AnswerAction `json:"answerAction"`
}

type BuyAction struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}
type AnswerAction struct {
	Action  string `json:"action"`
	CardIds []int  `json:"cardIds"`
}

var rp *RedisHandler

func init() {
	rp = NewRedisHandler(1)
}

// プレイヤー情報一覧を取得
func GetPlayers(gameId string) ([]Player, error) {
	r, e := rp.HVals(gameId)
	if e != nil {
		return []Player{}, e
	}

	players := strings.Split(r, ",")

	var ret []Player
	for _, v := range players {
		var player Player
		if e := json.Unmarshal([]byte(v), &player); e != nil {
			return []Player{}, errors.WithStack(e)
		}
		ret = append(ret, player)
	}
	return ret, nil
}

// プレイヤー情報を取得
func GetPlayer(gameId, playerId string) (Player, error) {
	r, e := rp.HGet(gameId, playerId)
	if e != nil {
		return Player{}, e
	}

	var ret Player
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return Player{}, errors.WithStack(e)
	}
	return ret, nil
}

// プレイヤー情報を追加
func AddPlayer(gameId string, player Player) (Player, error) {
	if b, e := ExistsGame(gameId); e != nil || b == false {
		if e != nil {
			return Player{}, errors.WithStack(e)
		}
		return Player{}, orgerrors.NewGameNotFoundError("")
	}

	b, _ := json.Marshal(player)
	str := *(*string)(unsafe.Pointer(&b)) // byteからstringに変換
	if _, e := rp.HSet(gameId, player.PlayerID, str); e != nil {
		return Player{}, e
	}

	return player, nil
}
