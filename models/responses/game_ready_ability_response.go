package responses

type GameReadyAbilityResponse struct {
	Status    string `json:"status"`
	Remaining int    `json:"remaining"`
}
