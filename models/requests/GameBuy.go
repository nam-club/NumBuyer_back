package requests

type GameBuy struct {
	PlayerID string `json:"playerId" validate:"required"`
}
