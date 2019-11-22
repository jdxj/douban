package modules

import (
	"douban/utils"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestReadChan(t *testing.T) {
	c := make(chan int)

	go func() {
		defer close(c)

		for i := 0; i < 100; i++ {
			c <- i
		}
	}()

	for {
		select {
		case data, ok := <-c:
			if ok {
				fmt.Println(data)
			} else {
				fmt.Println(data)
				time.Sleep(time.Second)
			}
		}
	}
}

func TestConnDB(t *testing.T) {
	err := utils.DB.Ping()
	if err != nil {
		panic(err)
	}

	defer utils.DB.Close()

	stmt, err := utils.DB.Prepare("insert into book_url (url) values (?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		_, err = stmt.Exec(strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func TestSQL(t *testing.T) {
	err := utils.DB.Ping()
	if err != nil {
		panic(err)
	}

	defer utils.DB.Close()

	rows, err := utils.DB.Query("select id from book_url where url=?", "fadfasd")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	id := -1
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}

		fmt.Println("id:", id)
	}

	if id == -1 {
		fmt.Println("not found")
	}
}

func TestRandUserAgent(t *testing.T) {
	ua := RandUserAgent()
	fmt.Println(ua)
}
