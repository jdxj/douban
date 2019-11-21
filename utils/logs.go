package utils

import (
	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger

func init() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"douban.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	Logger = logs.GetBeeLogger()
}
