package utils

import (
	"encoding/json"
	"log"
	"nam-club/NumBuyer_back/models/orgerrors"

	"github.com/pkg/errors"
)

// インスタンスをレスポンス形式（JSON文字列）に変換する
func Response(val interface{}) string {
	retJson, _ := json.Marshal(val)
	return string(retJson)
}

// インスタンスをレスポンス形式（JSON文字列）に変換する
func ResponseError(err error) string {
	errUnwrap := errors.Unwrap(err)
	var retJson []byte
	switch e := errUnwrap.(type) {
	case *orgerrors.ValidationError, *orgerrors.GameNotFoundError:
		retJson, _ = json.Marshal(e)
	case *orgerrors.InternalServerError:
		// stacktraceを吐きたいためzapでないloggerを使用
		// TODO zapでstacktrace吐けそうだったらそっちを使う
		log.Printf("[ERROR] %+v\n", e)
		retJson, _ = json.Marshal(e)
	default:
		// stacktraceを吐きたいためzapでないloggerを使用
		// TODO zapでstacktrace吐けそうだったらそっちを使う
		log.Printf("[ERROR] %+v\n", err)
		retJson, _ = json.Marshal(errors.Unwrap(orgerrors.NewInternalServerError("", nil)))
	}

	return string(retJson)
}
