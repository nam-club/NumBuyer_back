package requests

type GamePlayerInfo struct {
	PlayerID string `json:"playerId" validate:"required"`
	RoomID   string `json:"roomId" validate:"required,min=8,max=16"`
}
