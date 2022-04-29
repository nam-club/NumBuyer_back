package abilities

import "nam-club/NumBuyer_back/db"

type AbilityCatastrophe struct{}

func (a *AbilityCatastrophe) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityCatastrophe) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	return false, nil, nil
}
