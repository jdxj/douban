package utils

import (
	"douban/utils/logs"

	"gopkg.in/ini.v1"
)

var iniFile *ini.File

func init() {
	file, err := ini.Load("conf.ini")
	if err != nil {
		logs.Logger.Critical("%s", err)
		logs.Logger.Flush()

		panic(err)
	}

	iniFile = file
}

type DBConf struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	DBName   string `ini:"dbName"`
}

func GetDBConf() (*DBConf, error) {
	dbConf := new(DBConf)
	if err := iniFile.Section("db").MapTo(dbConf); err != nil {
		return nil, err
	}

	return dbConf, nil
}

type ModeConf struct {
	Mode int `ini:"mode"`
}

func GetModeConf() (*ModeConf, error) {
	modeConf := new(ModeConf)
	if err := iniFile.Section("mode").MapTo(modeConf); err != nil {
		return nil, err
	}

	return modeConf, nil
}
