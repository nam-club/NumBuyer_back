package responses

var (
	ErrorValidation   = &Error{Code: "error.validation", Message: "validation error"}
	ErrorGameNotFound = &Error{Code: "error.game.notFound", Message: "not found game"}
	ErrorInternal     = &Error{Code: "error.internal", Message: "internal server error"}
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
