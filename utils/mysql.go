package utils

import (
	"database/sql"
	"douban/utils/logs"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	username := DBConfig().Key("username").String()
	password := DBConfig().Key("password").String()
	host := DBConfig().Key("host").String()
	dbName := DBConfig().Key("dbName").String()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		username,
		password,
		host,
		dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logs.Logger.Critical("%s", err)
		logs.Logger.Flush()

		panic(err)
	}
	DB = db
}
