package requests

type CreateMatch struct {
	PlayerName string `json:"playerName" validate:"required,min=1,max=20"`
}
