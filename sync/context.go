package sync

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ExampleWithValue() {
	ctx := context.Background()
	ctx = context.TODO()

	ctx = context.WithValue(ctx, "userID", 123)

	userID, ok := ctx.Value("userID").(int)
	if !ok {
		fmt.Println("not found")
	}
	fmt.Println("userID", userID)
}

func ExampleWithTimtout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("timeout", ctx.Err())
	case <-time.After(3 * time.Second):
		fmt.Println("done")
	}
}

func queyDb(ctx context.Context, query string) (string, error) {
	resultCh := make(chan string, 1)
	// errCh := make(chan error)

	go func() {
		time.Sleep(time.Duration(time.Now().UnixNano()%3+1) * time.Second)
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(rng.Intn(3)+1) * time.Second)

		if ctx.Err() != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:

			//发送后立即退出
			resultCh <- fmt.Sprintf("%s result", query)
			return
		}

	}()
	select {
	case result := <-resultCh:
		return result, nil
	// case err := <-errCh:
	// 	return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func aggregateData(ctx context.Context, date string) (string, error) {
	resultCh := make(chan string, 1)
	go func() {
		duration := 10 * time.Now().UnixNano() % 11
		time.Sleep(time.Duration(duration) * time.Second)
		if ctx.Err() != nil {
			return
		}
		resultCh <- fmt.Sprintf("%s result, duration:%d", date, duration)

	}()
	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func ExamplWithDeadLine() {
	deadline := time.Date(2025, 3, 15, 14, 20, 0, 0, time.UTC)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	taskStart := time.Now()
	fmt.Printf("task start at %v\n", taskStart.Format("2006-01-02 15:04:05"))
	result, err := aggregateData(ctx, "2025-03-15")
	if err != nil {
		return
	}
	fmt.Println("aggregateData res", result)

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 启动两个 Worker
	wg.Add(2)
	go worker(ctx, "worker1", &wg)
	go worker(ctx, "worker2", &wg)

	time.Sleep(1 * time.Second)
	cancel() // 发送取消信号

	wg.Wait() // 等待所有 Worker 退出
}

func worker(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done() // 确保退出时通知 WaitGroup
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name, "stopped:", ctx.Err())
			return
		default:
			fmt.Println(name, "working...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
