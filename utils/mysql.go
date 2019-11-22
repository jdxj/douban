package utils

import (
	"database/sql"
	"douban/utils/logs"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	dbConf, err := GetDBConf()
	if err != nil {
		logs.Logger.Critical("%s", err)
		logs.Logger.Flush()

		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		dbConf.Username,
		dbConf.Password,
		dbConf.Host,
		dbConf.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logs.Logger.Critical("%s", err)
		logs.Logger.Flush()

		panic(err)
	}
	DB = db
}
