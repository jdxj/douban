package modules

import "net/http"

func GenHTTPClient() *http.Client {
	// todo: cookie
	c := &http.Client{}
	return c
}
