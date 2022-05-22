package responses

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type JoinResponse struct {
	RoomID    string                  `json:"roomId"`
	PlayerID  string                  `json:"playerId"`
	IsOwner   bool                    `json:"isOwner"`
	Abilities []JoinResponseAbilities `json:"abilities"`
}

type JoinResponseAbilities struct {
	AbilityId string `json:"abilityId"`
	Status    string `json:"status"`
	Remaining int    `json:"remaining"`
	Type      string `json:"type"`
	Trigger   string `json:"trigger"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GenerateJoinResponse(roomId string, player *db.Player) (*JoinResponse, error) {
	abilities := []JoinResponseAbilities{}
	for _, abDb := range player.Abilities {
		var abCs consts.Ability
		var err error
		if abCs, err = consts.ParseAbility(abDb.ID); err != nil {
			return nil, err
		}
		abilities = append(abilities, JoinResponseAbilities{abDb.ID, abDb.Status, abDb.Remaining, string(abCs.Type), string(abCs.Trigger)})
	}
	return &JoinResponse{RoomID: roomId, PlayerID: player.PlayerID, IsOwner: player.IsOwner, Abilities: abilities}, nil

}
