package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"time"
)

const (
	ThresholdSeconds = 5
	BonusCoin        = 5
)

type AbilityFiBoost struct{}

func (a *AbilityFiBoost) CanActivate(game *db.Game, player *db.Player, targetAbility *db.Ability) (bool, error) {
	if IsUsed(targetAbility) ||
		!IsReady(targetAbility) ||
		game.State.Phase != consts.PhaseCalculate.Value ||
		!player.AnswerAction.Correct {
		return false, nil
	}
	phaseChangedTime, e := time.Parse(time.RFC3339, game.State.PhaseChangedTime)
	if e != nil {
		return false, e
	}
	answerTime, e := time.Parse(time.RFC3339, player.AnswerAction.AnswerTime)
	if e != nil {
		return false, e
	}
	if answerTime.Before(phaseChangedTime.Add(time.Duration(ThresholdSeconds) * time.Second)) {
		return true, nil
	} else {
		return false, nil
	}
}

func (a *AbilityFiBoost) Fire(game *db.Game, player *db.Player, abilityIndex int) (bool, error) {
	if !IsActive(&player.Abilities[abilityIndex]) {
		return false, nil
	}
	player.Abilities[abilityIndex].Status = string(consts.AbilityStatusReady)
	player.Coin += BonusCoin
	if _, e := db.SetPlayer(game.RoomID, player); e != nil {
		return false, e
	} else {
		return true, nil
	}
}
