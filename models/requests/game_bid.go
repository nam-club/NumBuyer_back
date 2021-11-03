package requests

type GameBid struct {
	PlayerID string `json:"playerId" validate:"required"`
	Coin     int    `json:"coin" validate:"min=1,max=1000000"`
	Action   string `json:"action" validate:"required"`
}
