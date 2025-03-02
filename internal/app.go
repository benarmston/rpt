package rpt

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
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

func NewApp(version Version) *cli.App {
	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Usage:              "print the version",
		DisableDefaultText: true,
	}
	cli.VersionPrinter = func(ctx *cli.Context) {
		if ctx.Bool("verbose") {
			fmt.Printf("version=%s revision=%s date=%s\n", version.Version, version.Commit, version.Date)
		} else {
			fmt.Printf("%s\n", version.Version)
		}
	}
	appHelpTemplate := cli.AppHelpTemplate
	appHelpTemplate = strings.ReplaceAll(appHelpTemplate, "GLOBAL OPTIONS", "OPTIONS")
	appHelpTemplate = strings.ReplaceAll(appHelpTemplate, "global options", "options")

	app := &cli.App{
		CustomAppHelpTemplate: appHelpTemplate,
		Name:                  "rpt",
		Usage:                 "run the given command the given number of times",
		UsageText:             "rpt [OPTIONS] TIMES COMMAND [ARGUMENTS...]",
		Description:           "Run `COMMAND ARGUMENTS` TIMES times.",
		HideHelpCommand:       true,
		Version:               version.Version,
		Suggest:               true,
		Authors:               []*cli.Author{{Name: "Ben Armston", Email: ""}},
		Copyright:             "Copyright 2025 Ben Armston",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "wait `DURATION` between runs",
				Value:   0,
			},
			&cli.BoolFlag{
				Name:  "leading-edge",
				Usage: "if given, any provided delay is between the start of\n\tone command invocation and the start of the next.  If\n\tnot given, any provided delay is between the end of\n\tone command invocation and the start of the next\n\t",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "fail-fast",
				Usage: "if command fails exit immediately with the same exit\n\tcode as command",
				Value: false,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print debugging messages",
				Value:   false,
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.NArg() < 2 {
				return errors.New("insufficient arguments")
			}
			if _, err := strconv.Atoi(ctx.Args().First()); err != nil {
				return fmt.Errorf("%s is not an integer", ctx.Args().First())
			}
			return nil
		},
		Action: runRepeatedly,
	}
	return app
}

func runRepeatedly(ctx *cli.Context) error {
	times, _ := strconv.Atoi(ctx.Args().First())
	verbose := ctx.Bool("verbose")
	failFast := ctx.Bool("fail-fast")
	delay := ctx.Duration("delay")
	leadingEdge := ctx.Bool("leading-edge")
	for i := range times {
		var sleepChan <-chan time.Time
		if leadingEdge {
			if verbose {
				log.Printf("starting timer=%s\n", delay)
			}
			sleepChan = time.After(delay)
		}
		cmd := buildCommand(ctx)
		if verbose {
			log.Printf("Iteration=%d; running cmd=%s\n", i, cmd.String())
		}
		if err := runOnce(cmd); err != nil {
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

func buildCommand(ctx *cli.Context) *exec.Cmd {
	args := make([]string, ctx.Args().Len()-2)
	for i, arg := range ctx.Args().Slice() {
		if i >= 2 {
			args[i-2] = arg
		}
	}
	cmd := exec.Command(ctx.Args().Get(1), args...)
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
