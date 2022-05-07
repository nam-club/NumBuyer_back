package responses

type BuyUpdateResponse struct {
	PlayerID    string   `json:"playerId"`
	IsSuccessed bool     `json:"isSuccessed"` // 自身が落札者か
	Coin        int      `json:"coin"`
	Cards       []string `json:"cards"`
}
