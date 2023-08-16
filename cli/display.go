package cli

import (
	"context"
	"log/slog"
	"os"
)

func (app *App) Display(statusCh <-chan *Link, logger *slog.Logger) {
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
			logger.LogAttrs(
				context.Background(),
				slog.LevelError,
				"Error checking link",
				slog.String("file", link.File),
				slog.Any("error", link.Err),
				slog.Group("response",
					slog.String("URL", link.URL),
					slog.Int("status", link.Status),
				),
			)
		case link.Status != 200:
			logger.LogAttrs(
				context.Background(),
				slog.LevelWarn,
				"Non-200 response to link",
				slog.String("file", link.File),
				slog.Group("response",
					slog.String("URL", link.URL),
					slog.Int("status", link.Status),
				),
				slog.Any("error", link.Err),
			)
		}
	}
}

func displayAll(statusCh <-chan *Link, logger *slog.Logger) {
	outLogger := newLogger(os.Stdout)
	for link := range statusCh {
		switch {
		case link.Err != nil:
			logger.LogAttrs(
				context.Background(),
				slog.LevelError,
				"Error checking link",
				slog.String("file", link.File),
				slog.Any("error", link.Err),
				slog.Group("response",
					slog.String("URL", link.URL),
					slog.Int("status", link.Status),
				),
			)
		case link.Status != 200:
			logger.LogAttrs(
				context.Background(),
				slog.LevelWarn,
				"Non-200 response to link",
				slog.String("file", link.File),
				slog.Group("response",
					slog.String("URL", link.URL),
					slog.Int("status", link.Status),
				),
				slog.Any("error", link.Err),
			)
		default:
			outLogger.LogAttrs(
				context.Background(),
				slog.LevelWarn,
				"200 response to link",
				slog.String("file", link.File),
				slog.Group("response",
					slog.String("URL", link.URL),
					slog.Int("status", link.Status),
				),
				slog.Any("error", link.Err),
			)
		}
	}
}
