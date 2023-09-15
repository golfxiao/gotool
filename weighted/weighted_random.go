package weighted

import (
	"math/rand"
	"time"
)

type WeightedRandom struct {
	random        *rand.Rand // 随机数生成器
	defaultWeight int        // 初始权重
}

func NewWeightedRandom(defaultWeight int) *WeightedRandom {
	if defaultWeight <= 0 {
		defaultWeight = 50
	}
	return &WeightedRandom{
		random:        rand.New(rand.NewSource(time.Now().UnixNano())),
		defaultWeight: defaultWeight,
	}
}

func (this *WeightedRandom) Draw(winCount []int, num int) []int {
	if len(winCount) == 0 {
		return []int{}
	}
	if num < 0 {
		num = 1
	}
	if num > len(winCount) {
		num = len(winCount)
	}
	userWeights := this.calculateWeights(winCount)
	return this.draw(userWeights, num)
}

func (this *WeightedRandom) DrawOne(winCount []int) int {
	userWeights := this.calculateWeights(winCount)
	winned := this.draw(userWeights, 1)
	if len(winned) > 0 {
		return winned[0]
	} else {
		return 0
	}
}

// 计算权重
func (this *WeightedRandom) calculateWeights(winCount []int) []int {
	weights := make([]int, len(winCount))
	for i, count := range winCount {
		weight := this.defaultWeight >> count
		if weight < 1 {
			weight = 1
		}
		weights[i] = weight
	}
	return weights
}

func (this *WeightedRandom) draw(weights []int, num int) []int {
	winned := make([]int, 0, num)

	// 总权重
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// 循环抽取多个用户
	for i := 0; i < num; i++ {
		n := this.random.Intn(totalWeight)
		for j, weight := range weights {
			if contains(winned, j) {
				continue
			}
			n -= weight
			if n < 0 {
				totalWeight -= weight
				winned = append(winned, j)
				break
			}
		}
	}

	return winned
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
