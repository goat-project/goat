package consumer

import (
	"sync"
)

type ConsumerResult struct {
	ok  bool
	err error
}

type ResultsChannel <-chan ConsumerResult
type ErrHandler func(error)

func (cr ConsumerResult) Error() error {
	return cr.err
}

func (cr ConsumerResult) IsOK() bool {
	return cr.ok
}

func NewErrorResult(err error) ConsumerResult {
	return ConsumerResult{
		err: err,
		ok:  false,
	}
}

func NewOkResult() ConsumerResult {
	return ConsumerResult{
		err: nil,
		ok:  true,
	}
}

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
