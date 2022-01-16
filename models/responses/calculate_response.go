package responses

type CalculateResponse struct {
	IsCorrectAnswer bool      `json:"isCorrectAnswer"`
	PlayerID        string    `json:"playerId"`
	Coin            int       `json:"coin"`
	Cards           []string  `json:"cards"`
	AddedCoin       AddedCoin `json:"addedCoin"`
}

type AddedCoin struct {
	Total        int `json:"total"`
	CardNumBonus int `json:"cardNumBonus"`
}
