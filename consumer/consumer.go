package consumer

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/sirupsen/logrus"
)

// Consumer structure implements consumer Interface. The writer (template, JSON, XML) is specified by writer interface.
type Consumer struct {
	writer writer
}

// writer interface to differ writer implementations.
type writer interface {
	// get output directory
	outputDir() string
	// get count of records per file
	countPerFile() uint64
	// get slice of records
	records() []interface{}
	// save/set record to the slice of records
	save(interface{}, uint64)
	// wrap record to specific structure
	wrap(wrapper.RecordWrapper) (interface{}, error)
	// convert record to the specific format and write to the file
	convertAndWrite(io.Writer, uint64) error
}

// CreateConsumer creates Consumer with writer which specifies structure and format.
func CreateConsumer(w writer) *Consumer {
	return &Consumer{
		writer: w,
	}
}

// Consume processes all incoming records from a channel, transforms them to a specific format
// and saves them to the file. The directory is specified by a client id.
// The format of the file is specified based on the incoming records. The virtual machine records are converted
// to template format, template is specified in vm.tmpl file. The IP records are converted to JSON format and
// the storage records are converted to XML format.
func (c Consumer) Consume(ctx context.Context, id string, endOfWriting chan bool,
	records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)

	// ensure directory exists
	if err := ensureDirectoryExists(path.Join(c.writer.outputDir(), id)); err != nil {
		logrus.WithFields(logrus.Fields{"id": id, "error": err}).Error("unable to ensure directory exists")
		return nil, err
	}

	go func() {
		defer close(res)

		countInFile, filenameCounter := uint64(0), uint64(0)

		for {
			select {
			case data, ok := <-records:
				if !ok {
					// end of stream
					if countInFile > 0 {
						// but we have something to save!
						if err := c.writeFile(id, countInFile, filenameCounter); err != nil {
							trySendError(ctx, res, err)
						}
					}

					// signal the end of writing
					endOfWriting <- true
					return
				}

				// wrap received record
				rec, err := c.writer.wrap(data)
				if err != nil {
					trySendError(ctx, res, err)
				}

				// save wrapped data for later
				c.writer.save(rec, countInFile)

				countInFile++

				// if we already have this many records in the file
				if countInFile == c.writer.countPerFile() {
					if err := c.writeFile(id, countInFile, filenameCounter); err != nil {
						trySendError(ctx, res, err)
					}

					// increase filename counter
					filenameCounter++

					// reset record in file counter
					countInFile = 0
				}
			case <-ctx.Done():
				// goroutine has been canceled
				return
			}
		}
	}()

	return res, nil
}

// writeFile opens the file, converts records to specified format and writes them to the file.
func (c Consumer) writeFile(id string, countInFile, filenameCounter uint64) error {
	file, err := openFile(c.writer.outputDir(), id, filenameCounter)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("unable to create file")
		return err
	}

	if err = c.writer.convertAndWrite(file, countInFile); err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "id": id}).Error("unable to write data to file")
		return err
	}

	logrus.WithFields(logrus.Fields{"id": id, "file-name": file.Name(), "count-in-file": countInFile}).Info("write file")

	// close file
	return file.Close()
}

func ensureDirectoryExists(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && err != os.ErrExist {
		return err
	}

	return nil
}

// openFile creates unique name for the file and opens that file.
func openFile(outputDir, id string, filenameCounter uint64) (*os.File, error) {
	filename := path.Join(outputDir, path.Join(id, fmt.Sprintf(filenameFormat, filenameCounter)))

nameFinding:
	for {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break nameFinding
		}

		filenameCounter++
		filename = path.Join(outputDir, path.Join(id, fmt.Sprintf(filenameFormat, filenameCounter)))
	}

	return os.Create(filename)
}

// In case of the error on the server side, trySendError informs the client about that error.
func trySendError(ctx context.Context, res chan<- Result, err error) {
	logrus.WithField("error", err).Error("unable to finish")

	for {
		select {
		case res <- NewErrorResult(err):
			// error sent successfully
			return
		case <-ctx.Done():
			// goroutine has been canceled
			return
		}
	}
}
