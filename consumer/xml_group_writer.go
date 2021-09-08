package consumer

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/goat-project/goat/consumer/wrapper"
)

// NAMESPACE to eu emi storage record
const NAMESPACE = "http://eu-emi.eu/namespaces/2011/02/storagerecord"

type storagesXMLData struct {
	XMLName  xml.Name      `xml:"StorageUsageRecords"`
	XMLNs    string        `xml:"xmlns,attr"`
	XMLNsSr  string        `xml:"xmlns:sr,attr"`
	Storages []interface{} `xml:"StorageUsageRecord"`
}

// XMLGroupWriter converts each record to XML format and writes it to file.
// Multiple records may be written into a single file.
type XMLGroupWriter struct {
	// path to output directory
	outDir string
	// count the records per file
	count uint64
	// slice of records
	recs []interface{}
}

// NewXMLGroupWriter creates a new XMLGroupWriter.
func NewXMLGroupWriter(outputDir string, countPerFile uint64) XMLGroupWriter {
	return XMLGroupWriter{
		outDir: outputDir,
		count:  countPerFile,
		recs:   make([]interface{}, countPerFile),
	}
}

func (xgw XMLGroupWriter) outputDir() string {
	return xgw.outDir
}

func (xgw XMLGroupWriter) countPerFile() uint64 {
	return xgw.count
}

func (xgw XMLGroupWriter) records() []interface{} {
	return xgw.recs
}

func (xgw XMLGroupWriter) save(rec interface{}, index uint64) {
	if int(index) >= len(xgw.recs) {
		// should never happen
		logrus.WithFields(logrus.Fields{"index": index, "length": len(xgw.recs)}).Error("index out of range")
		return
	}

	xgw.recs[index] = rec
}

func (xgw XMLGroupWriter) wrap(data wrapper.RecordWrapper) (interface{}, error) {
	return data.AsXML()
}

func (xgw XMLGroupWriter) convertAndWrite(file io.Writer, countInFile uint64) error {
	if file == nil {
		return fmt.Errorf("unable to write data, file is nil")
	}

	// copy records to ensure the count in file
	newRecords := make([]interface{}, countInFile)
	copy(newRecords, xgw.recs)

	// convert to XML format
	xd, err := xml.MarshalIndent(storagesXMLData{
		XMLNs:    NAMESPACE,
		XMLNsSr:  NAMESPACE,
		Storages: newRecords}, "", " ")
	if err != nil {
		logrus.WithField("error", err).Error("unable to convert to XML format")
		return err
	}

	// write header and XML data to file
	_, err = file.Write([]byte(xml.Header + "\n"))
	if err != nil {
		logrus.WithField("error", err).Error("unable to write XML header")
		return err
	}

	_, err = file.Write(xd)
	return err
}
