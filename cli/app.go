package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// App stores information about the application's state.
type App struct {
	ExitValue     int
	FileProblems  int
	HelpWanted    bool
	LinkProblems  int
	Verbose       bool
	VersionWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}

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

// GatherLinks parses one or more files as HTML and returns links to test.
func (app *App) GatherLinks(text []byte) []string {
	if app.NoOp() {
		return []string{}
	}

	tkn := html.NewTokenizer(bytes.NewReader(text))
	links := []string{}

	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := tkn.Token()

			if t.Data != "a" {
				continue
			}

			link, ok := getLink(t)
			if !ok {
				continue
			}

			links = append(links, link)
		}
	}
}

func getLink(t html.Token) (string, bool) {
	link := ""
	ok := false

	for _, a := range t.Attr {
		if a.Key == "href" && strings.HasPrefix(a.Val, "http") {
			link = a.Val
			ok = true

			break
		}
	}

	return link, ok
}

// TestLinks runs GET requests on links to test for link rot.
func (app *App) TestLinks(links []string) {
	if app.NoOp() {
		return
	}

	ch := make(chan string)

	for _, link := range links {
		go app.testLink(link, ch)
	}

	for range links {
		msg := <-ch
		if app.Verbose || !strings.HasSuffix(msg, "200 OK") {
			fmt.Println(msg)
		}
	}
}

func (app *App) testLink(link string, ch chan<- string) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(link)
	if err != nil {
		ch <- fmt.Sprintf("%s: %s: %v", appName, link, err)
		app.LinkProblems++

		return
	}

	defer resp.Body.Close()
	ch <- fmt.Sprintf("%s: %q: %d %s", appName, link, resp.StatusCode,
		http.StatusText(resp.StatusCode))
}

// ExitStatus returns an appropriate exit status for an app.
func (app *App) ExitStatus() int {
	if app.ExitValue != exitSuccess || app.FileProblems > 0 || app.LinkProblems > 0 {
		return exitFailure
	}

	return exitSuccess
}
