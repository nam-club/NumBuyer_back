package abilities

import (
	"nam-club/NumBuyer_back/db"
)

const (
	ShutdownThresholdSeconds = 5
)

type AbilityShutdown struct{}

func (a *AbilityShutdown) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return true, nil
}

// shutdownの効果はスケジューラが発動するため、ここではステータスの更新だけ行う
func (a *AbilityShutdown) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	println("set force ready true")
	me.ForceReady = true
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}
}
