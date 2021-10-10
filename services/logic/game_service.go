package logic

import (
	"crypto/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
)

// 新規ゲームを生成する
func CreateNewGame(owner string) (*responses.Player, error) {

	var id string
	var e error
	if id, e = generateGameId(); e != nil {
		return nil, e
	}

	g := db.Game{
		RoomID: id,
		State: db.State{
			Phase:   consts.PhaseBeforeStart,
			Auction: "",
			Answer:  0,
		},
	}

	if _, e = db.SetGame(id, g); e != nil {
		return nil, e
	}
	var ret *responses.Player
	if ret, e = CreateNewPlayer(owner, id, true); e != nil {
		return nil, e
	}
	return ret, nil
}

// ランダムなゲームIDを一つ取得する
func GetRandomGameId() (string, error) {
	r, e := db.GetRandomGameId()
	if e != nil {
		return "", e
	}
	return r, nil
}

// ゲームIDを生成する
func generateGameId() (string, error) {
	const letters = "0123456789"

	for i := 0; i < 3; i++ {

		// 乱数を生成
		b := make([]byte, 10)
		if _, err := rand.Read(b); err != nil {
			return "", orgerrors.NewInternalServerError("")
		}

		var result string
		for _, v := range b {
			// index が letters の長さに収まるように調整
			result += string(letters[int(v)%len(letters)])
		}
		if b, _ := db.ExistsGame(result); b == false {
			return result, nil
		}
	}
	return "", orgerrors.NewInternalServerError("create room id error")
}
