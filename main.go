// Package main runs the command line and exits with an appropriate status.
package main

import (
	"os"

	"github.com/telemachus/linktest/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
