//go:build !docs
// +build !docs

package main

import (
	"log"
	"os"
)

func init() {
	// Don't include a timestamp in log output.
	log.SetFlags(0)
}

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatalf(
			"%s\n\nUsage:  %s\nSee  %s --help for help.",
			err, app.UsageText, app.HelpName,
		)
	}
}
