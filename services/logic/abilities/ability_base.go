package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type Ability interface {
	Fire(game *db.Game, player *db.Player, ability *db.Ability) (bool, error)
}

func IsReady(ability *db.Ability) bool {
	return ability.Status == string(consts.AbilityStatusReady)
}

func IsUsed(ability *db.Ability) bool {
	return ability.Status == string(consts.AbilityStatusUsed)
}
