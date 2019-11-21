package modules

func NewWormhole() *Wormhole {
	w := &Wormhole{
		tagURL:  make(chan string, 1),
		bookURL: make(chan string, 100),
	}
	return w
}

// 三个模块, 两个通道
type Wormhole struct {
	tagURL  chan string // 标准标签 url
	bookURL chan string // 书的详情 url
}

// 扔
func (w *Wormhole) ThrowTagURL() chan<- string {
	return w.tagURL
}

// 捞
func (w *Wormhole) FishTagURL() <-chan string {
	return w.tagURL
}

func (w *Wormhole) ThrowBookURL() chan<- string {
	return w.bookURL
}

func (w *Wormhole) FishBookURL() <-chan string {
	return w.bookURL
}
