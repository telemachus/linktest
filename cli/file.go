package cli

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

// File stores information and data for a file.
type File struct {
	Name    string
	Content []byte
}

func (app *App) FileGen(files []string, logger *slog.Logger) chan *File {
	if app.NoOp() {
		return nil
	}
	var wg sync.WaitGroup
	downstream := make(chan *File)
	yieldFile := func(file string, fileCh chan *File) {
		defer wg.Done()
		content, err := getContent(file)
		if err != nil {
			app.ExitValue = exitFailure
			logger.LogAttrs(
				context.Background(),
				slog.LevelWarn,
				"Skipping a file because of a problem",
				slog.Group("skip",
					slog.String("file", file),
					slog.Any("error", err),
				),
			)
			return
		}
		fileCh <- &File{
			Name:    file,
			Content: content,
		}
	}
	wg.Add(len(files))
	for _, file := range files {
		f := file
		go yieldFile(f, downstream)
	}
	go func() {
		wg.Wait()
		close(downstream)
	}()
	return downstream
}

func getContent(file string) ([]byte, error) {
	fh, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	doc, err := io.ReadAll(fh)
	if err != nil {
		return []byte{}, err
	}
	return doc, nil
}
