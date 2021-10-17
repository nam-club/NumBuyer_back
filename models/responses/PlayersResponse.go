package responses

type PlayersResponse struct {
	Players []Players `json:"players"`
}
type Players struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	RoomID     string `json:"roomId"`
	Coin       int    `json:"coin"`
	CardNum    int    `json:"cardNum"`
}
