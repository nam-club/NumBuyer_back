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
	abilityReload      = new(abilities.AbilityReload)
	abilityShutdown    = new(abilities.AbilityShutdown)
	abilityCatastrophe = new(abilities.AbilityCatastrophe)
)

// アビリティのステータスを発動状態・発動可能状態にする
func ReadyAbility(roomId, playerId string, abilityId string) (*db.Ability, error) {

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	// パラメータのセット
	var ret *db.Ability
	var abilityDB *db.Ability
	var abilityDBIndex int
	var abilityConst *consts.Ability
	for i, v := range player.Abilities {
		if v.ID == abilityId {
			ab, e := consts.ParseAbility(v.ID)
			if e != nil {
				return nil, orgerrors.NewValidationError("invalid ability id")
			}
			abilityDB = &v
			abilityConst = &ab
			abilityDBIndex = i
			break
		}
	}

	if abilityDB == nil || abilityConst == nil {
		return nil, orgerrors.NewValidationError("invalid ability")
	}

	if abilityConst.Timing == consts.AbilityTimingSoon {
		// 発動タイミングがsoonならアビリティを発動する
		game, e := db.GetGame(roomId)
		if e != nil {
			return nil, e
		}

		ab, e := FireAbility(game, player, abilityConst.ID)
		return ab, e
	} else if abilityConst.Timing == consts.AbilityTimingWait {
		// 発動タイミングがwaitならステータスをreadyにする
		if abilityDB.Status != string(consts.AbilityStatusUnused) {
			return nil, orgerrors.NewValidationError("ability status is not unused")
		}
		if player.Abilities[abilityDBIndex].Remaining == 0 {
			return nil, orgerrors.NewValidationError("exceeded the number of ability usable")
		}
		player.Abilities[abilityDBIndex].Status = string(consts.AbilityStatusReady)
		player.Abilities[abilityDBIndex].Remaining = player.Abilities[abilityDBIndex].Remaining - 1
		ret = &player.Abilities[abilityDBIndex]
		_, e = db.SetPlayer(roomId, player)
		if e != nil {
			return nil, e
		}
		return ret, nil
	}

	return nil, orgerrors.NewInternalServerError("unexpected path")
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
	case consts.AbilityIdReload:
		ability = abilityReload
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

// 発動条件を満たしているアビリティを発動する
func FireAbilities(game *db.Game, player *db.Player) ([]*db.Ability, error) {
	firedAbilities := []*db.Ability{}
	for i, ab := range player.Abilities {
		var ability abilities.Ability
		switch ab.ID {
		case consts.AbilityIdFiBoost:
			ability = abilityFiBoost
		case consts.AbilityIdNumViolence:
			ability = abilityNumViolence
		case consts.AbilityIdReload:
			ability = abilityReload
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

// 発動条件を満たしているアビリティを発動する
func FireAbility(game *db.Game, player *db.Player, abilityId string) (*db.Ability, error) {
	for i, ab := range player.Abilities {
		if ab.ID != abilityId {
			continue
		}

		var ability abilities.Ability
		switch ab.ID {
		case consts.AbilityIdFiBoost:
			ability = abilityFiBoost
		case consts.AbilityIdNumViolence:
			ability = abilityNumViolence
		case consts.AbilityIdReload:
			ability = abilityReload
		case consts.AbilityIdShutdown:
			ability = abilityShutdown
		case consts.AbilityIdCatastrophe:
			ability = abilityCatastrophe
		default:
			return nil, orgerrors.NewValidationError("ability parse error. " + ab.ID)
		}
		if fired, firedAbility, err := ability.Fire(game, player, i); fired {
			return firedAbility, nil
		} else if err != nil {
			return nil, err
		} else if !fired {
			return nil, nil
		}

	}
	// 通常呼び出されないパス
	return nil, orgerrors.NewValidationError("ability not found")
}
