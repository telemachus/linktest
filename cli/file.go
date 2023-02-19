package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// File stores information and data for a file.
type File struct {
	Name    string
	Content []byte
}

func (app *App) FileGen(files []string) chan *File {
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

			msg := fmt.Sprintf("%s: skipping [%s]: %v", appName, file, err)
			advise(msg)

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
