package cli_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/telemachus/linktest/cli"
)

const (
	goldenFile = "testdata/telemachus.html"
	emptyFile  = "testdata/empty.html"
)

func mixedLinks() []string {
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

func TestMixedLinks(t *testing.T) {
	t.Parallel()

	wanted := mixedLinks()
	app := &cli.App{ExitValue: 0}
	got := app.GetLinks(goldenFile)

	if !cmp.Equal(wanted, got) {
		t.Error(cmp.Diff(wanted, got))
	}
}

func TestEmptyFile(t *testing.T) {
	t.Parallel()

	wanted := []string{}
	app := &cli.App{ExitValue: 0}
	got := app.GetLinks(emptyFile)

	if !cmp.Equal(wanted, got) {
		t.Error(cmp.Diff(wanted, got))
	}
}
