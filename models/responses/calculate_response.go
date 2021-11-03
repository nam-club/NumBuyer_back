package responses

type CalculateResponse struct {
	IsCorrectAnswer bool     `json:"isCorrectAnswer"`
	PlayerID        string   `json:"playerId"`
	Coin            int      `json:"coin"`
	Cards           []string `json:"cards"`
}
