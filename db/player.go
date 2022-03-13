// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"nam-club/NumBuyer_back/models/orgerrors"
	"unsafe"

	"github.com/pkg/errors"
)

type Player struct {
	PlayerID     string       `json:"playerId"`
	PlayerName   string       `json:"playerName"`
	IsOwner      bool         `json:"isOwner"`
	Coin         int          `json:"coin"`
	Cards        []string     `json:"cards"`
	Abilities    []Ability    `json:"abilities"`
	BuyAction    BuyAction    `json:"buyAction"`
	AnswerAction AnswerAction `json:"answerAction"`
	Ready        bool         `json:"ready"`
}

type Ability struct {
	ID     string `json:"id"`
	Status string `json:"status"` // 実行状態
	Remain int    `json:"remain"` // 残使用回数 -1なら無限に実行可能
}
type BuyAction struct {
	Action   string `json:"action"`
	Value    string `json:"value"`
	BidCount int    `json:"bidCount"`
}
type AnswerAction struct {
	Action     string   `json:"action"`
	Cards      []string `json:"cards"`
	PlusCoin   int      `json:"plusCoin"`
	AnswerTime string   `json:"answerTime"`
	Correct    bool     `json:"correct"`
}

var rp *RedisHandler

func init() {
	rp = NewRedisHandler( /*index=*/ 1)
}

// プレイヤー情報一覧を取得
func GetPlayers(roomId string) ([]Player, error) {
	r, e := rp.HVals(roomId)
	if e != nil {
		return []Player{}, e
	}

	var ret []Player
	for _, v := range r {
		var player Player
		if e := json.Unmarshal(v, &player); e != nil {
			return []Player{}, errors.WithStack(e)
		}
		ret = append(ret, player)
	}
	return ret, nil
}

// プレイヤー情報を取得
func GetPlayer(roomId, playerId string) (*Player, error) {
	r, e := rp.HGet(roomId, playerId)
	if e != nil {
		return nil, e
	}

	var ret *Player
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return nil, errors.WithStack(e)
	}
	return ret, nil
}

// プレイヤー情報を追加
func SetPlayer(roomId string, player *Player) (*Player, error) {
	if b, e := ExistsGame(roomId); e != nil || b == false {
		if e != nil {
			return nil, errors.WithStack(e)
		}
		return nil, orgerrors.NewGameNotFoundError("")
	}

	b, _ := json.Marshal(player)
	str := *(*string)(unsafe.Pointer(&b)) // byteからstringに変換
	if _, e := rp.HSet(roomId, player.PlayerID, str); e != nil {
		return nil, e
	}

	if _, e := rp.HGet(roomId, player.PlayerID); e != nil {
		return nil, e
	}
	return player, nil
}

// プレイヤー情報を削除
func DeletePlayers(id string) (int, error) {
	return rp.Delete(id)
}
