package modules

import (
	"io"
	"math/rand"
	"net/http"
)

const (
	UserAgent = "User-Agent"
	Accept    = "Accept"
)

var (
	UAName = []string{
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
	}
)

func GenHTTPClient() *http.Client {
	// todo: cookie
	c := &http.Client{}
	return c
}

// todo: 需要更多数据
func RandUserAgent() string {
	idx := rand.Intn(len(UAName))
	return UAName[idx]
}

// NewRequest 是对 http.NewRequest 的封装,
// 主要用于向 request 中添加 header 数据.
func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 添加 User-Agent
	req.Header.Set(UserAgent, RandUserAgent())
	return req, nil
}

func NewRequestAndDo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	client := GenHTTPClient()
	return client.Do(req)
}
