package responses

type BuyUpdateResponse struct {
	PlayerID string   `json:"playerId"`
	Coin     int      `json:"coin"`
	Cards    []string `json:"cards"`
}
