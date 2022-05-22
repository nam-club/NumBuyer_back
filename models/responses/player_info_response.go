package responses

import "nam-club/NumBuyer_back/db"

type PlayerInfoResponse struct {
	PlayerId string   `json:"playerId"`
	Coin     int      `json:"coin"`
	Cards    []string `json:"cards"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GeneratePlayerInfoResponse(player db.Player) *PlayerInfoResponse {
	ret := &PlayerInfoResponse{}
	ret.PlayerId = player.PlayerID
	ret.Coin = player.Coin
	ret.Cards = player.Cards
	return ret

}
