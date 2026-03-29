package rpt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
		UsageText:                     "rpt [OPTIONS] TIMES COMMAND [ARGUMENTS...]",
		Description:                   "Run `COMMAND ARGUMENTS` TIMES times.",
		HideHelpCommand:               true,
		Version:                       version.Version,
		Suggest:                       true,
		Authors:                       []any{"Ben Armston"},
		Copyright:                     "Copyright 2025 Ben Armston. Licensed under the MIT License.",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "wait `DURATION` between runs",
				Value:   0,
			},
			&cli.BoolFlag{
				Name:  "leading-edge",
				Usage: "if given, any provided delay is between the\nstart of one command invocation and the start the next. If not given,\nany provided delay is between the end of one invocation and the start of\nthe next",
				Value: false,
			},
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
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.NArg() < 2 {
				return ctx, errors.New("insufficient arguments")
			}
			if _, err := strconv.Atoi(cmd.Args().First()); err != nil {
				return ctx, fmt.Errorf("%s is not an integer", cmd.Args().First())
			}
			return ctx, nil
		},
		Action: runRepeatedly,
	}
	return app
}

func runRepeatedly(ctx context.Context, cmd *cli.Command) error {
	times, _ := strconv.Atoi(cmd.Args().First())
	verbose := cmd.Bool("verbose")
	failFast := cmd.Bool("fail-fast")
	delay := cmd.Duration("delay")
	leadingEdge := cmd.Bool("leading-edge")
	for i := range times {
		var sleepChan <-chan time.Time
		if leadingEdge {
			if verbose {
				log.Printf("starting timer=%s\n", delay)
			}
			sleepChan = time.After(delay)
		}
		exe := buildExecCommand(cmd)
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

func buildExecCommand(c *cli.Command) *exec.Cmd {
	args := make([]string, c.Args().Len()-2)
	for i, arg := range c.Args().Slice() {
		if i >= 2 {
			args[i-2] = arg
		}
	}
	cmd := exec.Command(c.Args().Get(1), args...)
	return cmd
}

func runOnce(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("creating stderr pipe: %w", err)
	}
	go func() {
		io.Copy(os.Stdout, stdout) //nolint:errcheck
	}()
	go func() {
		io.Copy(os.Stderr, stderr) //nolint:errcheck
	}()
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}
	return cmd.Wait()
}
