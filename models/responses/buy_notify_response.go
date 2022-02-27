package responses

type BuyNotifyResponse struct {
	PlayerName   string   `json:"playerName"`
	Coin         int      `json:"coin"`
	AuctionCards []string `json:"auctionCards"`
	IsPassAll    bool     `json:"isPassAll"`
}
