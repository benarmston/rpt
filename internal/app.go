package rpt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/muesli/reflow/wordwrap"
	"github.com/urfave/cli/v3"
)

type Version struct {
	Version string
	Commit  string
	Date    string
}

var DefaultVersion = Version{
	Version: "dev",
	Commit:  "unknown",
	Date:    "unknown",
}

func NewApp(version Version) *cli.Command {
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "print the version",
	}
	cli.VersionPrinter = func(cmd *cli.Command) {
		if cmd.Bool("verbose") {
			fmt.Printf("version=%s revision=%s date=%s\n", version.Version, version.Commit, version.Date)
		} else {
			fmt.Printf("%s\n", version.Version)
		}
	}
	helpTemplate := cli.RootCommandHelpTemplate
	helpTemplate = strings.ReplaceAll(helpTemplate, "GLOBAL OPTIONS", "OPTIONS")
	helpTemplate = strings.ReplaceAll(helpTemplate, "global options", "options")

	app := &cli.Command{
		CustomRootCommandHelpTemplate: helpTemplate,
		Name:                          "rpt",
		Usage:                         "run the given command the given number of times",
		UsageText:                     "rpt [OPTIONS] TIMES COMMAND [-- ARGUMENTS...]",
		Description:                   wordwrap.String("Run `COMMAND ARGUMENTS` TIMES times.\n\nIf the '--delay' option is given, there will be a delay of the given DURATION after one run ends and the next starts. This provides a guaranteed delay between runs.\n\nIf the '--every' option is given, COMMAND will be run every DURATION. If COMMAND takes longer to run than the given DURATION the next run will start immediately once the current run has completed. This provides a predictable start time for each run (provided COMMAND consistently completes in under DURATION).", 74),
		HideHelpCommand:               true,
		Version:                       version.Version,
		Suggest:                       true,
		Authors:                       []any{"Ben Armston"},
		Copyright:                     "Copyright 2025 Ben Armston. Licensed under the MIT License.",
		MutuallyExclusiveFlags: []cli.MutuallyExclusiveFlags{
			{
				Flags: [][]cli.Flag{
					{
						&cli.DurationFlag{
							Name:    "delay",
							Aliases: []string{"d"},
							Usage:   "wait `DURATION` between one run ending and\nthe next starting",
							Value:   0,
						},
					},
					{
						&cli.DurationFlag{
							Name:    "every",
							Aliases: []string{"e"},
							Usage:   "run COMMAND every `DURATION`. The next run\nwill start DURATION after the previous run started or as soon as the\nprevious run ends if it takes longer than DURATION",
							Value:   0,
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "fail-fast",
				Usage: "if COMMAND fails, exit immediately with the same exit code\nas COMMAND",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print debugging messages",
			},
		},
		Arguments: []cli.Argument{
			&cli.IntArg{Name: "times", UsageText: "TIMES"},
			&cli.StringArg{Name: "command", UsageText: "COMMAND"},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.NArg() < 2 {
				return ctx, errors.New("insufficient arguments")
			}
			times, err := strconv.Atoi(cmd.Args().First())
			if err != nil {
				return ctx, fmt.Errorf("TIMES must be an integer; got %s", cmd.Args().First())
			}
			if times < 1 {
				return ctx, fmt.Errorf("TIMES must be at least 1; got %s", cmd.Args().First())
			}
			return ctx, nil
		},
		Action: runRepeatedly,
	}
	return app
}

func runRepeatedly(ctx context.Context, cmd *cli.Command) error {
	times := cmd.IntArg("times")
	verbose := cmd.Bool("verbose")
	failFast := cmd.Bool("fail-fast")

	var delay time.Duration
	var leadingEdge bool
	if slices.Contains(cmd.FlagNames(), "delay") {
		delay = cmd.Duration("delay")
		leadingEdge = false
	} else if slices.Contains(cmd.FlagNames(), "every") {
		delay = cmd.Duration("every")
		leadingEdge = true
	}

	for i := range times {
		var sleepChan <-chan time.Time
		if leadingEdge {
			if verbose {
				log.Printf("starting timer=%s\n", delay)
			}
			sleepChan = time.After(delay)
		}
		exe := exec.Command(cmd.StringArg("command"), cmd.Args().Slice()...)
		if verbose {
			log.Printf("Iteration=%d; running cmd=%s\n", i, exe.String())
		}
		if err := runOnce(exe); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if failFast && exitErr.ExitCode() != 0 {
					break
				}
			} else {
				log.Printf("%s\n", err)
			}
		}
		if i != times-1 {
			if verbose && leadingEdge {
				log.Printf("waiting for timer")
			} else if verbose {
				log.Printf("sleeping=%s\n", delay)
			}
			if !leadingEdge {
				sleepChan = time.After(delay)
			}
			<-sleepChan
		}
	}
	return nil

}

func runOnce(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}
	return cmd.Wait()
}
