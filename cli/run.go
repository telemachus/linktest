// Package cli organizes and implements a command line program.
package cli

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
		links := app.GetLinks(file)
		app.TestLinks(links)
	}

	return app.ExitStatus()
}
