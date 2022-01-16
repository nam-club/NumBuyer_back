package logic

import (
	"fmt"
	"math/rand"
	"nam-club/NumBuyer_back/consts"
	"nam-club/NumBuyer_back/db"
	"nam-club/NumBuyer_back/models/orgerrors"
	"nam-club/NumBuyer_back/models/responses"
	"nam-club/NumBuyer_back/utils"
	"strconv"
	"time"
)

// 正解者を全て取得する
func PickAllCorrector(roomId string) ([]db.Player, error) {

	players, e := db.GetPlayers(roomId)
	if e != nil {
		return nil, e
	}

	var correctors []db.Player
	for _, p := range players {
		if p.AnswerAction.Correct {
			correctors = append(correctors, p)
		}
	}

	return correctors, nil
}

// 計算結果をクリアする
func ClearCalculateAction(roomId string) error {
	players, e := db.GetPlayers(roomId)
	if e != nil {
		return e
	}

	for _, player := range players {
		player.Ready = false
		player.AnswerAction = db.AnswerAction{}
		db.SetPlayer(roomId, &player)
	}

	return nil
}

// 解答をシャッフルする
func ShuffleAnswer(roomId string) (string, error) {

	// ランダムな解答を生成する
	game, e := db.GetGame(roomId)
	if e != nil {
		return "", e
	}
	newAnswer := strconv.Itoa(rand.Intn(consts.TargetMax-consts.TargetMin) + consts.TargetMin)

	game.State.Answer = newAnswer
	_, e = db.SetGame(roomId, game)
	if e != nil {
		return "", e
	}

	return newAnswer, nil
}

func CalculateSubmits(roomId, playerId string, action consts.CalculateAction, submits []string) (*responses.CalculateResponse, error) {
	if !CheckPhase(roomId, consts.PhaseCalculate) {
		return nil, orgerrors.NewValidationError("not calculate phase")
	}

	game, e := db.GetGame(roomId)
	if e != nil {
		return nil, e
	}

	player, e := db.GetPlayer(roomId, playerId)
	if e != nil {
		return nil, e
	}

	if player.AnswerAction.Correct {
		return nil, orgerrors.NewValidationError("player already correctly answered")
	}

	if action == consts.CalculateActionPass {
		// actionがpassの場合ステータスをpassに更新してリターン
		player.AnswerAction.Action = action.String()
		player.Ready = true
		player, e = db.SetPlayer(roomId, player)
		if e != nil {
			return nil, e
		}
		return &responses.CalculateResponse{
			IsCorrectAnswer: false,
			PlayerID:        playerId,
			Coin:            player.Coin,
			Cards:           player.Cards,
		}, nil
	} else {
		// actionがanswerの場合

		// カードが正しいかバリデーション
		validateCards := make([]string, len(player.Cards))
		copy(validateCards, player.Cards)
		for _, s := range submits {
			if i := utils.ContainsStringWithIndex(validateCards, s); i != -1 {
				validateCards = utils.DeleteSliceElement(validateCards, i)
			} else {
				return nil, orgerrors.NewValidationError("player is not have submitted cards")
			}
		}

		// 計算
		calculated, e := calculate(submits)
		if e != nil {
			return nil, e
		}

		player.AnswerAction.Action = action.String()
		player.AnswerAction.AnswerTime = time.Now().Format(time.RFC3339)

		// 結果の判定
		if game.State.Answer == strconv.Itoa(calculated) {
			// 正解した時、正解者のコイン数とカード情報を更新する
			player.Coin += calculated + len(submits)
			player.AnswerAction.Correct = true
			player.Ready = true
			for _, s := range submits {
				i := utils.ContainsStringWithIndex(player.Cards, s)
				player.Cards = utils.DeleteSliceElement(player.Cards, i)
			}
			player, e = db.SetPlayer(roomId, player)
			if e != nil {
				return nil, e
			}
			return &responses.CalculateResponse{
				IsCorrectAnswer: true,
				PlayerID:        playerId,
				Coin:            player.Coin,
				Cards:           player.Cards,
				AddedCoin: responses.AddedCoin{
					Total:        calculated,
					CardNumBonus: len(submits),
				},
			}, nil
		} else {
			// 不正解の時
			player.AnswerAction.Correct = false
			player, e = db.SetPlayer(roomId, player)
			if e != nil {
				return nil, e
			}
			return &responses.CalculateResponse{
				IsCorrectAnswer: false,
				PlayerID:        playerId,
				Coin:            player.Coin,
				Cards:           player.Cards,
			}, nil
		}
	}
}

// カードの入力から計算する
func calculate(submits []string) (int, error) {
	calculated := 0
	code := ""
	for i, submit := range submits {
		if i%2 == 0 {
			submitInt, e := strconv.Atoi(submit)
			if (e != nil) ||
				(submitInt < consts.TermMin || consts.TermMax <= submitInt) {
				return -1, orgerrors.NewValidationError("Invalid calculate card: " + submit)
			}
			switch code {
			case consts.CodePlus:
				calculated = calculated + submitInt
			case consts.CodeMinus:
				calculated = calculated - submitInt
			case consts.CodeTimes:
				calculated = calculated * submitInt
			case consts.CodeDivide:
				calculated = calculated / submitInt
			default:
				calculated = submitInt
			}
			code = ""
		} else {
			if utils.ContainsString(consts.Codes, submit) {
				code = submit
			} else {
				return -1, orgerrors.NewValidationError("Invalid calculate card: " + submit)
			}
		}
	}
	if code != "" {
		return -1, orgerrors.NewValidationError(fmt.Sprintf("Invalid calculate submits: %v", submits))
	}
	return calculated, nil
}
