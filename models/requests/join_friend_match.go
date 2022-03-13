package requests

type JoinFriendMatch struct {
	PlayerName string   `json:"playerName" validate:"required,min=1,max=20"`
	RoomID     string   `json:"roomId" validate:"required,min=8,max=16"`
	AbilityIds []string `json:"abilityIds"`
}
