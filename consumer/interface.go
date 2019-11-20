package consumer

import (
	"context"

	"github.com/goat-project/goat/consumer/wrapper"
)

// Interface (consumer Interface) to processes records.
type Interface interface {
	// Consume processes all incoming records from a channel, transforms them to a specific format
	// and saves them to the file. The directory is specified by a client id.
	// The format of the file is specified based on the incoming records. The virtual machine records are converted
	// to template format, template is specified in vm.tmpl file. The IP records are converted to JSON format and
	// the storage records are converted to XML format.
	Consume(context.Context, string, chan bool, <-chan wrapper.RecordWrapper) (ResultsChannel, error)
}
