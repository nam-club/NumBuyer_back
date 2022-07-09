package abilities

import (
	"math/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
)

// do not mutate it
var subtractPatterns = []int{10, 20, 30, 40, 50}

const (
	// 4ターン目移行ならカタストロフィ使用可能
	CatastropheEnableMinTurn = 4
)

type AbilityCatastrophe struct{}

func (a *AbilityCatastrophe) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	if IsUsed(targetAbility) ||
		!IsReady(targetAbility) ||
		targetAbility.Remaining == 0 {
		return false, nil
	}
	return true, nil
}

func (a *AbilityCatastrophe) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if !IsReady(&me.Abilities[abilityIndex]) {
		return false, nil, nil
	}

	if game.State.CurrentTurn < CatastropheEnableMinTurn {
		return false, nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_ABILITY_CATASTROPHE_NOT_MEET_TURN, "", nil)
	}

	players, e := db.GetPlayers(game.RoomID)
	if e != nil {
		return false, nil, e
	}
	for _, player := range players {
		if player.PlayerID == me.PlayerID {
			continue
		}

		// ランダムな数値コインを減算
		subtract := subtractPatterns[rand.Intn(len(subtractPatterns))]
		player.Coin = player.Coin - subtract
		if player.Coin < 0 {
			player.Coin = 0
		}

		db.SetPlayer(game.RoomID, &player)
	}

	// 発動者のアビリティが残り0なら使用済みのそうでないなら未使用にする
	me, e = db.GetPlayer(game.RoomID, me.PlayerID)
	if e != nil {
		return false, nil, e
	}
	if me.Abilities[abilityIndex].Remaining == 0 {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUsed)
	} else {
		me.Abilities[abilityIndex].Status = string(consts.AbilityStatusUnused)
	}
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	} else {
		return true, &me.Abilities[abilityIndex], nil
	}

}
