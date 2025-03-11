package ninepod

import (
	"math/rand"
	"time"
)

type Prize struct {
	Name        string
	Probability float64
}

var prizes = []Prize{
	{"a", 0.01},
	{"b", 0.10},
	{"c", 0.20},
	{"d", 0.20},
	{"e", 0.15},
	{"f", 0.30},
	{"g", 0.05},
	{"h", 0.05},
	{"i", 0.001},
}

func Lottery(prizes []Prize) string {
	rand.Seed(time.Now().UnixNano())
	totalProbability := 0.0
	// 计算总概率
	for _, prize := range prizes {
		totalProbability += prize.Probability
	}

	// 生成 [0, totalProbability) 范围内的随机数
	randomNum := rand.Float64() * totalProbability
	currentProbability := 0.0
	// 遍历奖品列表，确定中奖奖品
	for _, prize := range prizes {
		currentProbability += prize.Probability
		if randomNum < currentProbability {
			return prize.Name
		}
	}
	return ""
}
