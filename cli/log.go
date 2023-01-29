package cli

import (
	"io"
	stdlog "log"

	kitlog "github.com/go-kit/log"
)

// newLogger returns a configured go-kit logger.
func newLogger(w io.Writer) kitlog.Logger {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(w))
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
	logger = kitlog.With(logger,
		"program", appName,
	)

	return logger
}
