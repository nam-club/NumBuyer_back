package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
)

// アビリティを発動準備状態にする
func ReadyAbility(roomId, playerId string, abilityId string) (*db.Ability, error) {

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	var ret *db.Ability
	for i, v := range player.Abilities {
		if v.ID == abilityId {
			if v.Status != string(consts.AbilityStatusUnused) {
				return nil, orgerrors.NewValidationError("ability status is not unused")
			}
			if player.Abilities[i].Remaining <= 0 {
				return nil, orgerrors.NewValidationError("exceeded the number of ability usable")
			}
			player.Abilities[i].Status = string(consts.AbilityStatusReady)
			player.Abilities[i].Remaining = player.Abilities[i].Remaining - 1
			ret = &player.Abilities[i]
		}
	}
	fmt.Printf("%v\n", player)

	player, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}
	return ret, nil
}
