package cli

import (
	"fmt"
	"log"
	"os"
)

var (
	inform = log.New(os.Stdout, fmt.Sprintf("%6s | ", "Info"), 0).Println
	advise = log.New(os.Stderr, fmt.Sprintf("%6s | ", "Error"), 0).Println
)
