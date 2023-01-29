package cli

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func (app *App) Display(statusCh <-chan *Link, logger kitlog.Logger) {
	if app.NoOp() {
		return
	}

	if app.Verbose {
		displayAll(statusCh, logger)

		return
	}

	for link := range statusCh {
		switch {
		case link.Err != nil:
			level.Error(logger).Log(
				"file", link.File,
				"URL", link.URL,
				"status", link.Status,
				"error", link.Err,
			)
		case link.Status != 200:
			level.Warn(logger).Log(
				"file", link.File,
				"URL", link.URL,
				"status", link.Status,
				"error", link.Err,
			)
		}
	}
}

func displayAll(statusCh <-chan *Link, errLogger kitlog.Logger) {
	outLogger := newLogger(os.Stdout)

	for link := range statusCh {
		switch {
		case link.Err != nil:
			level.Error(errLogger).Log(
				"file", link.File,
				"URL", link.URL,
				"status", link.Status,
				"error", link.Err,
			)
		case link.Status != 200:
			level.Warn(errLogger).Log(
				"file", link.File,
				"URL", link.URL,
				"status", link.Status,
				"error", link.Err,
			)
		default:
			level.Info(outLogger).Log(
				"file", link.File,
				"URL", link.URL,
				"status", link.Status,
				"error", link.Err,
			)
		}
	}
}
