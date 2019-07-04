package consumer

import (
	"context"

	"github.com/goat-project/goat/consumer/wrapper"
)

type storagesXMLData struct {
	Storages []interface{}
}

// XMLGroupWriter converts each record to XML format and writes it to file.
// Multiple records may be written into a single file.
type XMLGroupWriter struct {
	outputDir    string
	countPerFile uint64
	records      []interface{}
}

// NewXMLGroupWriter creates a new XMLGroupWriter.
func NewXMLGroupWriter(outputDir string, countPerFile uint64) XMLGroupWriter {
	return XMLGroupWriter{
		outputDir:    outputDir,
		countPerFile: countPerFile,
		records:      make([]interface{}, countPerFile),
	}
}

// Consume converts each record to XML and writes it to file.
// Multiple records may be written into a single file.
func (xgw XMLGroupWriter) Consume(ctx context.Context, id string,
	records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	// TODO implement consume for storage records
	return nil, nil
}
