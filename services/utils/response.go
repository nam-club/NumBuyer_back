package utils

import (
	"encoding/json"
)

// インスタンスをレスポンス形式（JSON文字列）に変換する
func ToResponseFormat(val interface{}) string {
	// encode json
	ret_json, _ := json.Marshal(val)
	return string(ret_json)
}
