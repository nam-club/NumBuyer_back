package routes

import (
	"encoding/json"
	"fmt"
	"nam-club/NumBuyer_back/models/orgerrors"

	socketio "github.com/googollee/go-socket.io"
	"gopkg.in/go-playground/validator.v9"
)

type RouteBase struct {
	server *socketio.Server
}

func NewRouteBase(server *socketio.Server) *RouteBase {
	return &RouteBase{server: server}
}

func (o *RouteBase) path(path string, f func(socketio.Conn, string)) {
	wrapFunc := func(s socketio.Conn, msg string) {
		fmt.Printf("[START REQUEST]%v\n", path)
		fmt.Printf("msg=%v\n", msg)
		f(s, msg)
		fmt.Printf("[END REQUEST]%v\n", path)
	}
	o.server.OnEvent("/", path, wrapFunc)
}

// リクエストメッセージの構造体への変換 & バリデーション
func Valid(reqBody string, result interface{}) error {
	if e := json.Unmarshal([]byte(reqBody), result); e != nil {
		return orgerrors.NewValidationError(e.Error())
	}

	v := validator.New()
	if e := v.Struct(result); e != nil {
		return orgerrors.NewValidationError(e.Error())
	}

	return nil
}
