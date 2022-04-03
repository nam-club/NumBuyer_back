package abilities

import "nam-club/NumBuyer_back/db"

type AbilityShutdown struct{}

func (a *AbilityShutdown) Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error) {
	return false, nil
}
