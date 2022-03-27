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
	PlayerId        string   `json:"playerId"`
	PlayerName      string   `json:"playerName"`
	Coin            int      `json:"coin"`
	CardNum         int      `json:"cardNum"`
	FiredAbilityIds []string `json:"firedAbilityIds"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateUpdateStateResponse(players []db.Player, phase consts.Phase) *UpdateStateResponse {

	ret := &UpdateStateResponse{}
	ret.Phase = phase.Value
	for _, p := range players {
		abilityIds := []string{}
		for _, a := range p.Abilities {
			abilityIds = append(abilityIds, a.ID)
		}

		ret.Players = append(ret.Players,
			UpdateStateResponsePlayers{
				PlayerId:        p.PlayerID,
				PlayerName:      p.PlayerName,
				Coin:            p.Coin,
				CardNum:         len(p.Cards),
				FiredAbilityIds: abilityIds,
			})
	}
	return ret

}
