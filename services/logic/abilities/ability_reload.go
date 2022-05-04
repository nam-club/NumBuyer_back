package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/utils"
)

type AbilityReload struct{}

func (a *AbilityReload) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityReload) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if game.State.Phase == consts.PhaseAuction.Value {
		return false, nil, orgerrors.NewValidationError("ability can not fire when phase is aunction")
	}
	if me.Abilities[abilityIndex].Remaining == 0 {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUsed)
	} else {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUnused)
	}
	me.Coin = me.Coin / 2
	me.Cards = utils.GenerateRandomCard(len(me.Cards))
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}
}