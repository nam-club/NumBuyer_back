package abilities

import "nam-club/NumBuyer_back/db"

type AbilityFiBoost struct{}

func (a *AbilityFiBoost) Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error) {
	}
	return false
}

func (a *AbilityFiBoost) Fire(game *db.Game, player *db.Player, ability *db.Ability) {

}
