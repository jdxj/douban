package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		"root",
		"",
		"localhost",
		"douban",
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	DB = db
}
