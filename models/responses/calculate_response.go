package responses

type CalculateResponse struct {
	IsCorrectAnswer bool     `json:"isCorrectAnswer"`
	IsPassed        bool     `json:"isPassed"`
	PlayerID        string   `json:"playerId"`
	Coin            int      `json:"coin"`
	Cards           []string `json:"cards"`
}
