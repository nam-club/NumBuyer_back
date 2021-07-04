package responses

type Player struct {
	PlayerID   int    `json:"playerId"`
	PlayerName string `json:"playerName"`
	RoomName   string `json:"roomName"`
	Money      int    `json:"money"`
	Cards      struct {
	} `json:"cards"`
}
