// Package cli organizes and implements a command line program.
package cli

import (
	"fmt"
	"io"
	"os"
)

const (
	exitSuccess = 0
	exitFailure = 1
	appName     = "linktest"
	appVersion  = "v0.0.1"
	appUsage    = `usage: linktest [-verbose] file ...

options:
    -verbose     Show responses from all links rather than only non-200 responses
    -help, -h    Show this message
    -version     Show version`
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}
	files := app.ParseFlags(args)

	for _, file := range files {
		fh, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: skipping %q: %v\n", appName, file, err)
			app.FileProblems++

			continue
		}

		doc, err := io.ReadAll(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: skipping %q: %v\n", appName, file, err)
			app.FileProblems++

			continue
		}

		links := app.GetLinks(doc)
		app.TestLinks(links)
	}

	return app.ExitStatus()
}
