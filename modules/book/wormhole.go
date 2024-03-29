package book

import (
	"douban/modules"
	"douban/modules/types"
	"douban/utils"
	"douban/utils/logs"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	// 标签页
	TagsPage     = "https://book.douban.com/tag/?view=cloud"
	TagURLPrefix = "https://book.douban.com"
	PageLimit    = 400
)

const (
	CaptureBookURLMode = iota + 1
	CaptureBookMode
)

func NewWormhole() (*Wormhole, error) {
	modeConf, err := utils.GetModeConf()
	if err != nil {
		return nil, err
	}
	w := &Wormhole{
		Mode:   modeConf.Mode,
		signal: make(chan int),
	}
	return w, nil
}

type Wormhole struct {
	Mode   int
	signal chan int
}

func (w *Wormhole) Run() {
	defer close(w.signal)

	switch w.Mode {
	case CaptureBookURLMode:
		w.CaptureTags()
	case CaptureBookMode:
		// 先启动监控 goroutine
		go w.monitor()
		w.CaptureBook()
	default:
		logs.Logger.Warn("Invalid mode of operation")
	}
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

	logs.Logger.Info("Capture finish")
}

func (w *Wormhole) CaptureBookURL(tagURL string) {
	stmtBookInsert, err := utils.DB.Prepare("insert into book_url (url) values (?)")
	if err != nil {
		logs.Logger.Error("Prepare 'insert into book_url (url) values (?)' failed: %s", err)
		return
	}
	defer stmtBookInsert.Close()

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

				if _, err = stmtBookInsert.Exec(bookURL); err != nil {
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

// CaptureBook 该方法需要从数据库中读取 book url
func (w *Wormhole) CaptureBook() {
	rows, err := utils.DB.Query("select count(*) from book_url")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}

	var total int64
	for rows.Next() {
		if err = rows.Scan(&total); err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}
	rows.Close()

	if total <= 0 {
		logs.Logger.Warn("No book url available")
		return
	}

	row, err := utils.DB.Query("select log from log where type=1 limit 1")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}

	var start int64 = -1
	for row.Next() {
		if err = row.Scan(&start); err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}
	row.Close()

	if start < 0 { //从没运行过
		_, err = utils.DB.Exec("insert into log (log, type) values (0, 1)")
		if err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}
	start++

	stmtLogUpdate, err := utils.DB.Prepare("update log set log=? where type=1")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer stmtLogUpdate.Close()

	stmtQuery, err := utils.DB.Prepare("select url from book_url order by id limit ?,?")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer stmtQuery.Close()

	stmtBookInsert, err := utils.DB.Prepare("insert into book (title, author, press) values (?,?,?)")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer stmtBookInsert.Close()

	stmtOpiInsert, err := utils.DB.Prepare("insert into opinion (score, amount, one, two, three, four, five, type, ref) values (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer stmtOpiInsert.Close()

	client := modules.GenHTTPClient()
	for i := start; i < total; i++ {
		utils.Pause(utils.Pause5s)

		// 只有一行
		row, err := stmtQuery.Query(i, 1)
		if err != nil {
			logs.Logger.Error("%s", err)
			// todo: 是否继续?
			return
		}

		var url string
		for row.Next() {
			if err = row.Scan(&url); err != nil {
				logs.Logger.Error("%s", err)
				row.Close()
				return
			}
		}
		row.Close()

		if url == "" {
			continue
		}

		book, err := w.genBook(url, client)
		if err != nil {
			logs.Logger.Error("Failed to generate book, url: %s, error: %s", url, err)
			continue
		}

		result, err := stmtBookInsert.Exec(book.ToInsert()...)
		if err != nil {
			logs.Logger.Error("Insert into table failed, url: %s, error: %s", url, err)
			continue
		}

		book.ID, err = result.LastInsertId()
		if err != nil {
			logs.Logger.Error("Get last insert id failed, err: %s", err)
			continue
		}

		_, err = stmtOpiInsert.Exec(book.Opinion.ToInsert()...)
		if err != nil {
			logs.Logger.Error("Insert into opinion failed, err: %s", err)
			continue
		}

		// book, opinion 都插入后更新记录
		if _, err = stmtLogUpdate.Exec(i); err != nil {
			logs.Logger.Error("Update log failed, err: %s", err)
		}
	}
}

func (w *Wormhole) genBook(url string, client *http.Client) (*Book, error) {
	req, err := modules.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// book
	book := new(Book)

	rawTitle := doc.Find("h1").Text()
	title := utils.CleanAndJoin(rawTitle, "")
	if title == "" {
		return nil, fmt.Errorf("found invalid title, url: %s", url)
	}
	book.Title = title

	// info
	w.genInfo(book, doc)
	// opinion
	w.genOpinion(book, doc)

	return book, nil
}

func (w *Wormhole) genInfo(book *Book, doc *goquery.Document) {
	info := new(Info)
	book.Info = info

	rawInfo := doc.Find("#info").Text()
	info.Unmarshal(rawInfo)

	book.Author = info.Author
	book.Press = info.Press
}

func (w *Wormhole) genOpinion(book *Book, doc *goquery.Document) {
	opinion := new(modules.Opinion)
	book.Opinion = opinion

	rawScore := doc.Find(".rating_num").Text()
	rawScore = utils.CleanAndJoin(rawScore, "")
	score, err := strconv.ParseFloat(rawScore, 64)
	if err != nil {
		logs.Logger.Warn("rawScore: %s, err: %s", rawScore, err)
	}
	opinion.Score = score

	rawAmount := doc.Find(".rating_people").Text()
	rawAmount = utils.CleanAndJoin(rawAmount, "")
	rawAmount = strings.ReplaceAll(rawAmount, "人评价", "")
	amount, err := strconv.ParseInt(rawAmount, 10, 64)
	if err != nil {
		logs.Logger.Warn("rawAmount: %s, err: %s", rawAmount, err)
	}
	opinion.Amount = amount

	opinion.Type = types.B
	opinion.Ref = &book.ID

	doc.Find(".rating_per").Each(func(i int, sel *goquery.Selection) {
		rawPer := utils.CleanAndJoin(sel.Text(), "")
		rawPer = strings.ReplaceAll(rawPer, "%", "")
		per, err := strconv.ParseFloat(rawPer, 64)
		if err != nil {
			logs.Logger.Warn("rawPer: %s, err: %s", rawPer, err)
		}

		switch i {
		case 0:
			opinion.Five = per
		case 1:
			opinion.Four = per
		case 2:
			opinion.Three = per
		case 3:
			opinion.Two = per
		case 4:
			opinion.One = per
		}
	})
}

// monitor 主要用于监控程序运行状态, 为了能够及时发现程序异常.
// 同时定时统计一些所抓取数据, 来分析运行效果.
func (w *Wormhole) monitor() {
	// 首次启动通知一次
	go w.takeALook()

	// 每天8点, 20点统计并通知
	now := time.Now()
	today8AM := time.Unix(now.Unix()/86400*86400, 0)
	tom8AM := today8AM.Add(24 * time.Hour)
	dur := tom8AM.Sub(now)

	timer := time.NewTimer(dur)
	defer timer.Stop()

	for {
		select {
		case <-w.signal:
			logs.Logger.Info("Discovery process has ended")
			return
		case <-timer.C:
			go w.takeALook()
			timer.Reset(12 * time.Hour)
		}
	}
}

func (w *Wormhole) takeALook() {
	rows, err := utils.DB.Query("select count(*) from book")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer rows.Close()

	var amount uint64
	for rows.Next() {
		if err = rows.Scan(&amount); err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}

	rows, err = utils.DB.Query("select log from log where type=1")
	if err != nil {
		logs.Logger.Error("%s", err)
		return
	}
	defer rows.Close()

	var progress uint64
	for rows.Next() {
		if err = rows.Scan(&progress); err != nil {
			logs.Logger.Error("%s", err)
			return
		}
	}

	subject := "book 抓取状态"
	format := "当前 book 表中有 %d 条数据\n"
	format += "当前进度为: %d\n"
	utils.SendEmail(subject, format, amount, progress)
}
