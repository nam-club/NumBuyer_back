package consts

import (
	"github.com/kelseyhightower/envconfig"
)

// 環境変数
// https://github.com/kelseyhightower/envconfig を使う
// 参考： https://qiita.com/andromeda/items/c5195307cd08537d4fad
func init() {
	var e EnvConst
	envconfig.Process("", &e)
	Env = e
}

var (
	Env EnvConst
)

type EnvConst struct {
	RedisUrl string `default:"127.0.0.1:6379"`
}
