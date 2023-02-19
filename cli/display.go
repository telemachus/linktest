package cli

import (
	"fmt"
)

func (app *App) Display(statusCh <-chan *Link) {
	if app.NoOp() {
		return
	}

	if app.Verbose {
		displayAll(statusCh)

		return
	}

	for link := range statusCh {
		switch {
		case link.Err != nil, link.Status != 200:
			msg := fmt.Sprintf(
				"%s: %q: %q: %d: %v",
				appName,
				link.File,
				link.URL,
				link.Status,
				link.Err,
			)
			advise(msg)
		}
	}
}

func displayAll(statusCh <-chan *Link) {
	for link := range statusCh {
		switch {
		case link.Err != nil, link.Status != 200:
			msg := fmt.Sprintf(
				"%s: %q: %q: %d: %v",
				appName,
				link.File,
				link.URL,
				link.Status,
				link.Err,
			)
			advise(msg)
		default:
			msg := fmt.Sprintf(
				"%s: %q: %q: %d: %v",
				appName,
				link.File,
				link.URL,
				link.Status,
				link.Err,
			)
			inform(msg)
		}
	}
}
