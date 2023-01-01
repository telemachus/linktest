// Package cli organizes and implements a command line program.
package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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

	if app.CPUProfile != "" {
		f, err := os.Create(app.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

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

		links := app.GatherLinks(doc)
		app.TestLinks(links)
	}

	if app.MemProfile != "" {
		f, err := os.Create(app.MemProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	return app.ExitStatus()
}
