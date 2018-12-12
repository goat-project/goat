package consumer

import (
	"context"
	"github.com/goat-project/goat/consumer/wrapper"
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

// Consume transforms Record-s into text and writes them to a subdirectory of dir specified by WriterConsumer's dir field. Each IpRecord is written to its own file.
func (wc WriterConsumer) Consume(ctx context.Context, id string, records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)
	return res, nil
}
