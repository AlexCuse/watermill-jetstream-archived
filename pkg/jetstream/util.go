package jetstream

import (
	"sync"
	"time"
)

func waitGroupTimeout(wg *sync.WaitGroup, timeout *time.Ticker) bool {
	wgClosed := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		wgClosed <- struct{}{}
	}()

	select {
	case <-wgClosed:
		return false
	case <-timeout.C:
		return true
	}
}
