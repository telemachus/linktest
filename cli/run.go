// Package cli organizes and implements a command line program.
package cli

import "os"

const (
	exitSuccess = 0
	exitFailure = 1
	appName     = "linktest"
	appVersion  = "v0.1.0"
	appUsage    = `usage: linktest [-verbose] file ...

options:
    -verbose     Show responses from all links rather than only non-200 responses
    -help, -h    Show this message
    -version     Show version`
)

// Run creates an App, performs the App's tasks, and returns an exit value.
func Run(args []string) int {
	app := &App{ExitValue: exitSuccess}
	fileNames := app.ParseFlags(args)
	logger := newLogger(os.Stderr)
	fileCh := app.FileGen(fileNames, logger)
	linkCh := app.LinkGen(fileCh)
	statusCh := app.CheckStatus(linkCh)
	app.Display(statusCh, logger)
	return app.ExitValue
}
