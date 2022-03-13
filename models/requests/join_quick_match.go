package requests

type JoinQuickMatch struct {
	PlayerName string   `json:"playerName" validate:"required,min=1,max=20"`
	AbilityIds []string `json:"abilityIds"`
}
