package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// ParseFlags handles flags and options in my finicky way.
func (app *App) ParseFlags(args []string) []string {
	flags := flag.NewFlagSet(appName, flag.ContinueOnError)
	// Contrary to Go's defaults, I want usage to go to stdout if the user
	// explicitly asks for help. Therefore, I need to handle the `-help` flag
	// manually.
	// See https://github.com/golang/go/issues/41523 for discussion.
	flags.SetOutput(io.Discard)
	// The final argument to these functions contains the flag's usage string.
	// However, I define a custom usage message, so I don't need to define
	// usage here.
	// flag treats "-h" like "-help" by default, so I need to catch that too.
	flags.BoolVar(&app.HelpWanted, "h", false, "")
	flags.BoolVar(&app.HelpWanted, "help", false, "")
	flags.BoolVar(&app.Verbose, "verbose", false, "")
	flags.BoolVar(&app.VersionWanted, "version", false, "")
	err := flags.Parse(args)
	files := flags.Args()
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "%s: %s\n%s\n", appName, err, appUsage)
		app.ExitValue = exitFailure
	case app.HelpWanted:
		fmt.Println(appUsage)
	case app.VersionWanted:
		fmt.Printf("%s: %s\n", appName, appVersion)
	case len(files) < 1:
		fmt.Fprintf(os.Stderr, "%s: specify one or more HTML files to test\n", appName)
		app.ExitValue = exitFailure
	}
	return files
}
