package cli

import (
	"io"
	"log/slog"

	"github.com/telemachus/humane"
)

func removeTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

// newLogger returns a configured slog logger.
func newLogger(w io.Writer) *slog.Logger {
	opts := &humane.Options{
		ReplaceAttr: removeTime,
	}
	logger := slog.New(humane.NewHandler(w, opts))
	return logger.With(slog.String("program", appName))
}
