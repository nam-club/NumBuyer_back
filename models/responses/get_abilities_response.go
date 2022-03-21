package responses

import "nam-club/NumBuyer_back/consts"

type GetAbilitiesResponse struct {
	Abilities []GetAbilitiesResponseAbilities `json:"abilities"`
}
type GetAbilitiesResponseAbilities struct {
	AbilityId string `json:"abilityId"`
	Type      string `json:"type"`
	Trigger   string `json:"trigger"`
}

func GenerateGetAbilitiesResponse(abilities []consts.Ability) *GetAbilitiesResponse {
	var respAbilities []GetAbilitiesResponseAbilities
	for _, v := range abilities {
		respAbilities = append(respAbilities, GetAbilitiesResponseAbilities{v.ID, string(v.Type), string(v.Trigger)})
	}
	return &GetAbilitiesResponse{Abilities: respAbilities}
}
