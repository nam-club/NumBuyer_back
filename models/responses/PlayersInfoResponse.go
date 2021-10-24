package responses

import "nam-club/NumBuyer_back/db"

type PlayersInfoResponse struct {
	RoomID  string               `json:"roomId"`
	Players []PlayersInfoPlayers `json:"players"`
}
type PlayersInfoPlayers struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	Coin       int    `json:"coin"`
	CardNum    int    `json:"cardNum"`
}

// レスポンスを生成
// DB接続、分岐などのビジネスロジックは書かないこと
func GeneratePlayersInfoResponse(players []db.Player, roomId string) *PlayersInfoResponse {
	ret := &PlayersInfoResponse{}
	ret.RoomID = roomId
	for _, v := range players {
		ret.Players = append(ret.Players,
			PlayersInfoPlayers{
				PlayerID:   v.PlayerID,
				PlayerName: v.PlayerName,
				Coin:       v.Coin,
				CardNum:    len(v.Cards),
			})
	}
	return ret

}
