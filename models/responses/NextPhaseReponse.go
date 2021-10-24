package responses

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type NextPhaseResponse struct {
	Phase   string             `json:"phase"`
	Players []NextPhasePlayers `json:"playerList"`
}
type NextPhasePlayers struct {
	PlayerID string `json:"playerId"`
	Coin     int    `json:"coin"`
	CardNum  int    `json:"cardNum"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateNextPhaseResponse(players []db.Player, phase consts.Phase) *NextPhaseResponse {
	ret := &NextPhaseResponse{}
	ret.Phase = string(phase)
	for _, v := range players {
		ret.Players = append(ret.Players,
			NextPhasePlayers{
				PlayerID: v.PlayerID,
				Coin:     v.Coin,
				CardNum:  len(v.Cards),
			})
	}
	return ret

}