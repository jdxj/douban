package book

import (
	"douban/modules"
	"fmt"
	"testing"
)

func TestNewWormhole(t *testing.T) {
	wor, err := NewWormhole()
	if err != nil {
		panic(err)
	}

	wor.CaptureTags()
}

func TestGenBook(t *testing.T) {
	url := "https://book.douban.com/subject/25862578/"
	client := modules.GenHTTPClient()

	w := &Wormhole{}
	book, err := w.genBook(url, client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", book)
}
