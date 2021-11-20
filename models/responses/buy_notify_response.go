package responses

type BuyNotifyResponse struct {
	PlayerName string `json:"playerName"`
	PlayerID   string `json:"playerId"`
	Coin       int    `json:"coin"`
}
