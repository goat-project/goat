package consumer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/sirupsen/logrus"
)

type ipsJSONData struct {
	Ips []interface{}
}

// JSONGroupWriter converts each record to json format and writes it to file.
// Multiple records may be written into a single file.
type JSONGroupWriter struct {
	// path to output directory
	outDir string
	// count the records per file
	count uint64
	// slice of records
	recs []interface{}
}

// NewJSONGroupWriter creates a new JSONGroupWriter.
func NewJSONGroupWriter(outputDir string, countPerFile uint64) JSONGroupWriter {
	return JSONGroupWriter{
		outDir: outputDir,
		count:  countPerFile,
		recs:   make([]interface{}, countPerFile),
	}
}

func (jgw JSONGroupWriter) outputDir() string {
	return jgw.outDir
}

func (jgw JSONGroupWriter) countPerFile() uint64 {
	return jgw.count
}

func (jgw JSONGroupWriter) records() []interface{} {
	return jgw.recs
}

func (jgw JSONGroupWriter) save(rec interface{}, index uint64) {
	if int(index) >= len(jgw.recs) {
		// should never happen
		logrus.WithFields(logrus.Fields{"index": index, "length": len(jgw.recs)}).Error("index out of range")
		return
	}

	jgw.recs[index] = rec
}

func (jgw JSONGroupWriter) wrap(data wrapper.RecordWrapper) (interface{}, error) {
	return data.AsJSON()
}

func (jgw JSONGroupWriter) convertAndWrite(file io.Writer, countInFile uint64) error {
	if file == nil {
		return fmt.Errorf("unable to write data, file is nil")
	}

	// copy records to ensure the count in file
	newRecords := make([]interface{}, countInFile)
	copy(newRecords, jgw.recs)

	// convert to JSON format
	jd, err := json.MarshalIndent(ipsJSONData{Ips: newRecords}, "", " ")
	if err != nil {
		logrus.WithField("error", err).Error("unable to convert to JSON format")
		return err
	}

	// write JSON data to file
	_, err = file.Write(jd)
	return err
}
