package consumer

import (
	"sync"
)

// DoneChannel represents a channel which is closed once a goroutine is done
type DoneChannel <-chan struct{}

// AndDone joins specified channels into single one. The resulting channel will be done once all channels are done.
func AndDone(chans ...DoneChannel) DoneChannel {
	var wg sync.WaitGroup

	result := make(chan struct{})

	multiplex := func(c DoneChannel) {
		<-c
		wg.Done()
	}

	wg.Add(len(chans))
	for _, ch := range chans {
		go multiplex(ch)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
