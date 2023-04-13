package cli

import (
	"bytes"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// Link stores information about a link.
type Link struct {
	File   string
	URL    string
	Status int
	Err    error
}

func (app *App) LinkGen(upstream chan *File) chan *Link {
	if app.NoOp() {
		return nil
	}
	var wg sync.WaitGroup
	downstream := make(chan *Link)
	for file := range upstream {
		wg.Add(1)
		go func(f *File, c chan *Link) {
			scanLink(f, c)
			wg.Done()
		}(file, downstream)
	}
	go func() {
		wg.Wait()
		close(downstream)
	}()
	return downstream
}

func scanLink(file *File, downstream chan *Link) {
	tkn := html.NewTokenizer(bytes.NewReader(file.Content))
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := tkn.Token()
			if t.Data != "a" {
				continue
			}
			link, ok := getLink(t)
			if !ok {
				continue
			}
			downstream <- &Link{
				File: file.Name,
				URL:  link,
			}
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

func (app *App) CheckStatus(linkCh <-chan *Link) <-chan *Link {
	if app.NoOp() {
		return nil
	}
	var wg sync.WaitGroup
	downstream := make(chan *Link)
	for li := range linkCh {
		wg.Add(1)
		go func(l *Link) {
			testLink(l, downstream)
			wg.Done()
		}(li)
	}
	go func() {
		wg.Wait()
		close(downstream)
	}()
	return downstream
}

func testLink(link *Link, ch chan<- *Link) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(link.URL)
	if err != nil {
		link.Err = err
		ch <- link
		return
	}
	defer resp.Body.Close()
	link.Status = resp.StatusCode
	ch <- link
}
