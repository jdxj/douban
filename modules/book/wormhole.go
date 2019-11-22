package book

import (
	"douban/modules"
	"douban/utils"
	"douban/utils/logs"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	// 标签页
	TagsPage     = "https://book.douban.com/tag/?view=cloud"
	TagURLPrefix = "https://book.douban.com"
	PageLimit    = 400
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

	resp, err := modules.NewRequestAndDo("GET", TagsPage, nil)
	if err != nil {
		logs.Logger.Error("Access tags page failed: %s", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Logger.Error("Can not create goquery from tags page%s", err)
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

	utils.Pause(utils.Pause5s)

	for _, url := range tagURLs {
		logs.Logger.Debug("CaptureBookURL: %s", url)
		w.CaptureBookURL(url)
	}

	// todo: 在中断信号中处理
	utils.DB.Close()
	logs.Logger.Info("Capture finish")
}

func (w *Wormhole) CaptureBookURL(tagURL string) {
	stmtInsert, err := utils.DB.Prepare("insert into book_url (url) values (?)")
	if err != nil {
		logs.Logger.Error("Prepare 'insert into book_url (url) values (?)' failed: %s", err)
		return
	}
	defer stmtInsert.Close()

	stmtQuery, err := utils.DB.Prepare("select id from book_url where url=?")
	if err != nil {
		logs.Logger.Error("Prepare 'select id from book_url where url=?' failed: %s", err)
		return
	}
	defer stmtQuery.Close()

	client := modules.GenHTTPClient()
	for i := 0; i < 20*PageLimit; i += 20 {
		url := tagURL + fmt.Sprintf("?start=%d&type=T", i)
		logs.Logger.Debug("In capture book url:", url)

		req, err := modules.NewRequest("GET", url, nil)
		if err != nil {
			logs.Logger.Error("Can not create '%s' req: %s", url, err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			logs.Logger.Error("Access '%s' fail: %s", url, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			logs.Logger.Error("Can not create '%s' goquery: %s", err)
			continue
		}

		var mark bool
		doc.Find(".nbg").Each(func(i int, sel *goquery.Selection) {
			mark = true

			if bookURL, ok := sel.Attr("href"); ok {
				// 去重
				rows, err := stmtQuery.Query(bookURL)
				if err != nil {
					logs.Logger.Error("Check for duplicate url failure: %s", err)
					return
				}
				defer rows.Close()

				if rows.Next() {
					logs.Logger.Debug("Found duplicate url: %s", bookURL)
					return
				}

				if _, err = stmtInsert.Exec(bookURL); err != nil {
					logs.Logger.Error("Insert failed, url: %s, err: %s", bookURL, err)
				}
			}
		})
		resp.Body.Close()

		if !mark {
			break
		}

		utils.Pause(utils.Pause5s)
	}
}
