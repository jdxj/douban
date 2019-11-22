package book

import (
	"douban/utils"
	"douban/utils/logs"
	"strings"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Press  string

	Info   *Info
	Starts *Stars
}

func (b *Book) ToScan() []interface{} {
	var row []interface{}
	row = append(row, &b.ID)
	row = append(row, &b.Title)
	row = append(row, &b.Author)
	row = append(row, &b.Press)
	return row
}

func (b *Book) ToInsert() []interface{} {
	var row []interface{}
	row = append(row, &b.Title)
	row = append(row, &b.Author)
	row = append(row, &b.Press)
	return row
}

// 已知属性
var Attr = map[string]struct{}{
	"作者":   struct{}{},
	"出版社":  struct{}{},
	"出品方":  struct{}{},
	"副标题":  struct{}{},
	"原作名":  struct{}{},
	"译者":   struct{}{},
	"出版年":  struct{}{},
	"页数":   struct{}{},
	"定价":   struct{}{},
	"装帧":   struct{}{},
	"丛书":   struct{}{},
	"ISBN": struct{}{},
	"统一书号": struct{}{},
}

// Info 是比 Book 具有更多书籍信息的结构,
// 但是 Info 不存储 title.
type Info struct {
	Author     string
	Press      string
	Producer   string
	Subtitle   string
	OriTitle   string
	Translator string
	Year       string // 出版年
	Pages      string
	Price      string
	Binding    string
	Series     string
	ISBN       string
	Number     string // 统一书号
}

// Unmarshal 用于从 content 解析数据到 info 中.
func (info *Info) Unmarshal(rawInfo string) {
	rows := utils.CleanAndSplit(rawInfo)
	logs.Logger.Debug("length:", len(rows))
	logs.Logger.Debug("rows:", rows)
	for _, row := range rows {
		cols := strings.Split(row, ":")
		if len(cols) < 2 {
			logs.Logger.Warn("Can not unmarshal cols: %s", row)
			continue
		}

		key := cols[0]
		val := cols[1]
		switch key {
		case "作者":
			info.Author = val
		case "出版社":
			info.Press = val
		case "出品方":
			info.Producer = val
		case "副标题":
			info.Subtitle = val
		case "原作名":
			info.OriTitle = val
		case "译者":
			info.Translator = val
		case "出版年":
			info.Year = val
		case "页数":
			info.Pages = val
		case "定价":
			info.Price = val
		case "装帧":
			info.Binding = val
		case "丛书":
			info.Series = val
		case "ISBN":
			info.ISBN = val
		case "统一书号":
			info.Number = val
		default:
			logs.Logger.Warn("Discover new properties")
		}
	}
}

type Stars struct {
	One, Two, Three, Four, Five float64
}
