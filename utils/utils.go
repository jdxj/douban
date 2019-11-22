package utils

import (
	"strings"
	"time"
)

const (
	Pause2s = 2 * time.Second
	Pause5s = 5 * time.Second
)

func Pause(dur time.Duration) {
	time.Sleep(dur)
}

// CleanAndJoin 用于删除 s 中的空白,
// 并使用 sep 将 s 的每个部分连接成一个字符串.
func CleanAndJoin(s, sep string) string {
	ss := strings.Fields(s)
	return strings.Join(ss, sep)
}

// CleanAndSplit 将去除 s 中的空白,
// 并将 s 根据空白分割.
func CleanAndSplit(s string) []string {
	return strings.Fields(s)
}
