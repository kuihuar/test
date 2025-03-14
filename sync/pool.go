package sync

import (
	"encoding/json"
	"fmt"
	"sync"
)

// 1. 对象复用（手动管理）
// 直接复用已有的对象，避免重复分配内存。例如，在高频解析 JSON 的代码中复用结构体：
// 手动管理，多GOROUTINE下可能 经发数据竞争，简单，适合单GOROUTINE场景
type Request struct {
	Data map[string]interface{}
}

var requestPool = &Request{Data: make(map[string]interface{})}

func processRequest(data []byte) {
	req := requestPool
	clear(req.Data)
	if err := json.Unmarshal(data, &req.Data); err != nil {
		return
	}
	fmt.Println(req.Data)
}

// 3. 使用 sync.Pool 的注意事项
// 对象生命周期：sync.Pool 中的对象可能被 GC 回收，不能依赖其长期存在。

// 重置状态：放回 Pool 前必须重置对象状态，避免脏数据。

// 适用场景：适合复用高频创建/销毁的临时对象（如解析器、缓冲区）。

// 避免滥用：长期存活的对象不适合用 sync.Pool，可能占用内存。

var requestPool2 = sync.Pool{
	New: func() interface{} {
		return &Request{
			Data: make(map[string]interface{}),
		}
	},
}

func ProcessRequest2(data []byte) {
	req := requestPool2.Get().(*Request)
	defer requestPool2.Put(req)

	clear(req.Data)
	if err := json.Unmarshal(data, &req.Data); err != nil {
		return
	}
	fmt.Println(req.Data)

}
