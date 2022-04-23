package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"time"
)

const (
	ShutdownThresholdSeconds = 5
)

type AbilityShutdown struct{}

func (a *AbilityShutdown) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	if IsUsed(targetAbility) ||
		!IsReady(targetAbility) ||
		game.State.Phase != consts.PhaseCalculate.Value ||
		!me.AnswerAction.Correct {
		return false, nil
	}
	phaseChangedTime, e := time.Parse(time.RFC3339, game.State.PhaseChangedTime)
	if e != nil {
		return false, e
	}
	answerTime, e := time.Parse(time.RFC3339, me.AnswerAction.AnswerTime)
	if e != nil {
		return false, e
	}
	if answerTime.Before(phaseChangedTime.Add(time.Duration(ShutdownThresholdSeconds) * time.Second)) {
		return true, nil
	} else {
		return false, nil
	}
}

// shutdownの効果はスケジューラが発動するため、ここではステータスの更新だけ行う
func (a *AbilityShutdown) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if !IsActive(&me.Abilities[abilityIndex]) {
		return false, nil, nil
	}
	me.Abilities[abilityIndex].Status = string(consts.AbilityStatusReady)
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}
}
