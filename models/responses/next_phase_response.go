package responses

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type NextPhaseResponse struct {
	Phase   string             `json:"phase"`
	Players []NextPhasePlayers `json:"players"`
}
type NextPhasePlayers struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	Coin       int    `json:"coin"`
	CardNum    int    `json:"cardNum"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateNextPhaseResponse(players []db.Player, phase consts.Phase) *NextPhaseResponse {
	ret := &NextPhaseResponse{}
	ret.Phase = phase.Value
	for _, v := range players {
		ret.Players = append(ret.Players,
			NextPhasePlayers{
				PlayerID:   v.PlayerID,
				PlayerName: v.PlayerName,
				Coin:       v.Coin,
				CardNum:    len(v.Cards),
			})
	}
	return ret

}
