package consumer

import (
	"context"
	"github.com/goat-project/goat/consumer/wrapper"
)

// Consumer processes records
type Consumer interface {
	// Consume processes all records from specified channel and specified client id
	Consume(ctx context.Context, id string, records <-chan wrapper.RecordWrapper) (ResultsChannel, error)
}
