package logic

import (
	"fmt"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/services/logic/abilities"
)

var (
	abilityFiBoost       = new(abilities.AbilityFiBoost)
	abilityNumViolence   = new(abilities.AbilityNumViolence)
	abilityBringYourself = new(abilities.AbilityBringYourself)
	abilityShutdown      = new(abilities.AbilityShutdown)
	abilityShakeShake    = new(abilities.AbilityShakeShake)
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

func FireAbilityIfPossible(game *db.Game, player *db.Player) ([]string, error) {
	firedAbilityId := []string{}
	for _, ab := range player.Abilities {
		var ability abilities.Ability
		switch ab.ID {
		case consts.AbilityIdFiBoost:
			ability = abilityFiBoost
		case consts.AbilityIdNumViolence:
			ability = abilityNumViolence
		case consts.AbilityIdBringYourself:
			ability = abilityBringYourself
		case consts.AbilityIdShutdown:
			ability = abilityShutdown
		case consts.AbilityIdShakeShake:
			ability = abilityShakeShake
		default:
			return nil, orgerrors.NewValidationError("ability parse error. " + ab.ID)
		}
		if ability.IsFirable(game, player, &ab) {
			ability.Fire(game, player, &ab)
			firedAbilityId = append(firedAbilityId, ab.ID)
		}
	}
	return firedAbilityId, nil
}
