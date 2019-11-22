package book

import (
	"douban/modules"
	"douban/utils"
	"douban/utils/logs"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	// 标签页
	TagsPage     = "https://book.douban.com/tag/?view=cloud"
	TagURLPrefix = "https://book.douban.com"
	PageLimit    = 400
	PauseDur     = 5 * time.Second
)

func NewWormhole() *Wormhole {
	w := &Wormhole{}
	return w
}

type Wormhole struct {
}

func (w *Wormhole) CaptureTags() {
	// todo: 是否要写入数据库中?
	var tagURLs []string

	req, err := modules.NewRequest("GET", TagsPage, nil)
	if err != nil {
		logs.Logger.Error("can not create tags page req: %s", err)
		return
	}

	client := modules.GenHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		logs.Logger.Error("access tags page fail: %s", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Logger.Error("can not create goquery from tags page%s", err)
		return
	}

	sel := doc.Find(".tagCol").Find("tbody a")
	sel.Each(func(i int, sel *goquery.Selection) {
		if suffix, ok := sel.Attr("href"); ok {
			href := TagURLPrefix + suffix

			tagURLs = append(tagURLs, href)
		}
	})
	resp.Body.Close()

	pause(PauseDur)

	for _, url := range tagURLs {
		logs.Logger.Debug("CaptureBookURL: %s", url)
		w.CaptureBookURL(url)
	}

	// todo: 处理中断信号
	utils.DB.Close()
	logs.Logger.Info("finish")
}

func (w *Wormhole) CaptureBookURL(tagURL string) {
	stmtInsert, err := utils.DB.Prepare("insert into book_url (url) values (?)")
	if err != nil {
		logs.Logger.Error("prepare 'insert into book_url (url) values (?)' fail: %s", err)
		return
	}
	defer stmtInsert.Close()

	stmtQuery, err := utils.DB.Prepare("select id from book_url where url=?")
	if err != nil {
		logs.Logger.Error("prepare 'select id from book_url where url=?' fail: %s", err)
		return
	}
	defer stmtQuery.Close()

	client := modules.GenHTTPClient()

	for i := 0; i < 20*PageLimit; i += 20 {
		url := tagURL + fmt.Sprintf("?start=%d&type=T", i)
		logs.Logger.Debug("in cap book url:", url)

		req, err := modules.NewRequest("GET", url, nil)
		if err != nil {
			logs.Logger.Error("can not create '%s'req: %s", url, err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			logs.Logger.Error("access '%s' fail: %s", url, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Logger.Error("can not create '%s' goquery: %s", err)
			continue
		}

		var mark bool
		doc.Find(".nbg").Each(func(i int, sel *goquery.Selection) {
			mark = true

			if bookURL, ok := sel.Attr("href"); ok {
				// 去重
				rows, err := stmtQuery.Query(bookURL)
				if err != nil {
					logs.Logger.Error("查重失败: %s", err)
					return
				}
				defer rows.Close()

				if rows.Next() {
					logs.Logger.Debug("found 重复 url: %s", bookURL)
					return
				}

				if _, err = stmtInsert.Exec(bookURL); err != nil {
					logs.Logger.Error("insert fail, url: %s, err: %s", bookURL, err)
				}
			}
		})
		resp.Body.Close()

		if !mark {
			break
		}

		pause(PauseDur)
	}
}

func pause(dur time.Duration) {
	time.Sleep(dur)
}
