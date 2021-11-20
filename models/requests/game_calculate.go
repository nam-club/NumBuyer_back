package requests

type GameCalculate struct {
	RoomID         string   `json:"roomId" validate:"required,min=8,max=16"`
	PlayerID       string   `json:"playerId" validate:"required"`
	CalculateCards []string `json:"calculateCards"`
	Action         string   `json:"action" validate:"required"`
}
