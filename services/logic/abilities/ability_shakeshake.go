package abilities

import "nam-club/NumBuyer_back/db"

type AbilityShakeShake struct{}

func (a *AbilityShakeShake) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	return false, nil
}

func (a *AbilityShakeShake) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	return false, nil, nil
}
