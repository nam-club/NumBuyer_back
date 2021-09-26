package responses

type Player struct {
	PlayerID   int    `json:"playerId"`
	PlayerName string `json:"playerName"`
	RoomID     string `json:"roomId"`
	Money      int    `json:"money"`
	Cards      struct {
	} `json:"cards"`
}
