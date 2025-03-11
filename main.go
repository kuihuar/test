// package main

// import (
// 	"fmt"
// 	"test/ninepod"
// )

// func main() {
// 	fmt.Println("Hello World")

// 	alipay := &Alipay{}
// 	orderService := NewOrderService(alipay)
// 	orderService.PlaceOrder(100.55)
// 	m := map[string]string{"name": "zhangsan"}
// 	fmt.Println(m["name"])
// 	fmt.Println("---------------------------------")
// 	var prizes = []ninepod.Prize{
// 		{Name: "a", Probability: 0.01},
// 		{Name: "b", Probability: 0.10},
// 		{Name: "c", Probability: 0.20},
// 		{Name: "d", Probability: 0.20},
// 		{Name: "e", Probability: 0.15},
// 		{Name: "f", Probability: 0.30},
// 		{Name: "g", Probability: 0.05},
// 		{Name: "h", Probability: 0.05},
// 		{Name: "i", Probability: 0.001},
// 	}
// 	for i := 0; i < 10; i++ {
// 		result := ninepod.Lottery(prizes)
// 		fmt.Printf("第 %d 次抽奖结果: %v\n", i+1, result)
// 	}

// 	// prize := ninepod.Lottery()
// 	// fmt.Println(prize)
// }

// type Payment interface {
// 	Pay(amount float64)
// }

// type Alipay struct{}

// func (a *Alipay) Pay(amount float64) {
// 	fmt.Printf("支付宝支付了%.2f元\n", amount)
// }

// type Wechatpay struct{}

// func (w *Wechatpay) Pay(amount float64) {
// 	fmt.Printf("微信支付了%.2f元\n", amount)
// }

// type OrderService struct {
// 	payment Payment
// }

// func NewOrderService(payment Payment) *OrderService { return &OrderService{payment: payment} }

//	func (o *OrderService) PlaceOrder(amount float64) {
//		o.payment.Pay(amount)
//	}
package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Prize 奖品结构体
type Prize struct {
	ID          int     // 奖品ID
	Name        string  // 奖品名称
	Probability float64 // 概率（支持任意正数，自动计算比例）
}

// Lottery 根据概率抽奖
func Lottery(prizes []Prize) (Prize, error) {
	if len(prizes) == 0 {
		return Prize{}, fmt.Errorf("empty prize list")
	}

	// 计算概率总和并构建累加区间
	sum := 0.0
	accumulation := make([]float64, len(prizes))
	for i, p := range prizes {
		if p.Probability < 0 {
			return Prize{}, fmt.Errorf("negative probability not allowed")
		}
		sum += p.Probability
		accumulation[i] = sum
	}

	if sum <= 0 {
		return Prize{}, fmt.Errorf("total probability must be positive")
	}

	// 生成随机数（使用纳秒种子保证随机性）
	r := rand.Float64() * sum

	// 查找命中区间
	for i, acc := range accumulation {
		if r <= acc {
			return prizes[i], nil
		}
	}

	// 理论上不会执行到此处
	return prizes[len(prizes)-1], nil
}

func main() {
	// 示例奖品配置（概率值可自由设置，自动计算比例）
	prizes := []Prize{
		{1, "一等奖", 1},   // 1/(1+5+20+30) ≈ 1.79%
		{2, "二等奖", 5},   // 5/56 ≈ 8.93%
		{3, "三等奖", 10},  // 10/56 ≈ 17.86%
		{4, "谢谢参与", 40}, // 40/56 ≈ 71.43%
	}

	// 进行100万次抽奖测试概率分布
	result := make(map[int]int)
	totalTimes := 1000000

	for i := 0; i < totalTimes; i++ {
		prize, err := Lottery(prizes)
		if err != nil {
			panic(err)
		}
		result[prize.ID]++
	}

	// 打印统计结果
	fmt.Println("抽奖结果统计：")
	totalProb := 0.0
	for _, p := range prizes {
		prob := float64(result[p.ID]) / float64(totalTimes) * 100
		totalProb += prob
		fmt.Printf("%s(ID:%d) 理论概率：%.2f%%，实际概率：%.2f%%\n",
			p.Name, p.ID,
			p.Probability/56*100,
			prob)
	}
	fmt.Printf("总概率：%.2f%%\n", totalProb)
}

type Singleton struct {
}

var once sync.Once
var s *Singleton

func GetSingle() *Singleton {
	once.Do(func() {
		s = &Singleton{}
	})
	return s
}
