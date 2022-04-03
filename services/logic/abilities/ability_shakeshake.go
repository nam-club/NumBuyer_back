package abilities

import "nam-club/NumBuyer_back/db"

type AbilityShakeShake struct{}

func (a *AbilityShakeShake) Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error) {
	return false, nil
}
