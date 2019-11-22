package utils

import (
	"time"
)

const (
	Pause2s = 2 * time.Second
	Pause5s = 5 * time.Second
)

func Pause(dur time.Duration) {
	time.Sleep(dur)
}
