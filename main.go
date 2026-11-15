package main

import (
	"context"
	"errors"
	"log"
	"os"

	rpt "github.com/benarmston/rpt/internal"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	// Don't include a timestamp in log output.
	log.SetFlags(0)
}

func main() {
	version := rpt.Version{Version: version, Commit: commit, Date: date}
	cmd := rpt.NewApp(version)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if exitError, ok := errors.AsType[rpt.SilentExitError](err); ok {
			os.Exit(exitError.ExitCode)
		}
		log.Fatalf(
			"%s\n\nUsage:  %s\nSee  %s --help for help.",
			err, cmd.UsageText, cmd.Name,
		)
	}
}
