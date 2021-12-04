package responses

type GameStartResponse struct {
	RoomID   string `json:"roomId"`
	GoalCoin int    `json:"goalCoin:"`
}
