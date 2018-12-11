package importer

import (
	"github.com/goat-project/goat/consumer/wrapper"
)

// RecordStream is a stream of records. ReceiveIdentifier must be called first, then Receive can be called
type RecordStream interface {
	// ReceiveIdentifier receives client identifier from given stream.
	// If the first message is not client identifier, ErrNonFirstClientIdentifier is returned
	ReceiveIdentifier() (string, error)

	// Receive receives the next record from the stream.
	// If there is no more records in the stream, io.EOF is returned.
	// If an identifier is received, ErrNonFirstClientIdentifier is returned.
	Receive() (wrapper.RecordWrapper, error)

	// Close closes the stream
	Close() error
}
