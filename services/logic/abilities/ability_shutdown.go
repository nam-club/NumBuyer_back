package abilities

import "nam-club/NumBuyer_back/db"

type AbilityShutdown struct{}

func (a *AbilityShutdown) CanActivate(game *db.Game, player *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityShutdown) Fire(game *db.Game, player *db.Player, abilityIndex int) (bool, error) {
	return false, nil
}
