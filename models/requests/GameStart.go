package requests

type GameStart struct {
	RoomID string `json:"roomId" validate:"required,min=8,max=16"`
}
