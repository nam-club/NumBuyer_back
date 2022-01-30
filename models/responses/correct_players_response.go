package responses

type CorrectPlayersResponse struct {
	AnsPlayers    []CorrectPlayer `json:"ansPlayers"`
	ExistsCorrect bool            `json:"existsCorrect"`
}

type CorrectPlayer struct {
	PlayerName string    `json:"playerName"`
	AddedCoin  AddedCoin `json:"addedCoin"`
}

type AddedCoin struct {
	Total        int `json:"total"`
	CardNumBonus int `json:"cardNumBonus"`
}
