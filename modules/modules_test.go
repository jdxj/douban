package modules

import (
	"fmt"
	"testing"
	"time"
)

func TestReadChan(t *testing.T) {
	c := make(chan int)

	go func() {
		defer close(c)

		for i := 0; i < 100; i++ {
			c <- i
		}
	}()

	for {
		select {
		case data, ok := <-c:
			if ok {
				fmt.Println(data)
			} else {
				fmt.Println(data)
				time.Sleep(time.Second)
			}
		}
	}
}
