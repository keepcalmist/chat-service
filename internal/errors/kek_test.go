package errors

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	timeStart := time.Now()
	ch1, ch2 := worker(), worker()
	_, _ = <-ch1, <-ch2

	t.Log(time.Now().Sub(timeStart))
}

func worker() chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		time.Sleep(3 * time.Second)
		ch <- 1
	}()
	return ch
}
