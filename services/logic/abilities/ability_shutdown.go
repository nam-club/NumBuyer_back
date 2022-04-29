package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

const (
	ShutdownThresholdSeconds = 5
)

type AbilityShutdown struct{}

func (a *AbilityShutdown) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	if IsUsed(targetAbility) ||
		!IsReady(targetAbility) ||
		game.State.Phase != consts.PhaseCalculate.Value {
		return false, nil
	}
	return true, nil
}

// shutdownの効果はスケジューラが発動するため、ここではステータスの更新だけ行う
func (a *AbilityShutdown) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if !IsActive(&me.Abilities[abilityIndex]) {
		return false, nil, nil
	}
	me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUnused)
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}
}
