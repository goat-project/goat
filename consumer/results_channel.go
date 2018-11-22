package consumer

import (
	"sync"
)

// Result encapsulates a result of consume method. It is either erroneous result or successful
type Result struct {
	ok  bool
	err error
}

// ResultsChannel is a read-only channel of Results
type ResultsChannel <-chan Result

// ErrHandler is a function that handles specified error
type ErrHandler func(error)

// Error returns the underlying error in case this result is successful. Otherwise returns nil
func (cr Result) Error() error {
	return cr.err
}

// IsOK returns true if and only if this result is successful
func (cr Result) IsOK() bool {
	return cr.ok
}

// NewErrorResult creates a new unsuccessful result with specified underlying error
func NewErrorResult(err error) Result {
	return Result{
		err: err,
		ok:  false,
	}
}

// NewOkResult creates a new successful result
func NewOkResult() Result {
	return Result{
		err: nil,
		ok:  true,
	}
}

// NewResultFromError creates an unsuccessful result in case specified error is not nil. If the error is nil, it creates a successful result
func NewResultFromError(e error) Result {
	if e != nil {
		return NewErrorResult(e)
	}
	return NewOkResult()
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
