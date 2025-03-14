context包的基本作用:
Context主要用于在goroutine之间传递请求作用域的变量、取消信号和截止时间。这在处理HTTP请求或需要跨多个goroutine操作时特别有用，比如超时控制或取消操作
context的核心概念:
比如Context接口的结构，包括Deadline、Done、Err、Value这些方法。需要解释每个方法的作用，以及它们如何被用来控制goroutine的行为
context的创建方法:
用户需要知道如何使用context.Background()和context.TODO()作为根context，以及如何派生出子context，比如使用WithCancel、WithTimeout、WithDeadline和WithValue

实际应用场景:
比如超时控制、取消传播、请求跟踪和数据传递。每个场景需要举例子，比如在HTTP请求中设置超时，或者在多个goroutine之间传递跟踪ID。用户可能希望了解如何在实际项目中应用这些功能，避免资源泄漏或提高系统响应速度
context的实现原理:
需要了解底层结构，比如cancelCtx、timerCtx、valueCtx等。解释这些结构如何存储取消信号、超时信息或键值对，以及父子context之间的联动机制。比如，当父context被取消时，所有子context也会被取消，这是通过关闭Done通道来实现的。

可能被问到context的作用、WithValue的使用注意事项、如何实现超时控制，或者context对性能的影响。需要准备这些问题的答案，强调context的正确使用方法和潜在陷阱，比如WithValue只传递请求作用域的数据，避免滥用。

用户可能想知道context的使用注意事项，比如不要传递结构体、context的线程安全性，以及及时取消context释放资源。这些是实际开发中容易出错的地方，面试官可能会考察候选人的实际经验。

最后，总结context在Go并发编程中的重要性，强调它如何帮助管理goroutine的生命周期，避免泄漏，并提升系统的可靠性和可维护性。这部分需要简明扼要，突出关键点，帮助用户在面试中清晰表达。

一、context 的核心作用
传递请求作用域的数据
例如在 HTTP 请求链中传递请求 ID、用户身份等元数据。

控制协程生命周期
通过取消信号（Cancellation）、超时（Timeout）、截止时间（Deadline）等机制，优雅终止协程。

跨 API 边界传递信号
适用于 HTTP 服务、数据库操作、RPC 调用等多层调用链路


二、context.Context 接口
context.Context 定义了四个关键方法：
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)  // 返回设置的截止时间
    Done() <-chan struct{}                    // 返回一个 Channel，用于监听取消信号
    Err() error                               // 返回取消的原因（如超时、手动取消）
    Value(key interface{}) interface{}        // 获取关联的键值数据
}
```
三、context 的创建与派生
1. 根 Context
context.Background()
默认的根 Context，一般作为所有 Context 的起点（例如在 main 函数或测试中）。

context.TODO()
表示暂时不确定用途的占位 Context（通常在重构时使用）。

2. 派生 Context
通过根 Context 派生出子 Context，附加额外功能：

WithCancel(parent Context)
返回一个可手动取消的 Context 和取消函数 cancel()。


```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 手动触发取消
```
WithTimeout(parent Context, timeout time.Duration)
设置超时时间，自动取消（底层调用 WithDeadline）。

```go

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
```
WithDeadline(parent Context, d time.Time)
设置明确的截止时间，自动取消。

WithValue(parent Context, key, val interface{})
附加键值对数据，用于跨层传递请求作用域数据。


四、context 的实际应用场景
1. 超时控制
在 HTTP 请求或数据库操作中，避免长时间阻塞：
```go
func queryDatabase(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // 模拟数据库查询
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Query completed")
        return nil
    case <-ctx.Done():
        return ctx.Err() // 返回 "context deadline exceeded"
    }
}
```

2. 级联取消
父 Context 取消后，所有子 Context 自动取消：

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    go worker(ctx, "worker1")
    go worker(ctx, "worker2")
    
    time.Sleep(1 * time.Second)
    cancel() // 取消所有 worker
}

func worker(ctx context.Context, name string) {
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
```
3. 传递请求作用域数据
在中间件中注入数据（如请求 ID、用户 Token）：

```go
// 中间件设置数据
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), "requestID", uuid.New().String())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 业务逻辑获取数据
func handler(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value("requestID").(string)
    fmt.Fprintf(w, "Request ID: %s", requestID)
}
```
五、context 的实现原理
1. 底层结构
cancelCtx：实现取消功能的 Context，内部通过 done 通道和子 Context 链表管理取消信号。

timerCtx：继承 cancelCtx，增加定时器实现超时或截止时间。

valueCtx：存储键值对数据，通过链表结构支持层级查找
2. 取消传播机制
父 Context 被取消时，会遍历所有子 Context，关闭其 done 通道并设置错误原因。

通过 sync.Once 确保取消操作只执行一次。
六、常见面试问题
1. context 的作用是什么？
传递请求作用域数据、控制协程生命周期（超时/取消）。

2. context.Value 应该存储什么类型的数据？
仅限请求作用域的数据（如请求 ID、认证 Token），避免传递可选参数或函数依赖。

3. 如何实现跨多个 Goroutine 的级联取消？
使用 WithCancel 创建父 Context，子 Goroutine 监听 ctx.Done()。

4. context 是否线程安全？
是的，context 的取消操作和读取方法是线程安全的。

5. context 对性能的影响？
合理使用时性能开销极小，但滥用 WithValue 可能因反射查找数据降低性能。
七、使用注意事项
不要传递结构体
使用 WithValue 时，键应设为自定义类型（避免包间冲突），值尽量轻量。

及时调用 cancel()
防止 Context 和关联资源泄漏（结合 defer cancel() 使用）。

避免滥用
非请求作用域的数据传递应通过函数参数显式传递。
