package consumer

import (
	"context"
	"github.com/goat-project/goat/consumer/wrapper"
)

// IPConsumer processes IpRecord-s
type IPConsumer interface {
	// ConsumeIps processes all ip records from specified channel and specified client id
	ConsumeIps(ctx context.Context, id string, ips <-chan wrapper.RecordWrapper) (ResultsChannel, error)
}

// VMConsumer processes VmRecords
type VMConsumer interface {
	// ConsumeVMs processes all ip records from specified channel and specified client id
	ConsumeVms(ctx context.Context, id string, vms <-chan wrapper.RecordWrapper) (ResultsChannel, error)
}

// StorageConsumer processes StorageRecords
type StorageConsumer interface {
	// ConsumeIps processes all ip records from specified channel and specified client id
	ConsumeStorages(ctx context.Context, id string, sts <-chan wrapper.RecordWrapper) (ResultsChannel, error)
}
