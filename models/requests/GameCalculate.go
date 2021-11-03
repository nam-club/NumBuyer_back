package requests

type GameCalculate struct {
	PlayerID       string   `json:"playerId" validate:"required"`
	CalculateCards []string `json:"calculateCards"`
	Action         string   `json:"action" validate:"required"`
}
