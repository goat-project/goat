package importer

import (
	"github.com/goat-project/goat-proto-go"
	"github.com/goat-project/goat/consumer/wrapper"
)

// RecordStream is a stream of records
type RecordStream interface {
	// ReceiveIdentifier receives client identifier from given stream.
	// If the first message is not client identifier, ErrNonFirstClientIdentifier is returned
	ReceiveIdentifier() (string, error)
	// Receive receives the next record from the stream.
	// If there is no more records in the stream, io.EOF is returned.
	Receive() (wrapper.RecordWrapper, error)
	// Close closes the stream
	Close() error
}

// WrapGrpc wraps the GRPC stream using a RecordStream
func WrapGrpc(stream goat_grpc.AccountingService_ProcessServer) RecordStream {
	var rs RecordStream
	return rs
}
