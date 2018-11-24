package consumer

import (
	"context"
	"github.com/goat-project/goat/consumer/wrapper"
	"io"
	"os"
	"path"
	"text/template"
)

// WriterConsumer processes accounting data by transforming them according to supplied templates and subsequently writing to a file
type WriterConsumer struct {
	dir          string
	templatesDir string
	template     *template.Template
}

// NewWriter creates a new WriterConsumer
func NewWriter(dir, templatesDir string) WriterConsumer {
	return WriterConsumer{
		dir:          dir,
		templatesDir: templatesDir,
	}
}

func ensureDirectoryExists(path string) error {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil && err != os.ErrExist {
		return err
	}

	return nil
}

func (wc *WriterConsumer) initTemplates() error {
	if wc.template != nil {
		return nil
	}
	template, err := template.ParseGlob(path.Join(wc.templatesDir, "*.tmpl"))
	wc.template = template
	return err
}

func (wc WriterConsumer) processTo(record interface{}, wr io.Writer) error {
	return wc.template.Execute(wr, record)
}

func (wc WriterConsumer) write(id string, rw wrapper.RecordWrapper) error {

	file, err := os.Open(path.Join(path.Join(wc.dir, id), rw.Filename()))
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	if err != nil {
		return err
	}

	if err = wc.processTo(rw.AsTemplate(), file); err != nil {
		return err
	}

	return err
}

// Consume transforms records into text and writes them to a subdirectory of dir specified by WriterConsumer's dir field. Each record is written into its own file.
func (wc WriterConsumer) Consume(ctx context.Context, id string, records <-chan wrapper.RecordWrapper) (ResultsChannel, error) {
	res := make(chan Result)

	if err := ensureDirectoryExists(path.Join(wc.dir, id)); err != nil {
		return nil, err
	}

	if err := wc.initTemplates(); err != nil {
		return nil, err
	}

	go func() {
		defer close(res)

		for {
			select {
			case rec := <-records:
				res <- NewResultFromError(wc.write(id, rec))
			case <-ctx.Done():
				return
			}
		}

	}()

	return res, nil
}
