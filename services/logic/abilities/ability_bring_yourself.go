package abilities

import "nam-club/NumBuyer_back/db"

type AbilityBringYourself struct{}

func (a *AbilityBringYourself) Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error) {
	return false, nil
}
