package consumer

import (
	"sync"
)

// ConsumerResult encapsulates a result of consume method. It is either erroneous result or successful
type ConsumerResult struct {
	ok  bool
	err error
}

// ResultsChannel is a read-only channel of ConsumerResults
type ResultsChannel <-chan ConsumerResult

// ErrHandler is a function that handles specified error
type ErrHandler func(error)

// Error returns the underlying error in case this result is successful. Otherwise returns nil
func (cr ConsumerResult) Error() error {
	return cr.err
}

// IsOK returns true if and only if this result is successful
func (cr ConsumerResult) IsOK() bool {
	return cr.ok
}

// NewErrorResult creates a new unsuccessful result with specified underlying error
func NewErrorResult(err error) ConsumerResult {
	return ConsumerResult{
		err: err,
		ok:  false,
	}
}

// NewOkResult creates a new successful result
func NewOkResult() ConsumerResult {
	return ConsumerResult{
		err: nil,
		ok:  true,
	}
}

// CheckResults reads all results from given results channels. When an error is encountered, onError is called with the underlying error
func CheckResults(onError ErrHandler, chans ...ResultsChannel) {
	var wg sync.WaitGroup
	multiplex := func(c ResultsChannel) {
		for r := range c {
			if !r.IsOK() {
				onError(r.Error())
			}
		}
		wg.Done()
	}

	wg.Add(len(chans))
	for _, ch := range chans {
		go multiplex(ch)
	}

	wg.Wait()
}
