package cli_test

import (
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/linktest/cli"
)

const goldenFile = "testdata/telemachus.html"

func wantedLinks() []string {
	return []string{
		"https://telemachus.me/foobar",
		"https://telemachus.me/about",
		"https://git.sr.ht/~telemachus",
		"https://telemachus.me/aboot",
		"https://telemachus.me/git-can-break-your-tests",
		"https://telemachus.me/site-setup",
		"https://telemachus.me/reading-february-2021",
		"http://telemachus.me/reading-january-2021",
		"https://telemachus.me/python-descartes-descartes-python",
		"https://telemachus.me/atom.xml",
	}
}

func TestGetLinks(t *testing.T) {
	t.Parallel()

	fh, err := os.Open(goldenFile)
	if err != nil {
		t.Fatalf("cannot open %s: %v", goldenFile, err)
	}

	doc, err := io.ReadAll(fh)
	if err != nil {
		t.Fatalf("cannot read %s: %v", goldenFile, err)
	}

	wanted := wantedLinks()
	app := &cli.App{ExitValue: 0}
	got := app.GetLinks(doc)

	if !cmp.Equal(wanted, got) {
		t.Error(cmp.Diff(wanted, got))
	}
}
