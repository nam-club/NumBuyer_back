// ゲーム情報に関するDB操作
package db

import (
	"encoding/json"
	"unsafe"
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
type BuyAction struct {
	Action string `json:"action"`
	Value  string `json:"value"`
}
type AnswerAction struct {
	Action  string `json:"action"`
	CardIds []int  `json:"cardIds"`
}
type Player struct {
	PlayerID     int          `json:"playerId"`
	PlayerName   string       `json:"playerName"`
	Coin         int          `json:"coin"`
	Cards        []string     `json:"cards"`
	BuyAction    BuyAction    `json:"buyAction"`
	AnswerAction AnswerAction `json:"answerAction"`
}

// トランザクションを実行する
func Transaction(f func()) {
	Atomic(f)
}

// ゲーム情報をセット
func SetGame(id string, game Game) (string, error) {
	j, e := json.Marshal(game)
	if e != nil {
		return "", e
	}
	// byteからstringに変換
	str := *(*string)(unsafe.Pointer(&j))
	ret := Set(id, str)
	return ret, nil
}

// ゲーム情報を取得
func GetGame(id string) (Game, error) {
	r := Get(id)
	var ret Game
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return Game{}, e
	}
	return ret, nil
}

// プレイヤー情報を取得
func GetPlayers(id string) ([]Player, error) {
	r := Get(id)
	var ret Game
	if e := json.Unmarshal([]byte(r), &ret); e != nil {
		return []Player{}, e
	}
	return ret.Players, nil
}

// プレイヤー情報をセット
func SetPlayer(id string, player Player) (Player, error) {
	g, err := GetGame(id)
	if err != nil {
		return Player{}, err
	}
	g.Players = append(g.Players, player)

	str := *(*string)(unsafe.Pointer(&g)) // byteからstringに変換
	Set(id, str)
	return player, nil
}
