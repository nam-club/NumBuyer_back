package responses

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type UpdateStateResponse struct {
	Phase   string                       `json:"phase"`
	Players []UpdateStateResponsePlayers `json:"players"`
}
type UpdateStateResponsePlayers struct {
	PlayerId       string                              `json:"playerId"`
	PlayerName     string                              `json:"playerName"`
	Coin           int                                 `json:"coin"`
	CardNum        int                                 `json:"cardNum"`
	FiredAbilities []UpdateStateResponseFiredAbilities `json:"firedAbilities"`
}
type UpdateStateResponseFiredAbilities struct {
	AbilityId  string                            `json:"abilityId"`
	Status     string                            `json:"status"`
	Type       string                            `json:"type"`
	Trigger    string                            `json:"trigger"`
	Parameters []UpdateStateResponseAbilityParam `json:"parameters"`
}

type UpdateStateResponseAbilityParam struct {
	Key   string `json:"key"`
	To    string `json:"to"`
	Value string `json:"value"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
// firedAbilities: key=プレイヤーID, value=発動したアビリティの一覧
func GenerateUpdateStateResponse(players []db.Player,
	phase consts.Phase,
	firedAbilities map[string][]*db.Ability) *UpdateStateResponse {

	ret := &UpdateStateResponse{}
	ret.Phase = phase.Value
	for _, p := range players {
		fires := firedAbilities[p.PlayerID]
		respAb := []UpdateStateResponseFiredAbilities{}
		for _, firedAb := range fires {
			abilityParams := []UpdateStateResponseAbilityParam{}
			for _, param := range firedAb.Parameters {
				abilityParams = append(abilityParams, UpdateStateResponseAbilityParam{
					Key:   param.Key,
					To:    param.To,
					Value: param.Value,
				})
			}

			ab, _ := consts.ParseAbility(firedAb.ID)
			respAb = append(respAb, UpdateStateResponseFiredAbilities{
				AbilityId:  firedAb.ID,
				Status:     firedAb.Status,
				Type:       string(ab.Type),
				Trigger:    string(ab.Trigger),
				Parameters: abilityParams,
			})
		}
		ret.Players = append(ret.Players,
			UpdateStateResponsePlayers{
				PlayerId:       p.PlayerID,
				PlayerName:     p.PlayerName,
				Coin:           p.Coin,
				CardNum:        len(p.Cards),
				FiredAbilities: respAb,
			})
	}
	return ret

}
