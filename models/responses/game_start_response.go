package responses

type GameStartResponse struct {
	RoomID   string `json:"roomId"`
	GoalCoin string `json:"goalCoin:"`
}
