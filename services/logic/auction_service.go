package logic

import (
	"math/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"strconv"
)

// 入札する
func Bid(roomId, playerId string, bidAction consts.BidAction, coin int) (*responses.BidResponse, error) {
	// ゲーム操作をロック
	if locked, e := db.SetLock(db.CreateLockKey(roomId, playerId), consts.MutexTTL); locked || e != nil {
		if e != nil {
			return nil, e
		}
		return nil, orgerrors.NewMutexError("calculate mutex error", nil)
	}
	// ロック情報を最後に削除
	defer db.DeleteLock(db.CreateLockKey(roomId, playerId))

	if !CheckPhase(roomId, consts.PhaseAuction) {
		return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_BID_NOT_AUCTION_PHASE, "not auction phase", nil)
	}

	game, e := db.GetGame(roomId)
	if e != nil {
		return nil, e
	}

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	if player.BuyAction.Value == consts.BidActionPass.String() {
		return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_BID_ALREADY_PASSED, "player already passed", nil)
	}

	player.BuyAction.Action = bidAction.String()
	if bidAction == consts.BidActionBid {
		// バリデーション
		if maxBid, _ := strconv.Atoi(game.State.AuctionMaxBid); maxBid >= coin {
			return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_BID_INSUFFICIENT, "insufficient bid", nil)
		}
		if player.BuyAction.BidCount >= consts.AuctionMaxBidCount {
			return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_BID_EXCEED_MAX, "exceed max bid count", nil)
		}
		if game.State.AuctionLastBidPlayerId == playerId {
			return nil, orgerrors.NewValidationError(orgerrors.VALIDATION_ERROR_BID_PROHIBITED_IN_A_ROW, "cannot bid in a row", nil)
		}

		// Bid情報セット処理
		player.BuyAction.Value = strconv.Itoa(coin)
		player.BuyAction.BidCount = player.BuyAction.BidCount + 1
		game.State.AuctionMaxBid = player.BuyAction.Value
		game.State.AuctionLastBidPlayerId = playerId
		if _, e := db.SetGame(roomId, game); e != nil {
			return nil, e
		}

		// オークションフェーズ期限更新処理
		remainSeconds := 0 // 0なら更新なし（レスポンスに含めない）
		if isResetted, e := ResetTimer(roomId, consts.AuctionResetTimeRemains); e != nil {
			return nil, e
		} else if isResetted {
			remainSeconds = consts.AuctionResetTimeRemains
		}

		if _, e = db.SetPlayer(roomId, player); e != nil {
			return nil, e
		}
		return &responses.BidResponse{
				PlayerName:    player.PlayerName,
				Coin:          coin,
				RemainingTime: remainSeconds},
			nil
	} else if bidAction == consts.BidActionPass {
		player.Ready = true
		if _, e = db.SetPlayer(roomId, player); e != nil {
			return nil, e
		}
	}
	return nil, nil
}

// プレイヤーのオークション終了時に必要な情報を取得する
func FetchAuctionEndInfo(roomId, playerId string) (*responses.BuyUpdateResponse, error) {
	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	return &responses.BuyUpdateResponse{
			PlayerID:    player.PlayerID,
			IsSuccessed: player.BuyAction.IsBuyer,
			Coin:        player.Coin,
			Cards:       player.Cards},
		nil
}

// 落札者を決定する
// 最終入札者 = 落札者
func DetermineBuyer(roomId string) (*db.Player, error) {
	game, e := db.GetGame(roomId)
	if e != nil {
		return nil, e
	}
	if game.State.AuctionLastBidPlayerId != "" {
		buyer, e := db.GetPlayer(roomId, game.State.AuctionLastBidPlayerId)
		if e != nil {
			return nil, e
		}
		buyer.BuyAction.IsBuyer = true
		buyer, e = db.SetPlayer(roomId, buyer)
		if e != nil {
			return nil, e
		}

		return buyer, nil
	} else {
		return nil, nil
	}
}

// オークションの状態をクリアし再セットする
func ClearAndResetAuction(roomId string) error {

	// オークションの状態をクリア
	game, e := db.GetGame(roomId)
	if e != nil {
		return e
	}

	game.State.Auction = []string{}
	game.State.AuctionMaxBid = ""
	game.State.AuctionLastBidPlayerId = ""
	if _, e = db.SetGame(roomId, game); e != nil {
		return e
	}

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return e
	}

	for _, player := range players {
		player.Ready = false
		player.BuyAction = db.BuyAction{}
		db.SetPlayer(roomId, &player)
	}

	// オークションカードをシャッフル
	_, e = ShuffleAuctionCard(roomId)
	if e != nil {
		return e
	}

	return nil
}

// オークションカードをシャッフルする
func ShuffleAuctionCard(roomId string) ([]string, error) {

	game, e := db.GetGame(roomId)
	if e != nil {
		return []string{}, e
	}

	// ランダムなオークションカードを生成する
	game.State.Auction = utils.GenerateRandomCard(
		rand.Intn(consts.AuctionCardsNumMax-consts.AuctionCardsNumMin) + consts.AuctionCardsNumMin)

	game, e = db.SetGame(roomId, game)
	if e != nil {
		return []string{}, e
	}

	return game.State.Auction, nil
}
