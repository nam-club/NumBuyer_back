package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/services/logic/abilities"
	"nam-club/NumBuyer_back/utils"

	"go.uber.org/zap"
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
	_, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}
	return ret, nil
}

func TryActivateAbilitiesIfHave(game *db.Game, abilityId string) (err error) {
	players, err := db.GetPlayers(game.RoomID)
	if err != nil {
		return err
	}

	for _, p := range players {
		if err = TryActivateAbilityIfHave(game, &p, abilityId); err != nil {
			return err
		}
	}
	return nil
}
func TryActivateAbilityIfHave(game *db.Game, player *db.Player, abilityId string) (err error) {
	abilityIndex := -1
	for i, a := range player.Abilities {
		if a.ID == abilityId {
			abilityIndex = i
			break
		}
	}
	if abilityIndex == -1 {
		return nil
	}

	var ability abilities.Ability
	switch abilityId {
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
		return orgerrors.NewValidationError("ability parse error. " + abilityId)
	}

	enable := false
	if enable, err = ability.CanActivate(game, player, &player.Abilities[abilityIndex]); enable {
		player.Abilities[abilityIndex].Status = string(consts.AbilityStatusActive)
		if _, err = db.SetPlayer(game.RoomID, player); err != nil {
			return err
		}
	}
	return err
}

func FireAbility(game *db.Game, player *db.Player) ([]*db.Ability, error) {
	firedAbilities := []*db.Ability{}
	for i, ab := range player.Abilities {
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
		if fired, firedAbility, err := ability.Fire(game, player, i); fired {
			firedAbilities = append(firedAbilities, firedAbility)
		} else if err != nil {
			utils.Log.Error("ability fire failed", zap.String("error", err.Error()))
		}
	}
	return firedAbilities, nil
}
