package consumer

import (
	"context"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/goat-project/goat/consumer/wrapper"
)

const (
	templateFileExtension = "tmpl"
)

// TemplateGroupWriter converts each record to template and writes it to file. Multiple records may be written into a single file.
type TemplateGroupWriter struct {
	dir          string
	templatesDir string
	countPerFile uint64
	templateName string
	outExtension string
	template     *template.Template
}

// NewTemplateGroupWriter creates a new TemplateGroupWriter.
func NewTemplateGroupWriter(dir, templatesDir, templateName, outExtension string, countPerFile uint64) TemplateGroupWriter {
	return TemplateGroupWriter{
		dir:          dir,
		templatesDir: templatesDir,
		countPerFile: countPerFile,
		templateName: templateName,
		outExtension: outExtension,
		template:     nil,
	}
}

func ensureDirectoryExists(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil && err != os.ErrExist {
		return err
	}

	return nil
}

func (tgw *TemplateGroupWriter) initTemplate() error {
	if tgw.template != nil {
		return nil
	}

	template, err := template.ParseGlob(path.Join(tgw.templatesDir, fmt.Sprintf("*.%s", templateFileExtension)))
	tgw.template = template
	return err
}

func trySendError(ctx context.Context, res chan<- Result, err error) {
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

// Consume converts each record to template and writes it to file. Multiple records may be written into a single file.
func (tgw TemplateGroupWriter) Consume(ctx context.Context, id string, records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)

	if err := ensureDirectoryExists(path.Join(tgw.dir, id)); err != nil {
		return nil, err
	}

	if err := tgw.initTemplate(); err != nil {
		return nil, err
	}

	// open the initial file
	file, err := os.Create(path.Join(tgw.dir, path.Join(id, fmt.Sprintf("0.%s", tgw.outExtension))))
	if err != nil {
		return nil, err
	}

	go func() {

		defer func() {
			err := file.Close()
			trySendError(ctx, res, err)

			// res must be closed after the file, as it might be used to send the last error
			close(res)
		}()

		var countInFile, filenameCounter uint64
		countInFile, filenameCounter = 0, 0
		for {
			select {
			case templateData, ok := <-records:
				if !ok {
					// end of stream
					return
				}
				// write templateData to file
				err := tgw.template.ExecuteTemplate(file, tgw.templateName, templateData)
				if err != nil {
					trySendError(ctx, res, err)
					return
				}

				countInFile++

				// if we already have this many records in the file
				if countInFile == tgw.countPerFile {

					// close file
					err = file.Close()
					if err != nil {
						trySendError(ctx, res, err)
						// exit on error
						return
					}

					// increase filename counter
					filenameCounter++

					// open next file
					file, err = os.Create(path.Join(tgw.dir, path.Join(id, fmt.Sprintf("%d.%s", filenameCounter, tgw.outExtension))))
					if err != nil {
						trySendError(ctx, res, err)
						// exit on error
						return
					}

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
