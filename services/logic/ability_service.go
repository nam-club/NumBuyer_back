package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
)

// アビリティを発動準備状態にする
func ReadyAbilities(roomId, playerId string, abilityIds []string) (string, error) {

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return "", e
	}

	for i, v := range player.Abilities {
		for _, inV := range abilityIds {
			if v.ID == inV {
				if v.Status != string(consts.AbilityStatusUnused) {
					return "", orgerrors.NewValidationError("ability status is not unused")
				}
				player.Abilities[i].Status = string(consts.AbilityStatusReady)
			}
		}
	}
	fmt.Printf("%v\n", player)

	_, e = db.SetPlayer(roomId, player)
	if e != nil {
		return "", e
	}
	return string(consts.AbilityStatusReady), nil
}
