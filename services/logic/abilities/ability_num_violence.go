package abilities

import "nam-club/NumBuyer_back/db"

type AbilityNumViolence struct{}

func (a *AbilityNumViolence) Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error) {
	return false, nil
}
