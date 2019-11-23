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
	url := "https://book.douban.com/subject/1948901/"
	client := modules.GenHTTPClient()

	w := &Wormhole{}
	book, err := w.genBook(url, client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", book)
	fmt.Printf("%+v\n", book.Info)
	fmt.Printf("%+v\n", book.Opinion)
}

func TestClearColon(t *testing.T) {
	rows := []string{
		"author:",
		"jdxj",
		"title:",
		"abcff:",
		"defaa",
		"press:",
		"apple",
	}
	info := new(Info)
	rows = info.clearColon(rows)
	fmt.Println(rows)
}

func TestSQL(t *testing.T) {
}
