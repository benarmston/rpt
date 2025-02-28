package rpt

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
		Usage:                 "repeat running a command a number of times",
		UsageText:             "rpt [OPTIONS] COMMAND [ARGUMENTS...]",
		Description:           "Repeatedly run COMMAND with ARGUMENTS.  The number of times to run COMMAND\nis determined by OPTIONS.",
		HideHelpCommand:       true,
		Version:               version.Version,
		Suggest:               true,
		Authors:               []*cli.Author{{Name: "Ben Armston", Email: ""}},
		Copyright:             "Copyright 2025 Ben Armston",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    "times",
				Aliases: []string{"t"},
				Usage:   "number of `TIMES` to run COMMAND",
				Value:   1,
			},
			&cli.DurationFlag{
				Name:    "delay",
				Aliases: []string{"d"},
				Usage:   "wait `DURATION` between runs",
				Value:   0,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "print debugging messages",
				Value:   false,
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.NArg() < 1 {
				return errors.New("missing command")
			}
			return nil
		},
		Action: runRepeatedly,
	}
	return app
}

func runRepeatedly(ctx *cli.Context) error {
	verbose := ctx.Bool("verbose")
	times := ctx.Int64("times")
	delay := ctx.Duration("delay")
	for i := range times {
		if i != 0 {
			if verbose {
				log.Printf("sleeping=%s\n", delay)
			}
			time.Sleep(delay)
		}
		cmd := exec.Command(ctx.Args().First(), ctx.Args().Tail()...)
		if verbose {
			log.Printf("Iteration=%d; running cmd=%s\n", i, cmd.String())
		}
		if err := runOnce(cmd); err != nil {
			log.Printf("%s\n", err)
		}
	}
	return nil

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
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("waiting on command: %w", err)
	}
	return nil
}
