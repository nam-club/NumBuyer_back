package requests

type GameBid struct {
	RoomID   string `json:"roomId" validate:"required,min=8,max=16"`
	PlayerID string `json:"playerId" validate:"required"`
	Coin     int    `json:"coin" validate:"min=0,max=1000000"`
	Action   string `json:"action" validate:"required"`
}
