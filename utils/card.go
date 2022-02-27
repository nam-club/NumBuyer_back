package utils

import (
	"math/rand"
	"nam-club/NumBuyer_back/consts"
	"strconv"
)

// ランダムなカードを生成する
func GenerateRandomCard(num int) []string {

	generated := []string{}
	for i := 0; i < num; i++ {

		// ランダムなオークションカードを生成する
		if rand.Intn(100)+1 <= consts.AuctionCodeProbability {
			// 符号を生成
			index := rand.Intn(len(consts.Codes))
			generated = append(generated, consts.Codes[index])
		} else {
			// 数字を生成（最小値以上、最大値未満）
			generated = append(generated, strconv.Itoa(rand.Intn(consts.TermMax-consts.TermMin)+consts.TermMin))
		}
	}
	return generated
}
