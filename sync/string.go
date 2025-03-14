package sync

import "strings"

func concatStringsRegular(strs []string) string {
	result := ""
	for _, s := range strs {
		result += s // 每次拼接都生成新字符串
	}
	return result
}

func concatStringsOptimized(strs []string) string {
	var builder strings.Builder

	// 预分配足够的内存（可选，但推荐提前估算总长度）
	totalLen := 0
	for _, s := range strs {
		totalLen += len(s)
	}
	builder.Grow(totalLen) // 预分配内存，减少动态扩容次数

	// 拼接所有字符串
	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()

}
