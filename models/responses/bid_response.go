package responses

type BidResponse struct {
	PlayerName    string `json:"playerName"`
	Coin          int    `json:"coin"`
	RemainingTime int    `json:"remainingTime"`
}
