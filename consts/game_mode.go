package consts

import (
	"nam-club/NumBuyer_back/models/orgerrors"
)

type GameMode int

const (
	GameModeQuickMatch GameMode = iota
	GameModeFriendMatch
)

func (v GameMode) Valid() error {
	switch v {
	case GameModeQuickMatch, GameModeFriendMatch:
		return nil
	default:
		return orgerrors.NewValidationError("invalid bid action type")
	}
}

func ParseGameMode(s int) (v GameMode, err error) {
	v = GameMode(s)
	err = v.Valid()
	return
}
