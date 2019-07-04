package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/goat-project/goat/consumer/wrapper"
)

type ipsJSONData struct {
	Ips []interface{}
}

// JSONGroupWriter converts each record to json format and writes it to file.
// Multiple records may be written into a single file.
type JSONGroupWriter struct {
	outputDir    string
	countPerFile uint64
	records      []interface{}
}

// NewJSONGroupWriter creates a new JSONGroupWriter.
func NewJSONGroupWriter(outputDir string, countPerFile uint64) JSONGroupWriter {
	return JSONGroupWriter{
		outputDir:    outputDir,
		countPerFile: countPerFile,
		records:      make([]interface{}, countPerFile),
	}
}

func (jgw JSONGroupWriter) writeFile(id string, countInFile, filenameCounter uint64) error {
	newRecords := make([]interface{}, countInFile)
	copy(newRecords, jgw.records)
	jsonData := ipsJSONData{Ips: newRecords}
	filename := path.Join(jgw.outputDir, path.Join(id, fmt.Sprintf(filenameFormat, filenameCounter)))
	// open the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// convert to JSON format
	jd, err := json.MarshalIndent(jsonData, "", " ")
	if err != nil {
		return err
	}

	// write jsonData to file
	_, err = file.Write(jd)
	if err != nil {
		return err
	}

	// close file
	return file.Close()
}

// Consume converts each record to json and writes it to file.
// Multiple records may be written into a single file.
func (jgw JSONGroupWriter) Consume(ctx context.Context, id string,
	records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)

	if err := ensureDirectoryExists(path.Join(jgw.outputDir, id)); err != nil {
		return nil, err
	}

	go func() {
		defer close(res)

		var countInFile, filenameCounter uint64
		countInFile, filenameCounter = 0, 0
		for {
			select {
			case jsonData, ok := <-records:
				if !ok {
					// end of stream
					if countInFile > 0 {
						err := jgw.writeFile(id, countInFile, filenameCounter)
						// but we have something to save!
						if err != nil {
							trySendError(ctx, res, err)
						}
					}
					return
				}

				// convert received record to JSON
				rec, err := jsonData.AsJSON()
				if err != nil {
					trySendError(ctx, res, err)
				}

				// save it for later
				jgw.records[countInFile] = rec

				countInFile++

				// if we already have this many records in the file
				if countInFile == jgw.countPerFile {
					err := jgw.writeFile(id, countInFile, filenameCounter)
					if err != nil {
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
