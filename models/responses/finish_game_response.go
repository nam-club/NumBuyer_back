package responses

type FinishGameResponse struct {
	Players []FinishGamePlayers `json:"players"`
}
type FinishGamePlayers struct {
	PlayerName string `json:"playerName"`
	Rank       int    `json:"rank"`
	Coin       int    `json:"coin"`
}
