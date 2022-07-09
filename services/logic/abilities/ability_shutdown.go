package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

const (
	ShutdownMinCardsNum = 5
)

type AbilityShutdown struct{}

func (a *AbilityShutdown) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return true, nil
}

// shutdownの効果はスケジューラが発動するため、ここではステータスの更新だけ行う
func (a *AbilityShutdown) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if len(me.AnswerAction.Cards) < ShutdownMinCardsNum {
		return false, nil, nil
	}

	if me.Abilities[abilityIndex].Remaining == 0 {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUsed)
	} else {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusBackToReady)
	}
	me.ForceReady = true
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}
}
