package responses

type GameReadyAbilityResponse struct {
	AbilityId string `json:"abilityId"`
	Status    string `json:"status"`
	Remaining int    `json:"remaining"`
}
