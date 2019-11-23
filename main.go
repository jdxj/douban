package main

import (
	"douban/modules/book"
	"douban/utils/logs"
)

func main() {
	wor, err := book.NewWormhole()
	if err != nil {
		logs.Logger.Critical("%s", err)
		return
	}

	wor.Run()
}
