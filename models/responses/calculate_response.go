package responses

type CalculateResponse struct {
	ActionResult string   `json:"actionResult"`
	PlayerID     string   `json:"playerId"`
	Coin         int      `json:"coin"`
	Cards        []string `json:"cards"`
}
