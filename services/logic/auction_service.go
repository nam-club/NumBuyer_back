package logic

import (
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"strconv"
)

// 入札する
func Bid(roomId, playerId string, bidAction consts.BidAction, coin int) (*responses.BidResponse, error) {
	if !CheckPhase(roomId, consts.PhaseAuction) {
		return nil, orgerrors.NewValidationError("not auction phase")
	}

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}
	if player.BuyAction.Value == consts.BidActionPass {
		return nil, orgerrors.NewValidationError("player already passed")
	}

	player.BuyAction.Action = bidAction.String()
	if bidAction == consts.BidActionBid {
		player.BuyAction.Value = strconv.Itoa(coin)
	} else if bidAction == consts.BidActionPass {
		player.Ready = true
	}
	player, e = db.SetPlayer(roomId, player)
	if e != nil {
		return nil, e
	}

	return &responses.BidResponse{PlayerName: player.PlayerName, Coin: coin}, nil
}

// プレイヤーのオークション終了時に必要な情報を取得する
func FetchAuctionEndInfo(roomId, playerId string) (*responses.BuyUpdateResponse, error) {

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	return &responses.BuyUpdateResponse{PlayerID: player.PlayerID, Coin: player.Coin, Cards: player.Cards}, nil
}

// 落札者を決定する
func DetermineBuyer(roomId string) (*db.Player, error) {

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}

	var buyer db.Player
	maxBidCoin := 0
	for _, p := range players {
		if p.BuyAction.Action == consts.BidActionBid.String() {
			b, e := strconv.Atoi(p.BuyAction.Value)
			if e == nil && b > maxBidCoin {
				maxBidCoin = b
				buyer = p
			}
		}
	}
	if maxBidCoin > 0 {
		return &buyer, nil
	} else {
		return nil, nil
	}
}

// オークションの状態をクリアする
func ClearAuction(roomId string) error {

	game, e := db.GetGame(roomId)
	if e != nil {
		return e
	}

	game.State.Auction = ""
	game, e = db.SetGame(roomId, game)
	if e != nil {
		return e
	}

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return e
	}

	for _, player := range players {
		player.BuyAction = db.BuyAction{}
		db.SetPlayer(roomId, &player)
	}

	return nil
}

// オークションカードをシャッフルする
func ShuffleAuctionCard(roomId string) (string, error) {

	game, e := db.GetGame(roomId)
	if e != nil {
		return "", e
	}

	// ランダムなオークションカードを生成する
	game.State.Auction = utils.GenerateRandomCard()

	game, e = db.SetGame(roomId, game)
	if e != nil {
		return "", e
	}

	return game.State.Auction, nil
}
