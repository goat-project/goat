package consumer

import (
	"fmt"
	"io"
	"path"
	"text/template"

	"github.com/goat-project/goat/consumer/wrapper"
	"github.com/sirupsen/logrus"
)

const (
	templateFileName = "vm.tmpl"
	templateName     = "VMS"
	filenameFormat   = "%014d"
)

type vmsTemplateData struct {
	Vms []interface{}
}

// TemplateGroupWriter converts each record to template and writes it to file.
// Multiple records may be written into a single file.
type TemplateGroupWriter struct {
	// path to output directory
	outDir string
	// count the records per file
	count uint64
	// slice of records
	recs []interface{}
	// path to template directory
	templatesDir string
}

// NewTemplateGroupWriter creates a new TemplateGroupWriter.
func NewTemplateGroupWriter(outputDir, templatesDir string, countPerFile uint64) TemplateGroupWriter {
	return TemplateGroupWriter{
		outDir:       outputDir,
		count:        countPerFile,
		recs:         make([]interface{}, countPerFile),
		templatesDir: templatesDir,
	}
}

func (tgw TemplateGroupWriter) outputDir() string {
	return tgw.outDir
}

func (tgw TemplateGroupWriter) countPerFile() uint64 {
	return tgw.count
}

func (tgw TemplateGroupWriter) records() []interface{} {
	return tgw.recs
}

func (tgw TemplateGroupWriter) save(rec interface{}, index uint64) {
	if int(index) >= len(tgw.recs) {
		// should never happen
		logrus.WithFields(logrus.Fields{"index": index, "length": len(tgw.recs)}).Error("index out of range")
		return
	}

	tgw.recs[index] = rec
}

func (tgw TemplateGroupWriter) wrap(data wrapper.RecordWrapper) (interface{}, error) {
	return data.AsTemplate()
}

func (tgw TemplateGroupWriter) convertAndWrite(file io.Writer, countInFile uint64) error {
	if file == nil {
		return fmt.Errorf("unable to write data, file is nil")
	}

	// copy records to ensure the count in file
	newRecords := make([]interface{}, countInFile)
	copy(newRecords, tgw.recs)

	// initialize template
	tmpl, err := template.ParseGlob(path.Join(tgw.templatesDir, templateFileName))
	if err != nil {
		logrus.WithFields(logrus.Fields{"template-name": templateName, "error": err}).Error("unable to initialize template")
		return err
	}

	// write template data to file
	return tmpl.ExecuteTemplate(file, templateName, vmsTemplateData{Vms: newRecords})
}
