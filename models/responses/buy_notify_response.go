package responses

type BuyNotifyResponse struct {
	PlayerName  string `json:"playerName"`
	Coin        int    `json:"coin"`
	AuctionCard string `json:"auctionCard"`
	IsPassAll   bool   `json:"isPassAll"`
}
