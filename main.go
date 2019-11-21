package main

import "douban/modules/book"

func main() {
	wor := book.NewWormhole()
	wor.CaptureTags()
}
