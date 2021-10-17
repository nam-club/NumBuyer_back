package responses

import "nam-club/NumBuyer_back/db"

type NextTurnResponse struct {
	PlayerID    string   `json:"playerId"`
	Cards       []string `json:"cards"`
	Coin        int      `json:"coin"`
	TargetCard  string   `json:"targetCard"`
	AuctionCard string   `json:"auctionCard"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateNextTurnResponse(player db.Player, game db.Game) *NextTurnResponse {
	return &NextTurnResponse{
		PlayerID:   player.PlayerID,
		Cards:      player.Cards,
		Coin:       player.Coin,
		TargetCard: game.State.Answer,
	}

}
