package consumer

import (
	"context"
	"github.com/goat-project/goat/consumer/wrapper"
)

// TemplateGroupWriter converts each record to template and writes it to file. Multiple records may be written into a single file.
type TemplateGroupWriter struct {
	dir          string
	templatesDir string
	countPerFile uint64
}

// NewTemplateGroupWriter creates a new TemplateGroupWriter.
func NewTemplateGroupWriter(dir, templatesDir string, countPerFile uint64) TemplateGroupWriter {
	return TemplateGroupWriter{
		dir:          dir,
		templatesDir: templatesDir,
		countPerFile: countPerFile,
	}
}

// Consume converts each record to template and writes it to file. Multiple records may be written into a single file.
func (wc TemplateGroupWriter) Consume(ctx context.Context, id string, records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)
	return res, nil
}
