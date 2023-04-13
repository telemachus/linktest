package cli

import (
	"io"

	"github.com/telemachus/humane"
	"golang.org/x/exp/slog"
)

func removeTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

// newLogger returns a configured slog logger.
func newLogger(w io.Writer) *slog.Logger {
	ho := humane.Options{
		ReplaceAttr: removeTime,
	}
	logger := slog.New(ho.NewHandler(w))
	return logger.With(slog.String("program", appName))
}
