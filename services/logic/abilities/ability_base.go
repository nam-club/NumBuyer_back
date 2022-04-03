package abilities

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

// 以下の流れでアビリティが発動する
// 1. アビリティ毎適切なタイミングでActivateする
// 2. スケジューラがActiveなアビリティを検知して発動
type Ability interface {
	CanActivate(game *db.Game, player *db.Player, targetAbility *db.Ability) (bool, error)
	Fire(game *db.Game, player *db.Player, abilityIndex int) (bool, error)
}

func IsReady(ability *db.Ability) bool {
	return ability.Status == string(consts.AbilityStatusReady)
}

func IsActive(ability *db.Ability) bool {
	return ability.Status == string(consts.AbilityStatusActive)
}

func IsUsed(ability *db.Ability) bool {
	return ability.Status == string(consts.AbilityStatusUsed)
}
