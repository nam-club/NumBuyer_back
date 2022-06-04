package orgerrors

import "github.com/pkg/errors"

type ValidationError struct {
	Status     int               `json:"status"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

const (
	VALIDATION_ERROR_BID_ACTION                    = "bid.action"
	VALIDATION_ERROR_ABILITY_PARSE                 = "ability.parseError"
	VALIDATION_ERROR_CALCULATE_ACTION_RESULT       = "calculate.actionResult"
	VALIDATION_ERROR_CALCULATE_ACTION              = "calculate.action"
	VALIDATION_ERROR_GAME_MODE                     = "game.mode"
	VALIDATION_ERROR_ABILITY_ALREADY_USED          = "ability.alreadyUsed"
	VALIDATION_ERROR_ABILITY_NOT_REMAINING         = "ability.notRemaining"
	VALIDATION_ERROR_BID_NOT_AUCTION_PHASE         = "bid.notAuctionPhase"
	VALIDATION_ERROR_BID_ALREADY_PASSED            = "bid.alreadyPassed"
	VALIDATION_ERROR_BID_INSUFFICIENT              = "bid.insuffcient"
	VALIDATION_ERROR_BID_EXCEED_MAX                = "bid.exceedMax"
	VALIDATION_ERROR_BID_PROHIBITED_IN_A_ROW       = "bid.ProhibitedInARow"
	VALIDATION_ERROR_CALCULATE_NOT_CALCULATE_PHASE = "calculate.notCalculatePhase"
	VALIDATION_ERROR_CALCULATE_ALREADY_READY       = "calculate.alreadyReady"
	VALIDATION_ERROR_CALCULATE_ALREADY_CORRECTED   = "calculate.alreadyCorrected"
	VALIDATION_ERROR_CALCULATE_ALREADY_PASSED      = "calculate.alreadyPassed"
	VALIDATION_ERROR_GAME_NOT_WAITING_PHASE        = "game.notWaitingPhase"
	VALIDATION_ERROR_GAME_NOT_MEET_PLAYER_NUMS     = "game.notMeetPlayerNums"
	VALIDATION_ERROR_ABLITY_RELOAD_INVAlID_PHASE   = "ability.reload.invalidPhase"
)

func (e *ValidationError) Error() string { return e.Message }

// バリデーションエラーはエラーコードを指定できるようにする。
// codeSuffix: error.validationに続くエラーコード。指定しなくてもOK
func NewValidationError(codeSuffix string, message string, params map[string]string) error {
	status := 400
	var code string
	if codeSuffix != "" {
		code = "error.validation." + codeSuffix
	} else {
		code = "error.validation"
	}

	if message == "" {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: "validation error"})
	} else {
		return errors.WithStack(&ValidationError{Status: status, Code: code, Message: message, Parameters: params})
	}
}
