package responses

type Player struct {
	PlayerID   string   `json:"playerId"`
	PlayerName string   `json:"playerName"`
	RoomID     string   `json:"roomId"`
	Money      int      `json:"money"`
	Cards      []string `json:"cards"`
}
