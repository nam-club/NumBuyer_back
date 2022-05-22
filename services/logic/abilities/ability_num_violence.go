package abilities

import (
	"math/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
)

type AbilityNumViolence struct{}

func (a *AbilityNumViolence) CanActivate(game *db.Game, me *db.Player, targetAbility *db.Ability) (bool, error) {
	if IsUsed(targetAbility) ||
		!IsReady(targetAbility) ||
		game.State.Phase != consts.PhaseCalculate.Value ||
		!me.AnswerAction.Correct {
		return false, nil
	}

	players, e := db.GetPlayers(game.RoomID)
	if e != nil {
		return false, e
	}

	maxCards := 1
	for _, p := range players {
		if len(p.AnswerAction.Cards) > maxCards {
			maxCards = len(p.AnswerAction.Cards)
		}
	}
	return len(me.AnswerAction.Cards) == maxCards, nil
}

func (a *AbilityNumViolence) Fire(game *db.Game, me *db.Player, abilityIndex int) (bool, *db.Ability, error) {
	if !IsActive(&me.Abilities[abilityIndex]) {
		return false, nil, nil
	}
	me.Abilities[abilityIndex].Status = string(consts.AbilityStatusReady)
	if _, e := db.SetPlayer(game.RoomID, me); e != nil {
		return false, nil, e
	}

	players, e := db.GetPlayers(game.RoomID)
	if e != nil {
		return false, nil, e
	}

	for _, p := range players {
		if p.PlayerID == me.PlayerID {
			continue
		}

		if len(p.Cards) <= 2 {
			p.Cards = []string{}
			db.SetPlayer(game.RoomID, &p)
		} else {
			removeIndex1 := rand.Intn(len(p.Cards))
			removeIndex2 := rand.Intn(len(p.Cards))
			for {
				if removeIndex1 == removeIndex2 {
					removeIndex2 = rand.Intn(len(p.Cards))
				} else {
					break
				}
			}
			newCards := []string{}
			for i := 0; i < len(p.Cards); i++ {
				if i == removeIndex1 || i == removeIndex2 {
					continue
				}
				newCards = append(newCards, p.Cards[i])
			}
			p.Cards = newCards
			db.SetPlayer(game.RoomID, &p)
		}
	}
	return true, &me.Abilities[abilityIndex], nil

}
