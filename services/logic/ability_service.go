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
	abilityFiBoost     = new(abilities.AbilityFiBoost)
	abilityNumViolence = new(abilities.AbilityNumViolence)
	abilityReboot      = new(abilities.AbilityReboot)
	abilityShutdown    = new(abilities.AbilityShutdown)
	abilityCatastrophe = new(abilities.AbilityCatastrophe)
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
			if player.Abilities[i].Remaining == 0 {
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

func IsExistsActive(roomId string, abilityId string) (bool, error) {
	players, err := db.GetPlayers(roomId)
	if err != nil {
		return false, err
	}

	for _, p := range players {
		for _, a := range p.Abilities {
			if a.ID == abilityId {
				if abilities.IsActive(&a) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func TryActivateAbilitiesIfHave(game *db.Game, abilityId string) (err error) {
	players, err := db.GetPlayers(game.RoomID)
	if err != nil {
		return err
	}

	for _, p := range players {
		if _, err = TryActivateAbilityIfHave(game, &p, abilityId); err != nil {
			return err
		}
	}
	return nil
}
func TryActivateAbilityIfHave(game *db.Game, player *db.Player, abilityId string) (canActivate bool, err error) {
	abilityIndex := -1
	for i, a := range player.Abilities {
		if a.ID == abilityId {
			abilityIndex = i
			break
		}
	}
	if abilityIndex == -1 {
		return false, nil
	}

	var ability abilities.Ability
	switch abilityId {
	case consts.AbilityIdFiBoost:
		ability = abilityFiBoost
	case consts.AbilityIdNumViolence:
		ability = abilityNumViolence
	case consts.AbilityIdReboot:
		ability = abilityReboot
	case consts.AbilityIdShutdown:
		ability = abilityShutdown
	case consts.AbilityIdCatastrophe:
		ability = abilityCatastrophe
	default:
		return false, orgerrors.NewValidationError("ability parse error. " + abilityId)
	}

	enable := false
	if enable, _ = ability.CanActivate(game, player, &player.Abilities[abilityIndex]); enable {
		player.Abilities[abilityIndex].Status = string(consts.AbilityStatusActive)
		if _, err = db.SetPlayer(game.RoomID, player); err != nil {
			return false, err
		}
	}
	return true, nil
}

func HaveAbility(player *db.Player, abilityId string) bool {
	for _, a := range player.Abilities {
		if a.ID == abilityId {
			return true
		}
	}
	return false
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
		case consts.AbilityIdReboot:
			ability = abilityReboot
		case consts.AbilityIdShutdown:
			ability = abilityShutdown
		case consts.AbilityIdCatastrophe:
			ability = abilityCatastrophe
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
