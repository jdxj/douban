package book

import (
	"douban/modules"
	"douban/utils"
	"douban/utils/logs"
	"strings"
)

type Book struct {
	ID     int64
	Title  string
	Author string
	Press  string

	Info    *Info
	Opinion *modules.Opinion
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
	rawRows := utils.CleanAndSplit(rawInfo)
	rawRows = info.clearColon(rawRows)
	logs.Logger.Debug("length: %d, rawRows: %v", len(rawRows), rawRows)

	var rows []string
	var kv string
	for i, meta := range rawRows {
		if i == 0 && kv == "" { // 初始化第一个
			kv = meta
			continue
		}

		if !strings.HasSuffix(meta, ":") {
			kv += meta
		} else {
			rows = append(rows, kv)
			kv = meta
		}

		if i == len(rawRows)-1 { // 追加最后一个
			rows = append(rows, kv)
		}
	}
	logs.Logger.Debug("length: %d, rows: %v", len(rows), rows)

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
			logs.Logger.Warn("Discover new properties: %s", key)
		}
	}
}

// clearColon 用于清除属性值中的冒号, 防止与属性名后面跟随的冒号冲突.
// 例: `原作名: Principles: Life and Work` 将变成: `原作名: Principles@ Life and Work`
func (info *Info) clearColon(rows []string) []string {
	var last bool
	for i, meta := range rows {
		if i == 0 {
			// 默认第一个是带有冒号的属性名
			last = true
			continue
		}

		if last && strings.HasSuffix(meta, ":") {
			rows[i] = strings.ReplaceAll(meta, ":", "@")
			last = false
		} else {
			last = strings.HasSuffix(meta, ":")
		}
	}

	return rows
}
