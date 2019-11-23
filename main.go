package main

import (
	"douban/modules/book"
	"douban/utils"
	"douban/utils/logs"
	"time"
)

func main() {
	wor, err := book.NewWormhole()
	if err != nil {
		logs.Logger.Critical("%s", err)
		return
	}

	wor.Run()

	// 确保其他 goroutine 结束
	time.Sleep(time.Second)
	utils.Release()
}
