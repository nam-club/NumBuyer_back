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
		return orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_GAME_MODE, "invalid game mode", nil)
	}
}

func ParseGameMode(s int) (v GameMode, err error) {
	v = GameMode(s)
	err = v.Valid()
	return
}
