package responses

type CorrectPlayersResponse struct {
	AnsPlayers    []string `json:"ansPlayers"`
	ExistsCorrect bool     `json:"existsCorrect"`
}
