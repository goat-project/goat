package consumer

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path"

	"github.com/goat-project/goat/consumer/wrapper"
)

type storagesXMLData struct {
	XMLName  xml.Name      `xml:"STORAGES"`
	Storages []interface{} `xml:"STORAGE"`
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
	res := make(chan Result)

	if err := ensureDirectoryExists(path.Join(xgw.outputDir, id)); err != nil {
		return nil, err
	}

	go func() {
		defer close(res)

		var countInFile, filenameCounter uint64
		countInFile, filenameCounter = 0, 0
		for {
			select {
			case xmlData, ok := <-records:
				if !ok {
					// end of stream
					if countInFile > 0 {
						err := xgw.writeFile(id, countInFile, filenameCounter)
						// but we have something to save!
						if err != nil {
							trySendError(ctx, res, err)
						}
					}
					return
				}

				// convert received record to XML
				rec, err := xmlData.AsXML()
				if err != nil {
					trySendError(ctx, res, err)
				}

				// save it for later
				xgw.records[countInFile] = rec

				countInFile++

				// if we already have this many records in the file
				if countInFile == xgw.countPerFile {
					err := xgw.writeFile(id, countInFile, filenameCounter)
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

func (xgw XMLGroupWriter) writeFile(id string, countInFile, filenameCounter uint64) error {
	newRecords := make([]interface{}, countInFile)
	copy(newRecords, xgw.records)
	xmlData := storagesXMLData{Storages: newRecords}
	filename := path.Join(xgw.outputDir, path.Join(id, fmt.Sprintf(filenameFormat, filenameCounter)))
	// open the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// convert to XML format
	jd, err := xml.MarshalIndent(xmlData, "", " ")
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
