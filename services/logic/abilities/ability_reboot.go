package abilities

import "nam-club/NumBuyer_back/db"

type AbilityReboot struct{}

func (a *AbilityReboot) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityReboot) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	return false, nil, nil
}
