package config

import (
	"math/rand"
	"time"
)

func InitConfig() {
	//  randのシード値
	rand.Seed(time.Now().Unix())
}
