package consumer

import (
	"context"
	"github.com/goat-project/goat-proto-go"
)

// WriterConsumer processes accounting data by transforming them according to supplied templates and subsequently writing to a file
type WriterConsumer struct {
	dir          string
	templatesDir string
}

// NewWriter creates a new WriterConsumer
func NewWriter(dir, templatesDir string) WriterConsumer {
	return WriterConsumer{
		dir:          dir,
		templatesDir: templatesDir,
	}
}

// ConsumeIps transforms IpRecord-s into text and writes them to a subdirectory of dir specified by WriterConsumer's dir field. Each IpRecord is written to its own file.
func (wc WriterConsumer) ConsumeIps(ctx context.Context, id string, ips <-chan goat_grpc.IpRecord) (DoneChannel, error) {
	done := make(chan struct{})
	return done, nil
}

// ConsumeVms transforms VmRecord-s into text and writes them to a subdirectory of dir specified by WriterConsumer's dir field. Each VmRecord is written to its own file.
func (wc WriterConsumer) ConsumeVms(ctx context.Context, id string, vms <-chan goat_grpc.VmRecord) (DoneChannel, error) {
	done := make(chan struct{})
	return done, nil
}

// ConsumeStorages transforms StorageRecord-s into text and writes them to a subdirectory of dir specified by WriterConsumer's dir field. Each StorageRecord is written to its own file.
func (wc WriterConsumer) ConsumeStorages(ctx context.Context, id string, sts <-chan goat_grpc.StorageRecord) (DoneChannel, error) {
	done := make(chan struct{})
	return done, nil
}
