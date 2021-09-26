package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"

	"github.com/google/uuid"
)

// 新規ゲームを生成する
func CreateNewGame(owner string) *db.Game {

	id := uuid.Must(uuid.NewUUID()).String()
	g := db.Game{
		RoomID: id,
		State: db.State{
			Phase:   consts.PhaseBeforeStart,
			Auction: "",
			Answer:  0,
		},
	}
	db.Transaction(func() {
		db.SetGame(id, g)
		CreateNewPlayer(owner, id, true)
	})

	return &g
}
