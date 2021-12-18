package responses

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type UpdateStateResponse struct {
	Phase   string                       `json:"phase"`
	Players []UpdateStateResponsePlayers `json:"players"`
}
type UpdateStateResponsePlayers struct {
	PlayerName string `json:"playerName"`
	Coin       int    `json:"coin"`
	CardNum    int    `json:"cardNum"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateUpdateStateResponse(players []db.Player, phase consts.Phase) *UpdateStateResponse {
	ret := &UpdateStateResponse{}
	ret.Phase = phase.Value
	for _, v := range players {
		ret.Players = append(ret.Players,
			UpdateStateResponsePlayers{
				PlayerName: v.PlayerName,
				Coin:       v.Coin,
				CardNum:    len(v.Cards),
			})
	}
	return ret

}
