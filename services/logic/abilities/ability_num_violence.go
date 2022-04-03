package abilities

import "nam-club/NumBuyer_back/db"

type AbilityNumViolence struct{}

func (a *AbilityNumViolence) CanActivate(game *db.Game, player *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityNumViolence) Fire(game *db.Game, player *db.Player, abilityIndex int) (bool, error) {
	return false, nil
}
