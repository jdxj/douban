package book

import (
	"douban/modules"
	"douban/utils"
	"douban/utils/logs"
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
	rows, err := utils.DB.Query("select count(*) from book_url")
	if err != nil {
		logs.Logger.Error("sql: 'select count(*) from book_url', err: %s", err)
		return
	}

	var total int64
	for rows.Next() {
		if err = rows.Scan(&total); err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}
	rows.Close()

	if total <= 0 {
		logs.Logger.Warn("No book url available")
		return
	}

	stmtQuery, err := utils.DB.Prepare("select url from book_url order by id limit ?,?")
	if err != nil {
		logs.Logger.Error("sql: 'select url from book_url order by id limit ?,?', err: %s", err)
		return
	}
	defer stmtQuery.Close()

	stmtBookInsert, err := utils.DB.Prepare("insert into book (title, author, press) values (?,?,?)")
	if err != nil {
		logs.Logger.Error("sql: 'insert into book (title, author, press) values (?,?,?)', err: %s", err)
		return
	}
	defer stmtBookInsert.Close()

	stmtOpiInsert, err := utils.DB.Prepare("insert into opinion (score, amount, one, two, three, four, five, type, ref) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		logs.Logger.Error("sql: 'insert into opinion (score, amount, one, two, three, four, five, type, ref)', err: %s", err)
		return
	}
	defer stmtOpiInsert.Close()

	//for i := int64(0); i < total; i++ {
	//	utils.Pause(utils.Pause1s)
	//
	//	row, err := stmtQuery.Query(i, 1)
	//	if err != nil {
	//		logs.Logger.Error("%s", err)
	//		return
	//	}
	//
	//	var url string
	//	for row.Next() {
	//		if err = row.Scan(&url); err != nil {
	//			logs.Logger.Error("%s", err)
	//			row.Close()
	//			return
	//		}
	//	}
	//	row.Close()
	//
	//	fmt.Println(url)
	//}

	result, err := stmtBookInsert.Exec("iuy", "jdxj", "jdxj")
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	_, err = stmtOpiInsert.Exec(1.23, 4, 5.6, 7.8, 9.1, 11.12, 13.14, 1, id)
	if err != nil {
		panic(err)
	}
}

func TestSendEmail(t *testing.T) {
	w, err := NewWormhole()
	if err != nil {
		panic(err)
	}

	w.takeALook()
}
