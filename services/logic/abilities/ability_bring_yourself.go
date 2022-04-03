package abilities

import "nam-club/NumBuyer_back/db"

type AbilityBringYourself struct{}

func (a *AbilityBringYourself) CanActivate(game *db.Game, player *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityBringYourself) Fire(game *db.Game, player *db.Player, abilityIndex int) (bool, error) {
	return false, nil
}
