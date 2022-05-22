package requests

type GameReadyAbility struct {
	RoomID    string `json:"roomId" validate:"required,min=8,max=16"`
	PlayerID  string `json:"playerId" validate:"required"`
	AbilityId string `json:"abilityId"`
}
