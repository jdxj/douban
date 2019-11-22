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

func DBConfig() *ini.Section {
	sec := iniFile.Section("db")
	return sec
}
