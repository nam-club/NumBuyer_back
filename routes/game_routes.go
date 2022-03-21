package routes

import (
	"errors"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/requests"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/services/logic"
	"nam-club/NumBuyer_back/utils"

	socketio "github.com/googollee/go-socket.io"
)

func RoutesGame(r *RouteBase) {

	r.path(consts.TSJoinQuickMatch, func(s socketio.Conn, msg string) {
		req := &requests.JoinQuickMatch{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		roomId, e := logic.GetRandomRoomId()
		if e != nil {
			// 部屋が見つからなかった場合は新規作成
			switch errors.Unwrap(e).(type) {
			case *orgerrors.GameNotFoundError:
				abilities := consts.ParseAbilities(req.AbilityIds)
				resp, e := logic.CreateNewGame(req.PlayerName,
					consts.QuickMatchPlayersMin,
					consts.QuickMatchPlayersMax,
					consts.GameModeQuickMatch,
					abilities)
				if e != nil {
					s.Emit(consts.FSGameJoin, utils.ResponseError(e))
					return
				}

				// フェーズのタイマーをスタート
				if e := logic.CanCreateGameScheduler(resp.RoomID); e != nil {
					s.Emit(consts.FSGameJoin, utils.ResponseError(e))
					return
				}
				logic.NewPhaseScheduler(resp.RoomID, r.server).Start()

				s.LeaveAll()
				s.Join(resp.RoomID)

				s.Emit(consts.FSGameJoin, utils.Response(resp))
				return
			default:
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}
		} else {
			// ゲームに参加可能かチェック
			joinable := logic.CheckPhase(roomId, consts.PhaseWaiting)
			if !joinable {
				s.Emit(consts.FSGameJoin, orgerrors.NewValidationError("can not join game"))
				return
			}

			// 部屋が見つかった場合はその部屋に参加
			abilities := consts.ParseAbilities(req.AbilityIds)
			player, e := logic.CreateNewPlayer(req.PlayerName, roomId, false, abilities)
			if e != nil {
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}

			// 一つの部屋にのみ入室した状態にする
			s.LeaveAll()
			s.Join(roomId)
			resp := responses.JoinResponse{RoomID: roomId, PlayerID: player.PlayerID}
			s.Emit(consts.FSGameJoin, utils.Response(resp))

			// 人数が揃っていたらゲームを開始する
			players, e := logic.GetPlayersInfo(roomId)
			if e != nil {
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}
			game, e := logic.GetGame(roomId)
			if e != nil {
				s.Emit(consts.FSGameJoin, utils.ResponseError(e))
				return
			}
			if game.PlayersMin <= len(players.Players) {
				if e := logic.StartGame(roomId); e != nil {
					s.Emit(consts.FSGameStart, utils.ResponseError(e))
					return
				}
				respStart := responses.GenerateGameStartResponse(roomId, consts.CoinClearNum)
				r.server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameStart, utils.Response(respStart))
			}

		}
	})

	r.path(consts.TSJoinFriendMatch, func(s socketio.Conn, msg string) {
		req := &requests.JoinFriendMatch{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		// ゲームに参加可能かチェック
		joinable := logic.CheckPhase(req.RoomID, consts.PhaseWaiting)
		if !joinable {
			s.Emit(consts.FSGameJoin, orgerrors.NewValidationError("can not join game"))
			return
		}

		abilities := consts.ParseAbilities(req.AbilityIds)
		player, e := logic.CreateNewPlayer(req.PlayerName, req.RoomID, false, abilities)
		if e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		// 一つの部屋にのみ入室した状態にする
		s.LeaveAll()
		s.Join(req.RoomID)

		resp := responses.JoinResponse{RoomID: req.RoomID, PlayerID: player.PlayerID}
		s.Emit(consts.FSGameJoin, utils.Response(resp))
	})

	r.path(consts.TSCreateMatch, func(s socketio.Conn, msg string) {
		req := &requests.CreateMatch{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		abilities := consts.ParseAbilities(req.AbilityIds)
		resp, e := logic.CreateNewGame(req.PlayerName,
			consts.FriendMatchPlayersMin,
			consts.FriendMatchPlayersMax,
			consts.GameModeFriendMatch,
			abilities)
		if e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}

		// フェーズのタイマーをスタート
		if e := logic.CanCreateGameScheduler(resp.RoomID); e != nil {
			s.Emit(consts.FSGameJoin, utils.ResponseError(e))
			return
		}
		logic.NewPhaseScheduler(resp.RoomID, r.server).Start()

		s.LeaveAll()
		s.Join(resp.RoomID)
		s.Emit(consts.FSGameJoin, utils.Response(resp))
	})

	r.path(consts.TSGetAbilities, func(s socketio.Conn, msg string) {
		resp := responses.GenerateGetAbilitiesResponse(consts.GetAbilities())
		s.Emit(consts.FSGetAbilities, utils.Response(resp))
	})

	r.path(consts.TSGameReadyAbilities, func(s socketio.Conn, msg string) {
		req := &requests.GameReadyAbilities{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameReadyAbilities, utils.ResponseError(e))
			return
		}

		status, e := logic.ReadyAbilities(req.RoomID, req.PlayerID, req.AbilityIds)
		if e != nil {
			s.Emit(consts.FSGameReadyAbilities, utils.ResponseError(e))
			return
		}

		s.Emit(consts.FSGameReadyAbilities, utils.Response(responses.GameReadyAbilitiesResponse{Status: status}))
	})

	r.path(consts.TSGamePlayersInfo, func(s socketio.Conn, msg string) {
		req := &requests.GamePlayerInfo{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGamePlayersInfo, utils.ResponseError(e))
			return
		}
		resp, e := logic.GetPlayersInfo(req.RoomID)
		if e != nil {
			s.Emit(consts.FSGamePlayersInfo, utils.ResponseError(e))
			return
		}
		r.server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGamePlayersInfo, utils.Response(resp))
	})

	r.path(consts.TSGameStart, func(s socketio.Conn, msg string) {
		req := &requests.GameStart{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}

		if e := logic.StartGame(req.RoomID); e != nil {
			s.Emit(consts.FSGameStart, utils.ResponseError(e))
			return
		}

		resp := responses.GenerateGameStartResponse(req.RoomID, consts.CoinClearNum)
		r.server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameStart, utils.Response(resp))
	})

	r.path(consts.TSGameNextTurn, func(s socketio.Conn, msg string) {
		req := &requests.GameNextTurn{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameNextTurn, utils.ResponseError(e))
			return
		}
		resp, e := logic.FetchNextTurnInfo(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FSGameNextTurn, utils.ResponseError(e))
			return
		}
		s.Emit(consts.FSGameNextTurn, utils.Response(resp))
	})

	r.path(consts.TSGameBid, func(s socketio.Conn, msg string) {
		req := &requests.GameBid{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameBid, utils.ResponseError(e))
			return
		}
		bidAction, e := consts.ParseBidAction(req.Action)
		if e != nil {
			s.Emit(consts.FSGameBid, utils.ResponseError(e))
			return
		}

		resp, e := logic.Bid(req.RoomID, req.PlayerID, bidAction, req.Coin)
		if e != nil {
			s.Emit(consts.FSGameBid, utils.ResponseError(e))
			return
		}

		// Bid時のみレスポンスを返却
		if bidAction == consts.BidActionBid {
			r.server.BroadcastToRoom("/", s.Rooms()[0], consts.FSGameBid, utils.Response(resp))
			return
		}
	})

	r.path(consts.TSGameBuy, func(s socketio.Conn, msg string) {
		req := &requests.GameBuy{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameBuyUpdate, utils.ResponseError(e))
			return
		}
		resp, e := logic.FetchAuctionEndInfo(req.RoomID, req.PlayerID)
		if e != nil {
			s.Emit(consts.FSGameBuyUpdate, utils.ResponseError(e))
			return
		}
		s.Emit(consts.FSGameBuyUpdate, utils.Response(resp))
	})

	r.path(consts.TSGameCalculate, func(s socketio.Conn, msg string) {
		req := &requests.GameCalculate{}
		if e := Valid(msg, req); e != nil {
			s.Emit(consts.FSGameCalculateResult, utils.ResponseError(e))
			return
		}

		action, e := consts.ParseCalculateAction(req.Action)
		if e != nil {
			s.Emit(consts.FSGameCalculateResult, utils.ResponseError(e))
			return
		}

		resp, e := logic.CalculateSubmits(req.RoomID, req.PlayerID, action, req.CalculateCards)
		if e != nil {
			s.Emit(consts.FSGameCalculateResult, utils.ResponseError(e))
			return
		}

		// Answer時のみレスポンスを返却
		if action == consts.CalculateActionAnswer {
			s.Emit(consts.FSGameCalculateResult, utils.Response(resp))
			return
		}

	})
}
