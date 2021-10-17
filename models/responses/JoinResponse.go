package responses

type JoinResponse struct {
	RoomID   string `json:"roomId"`
	PlayerID string `json:"playerId"`
}
