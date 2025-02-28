//go:build docs
// +build docs

// Based on code taken from https://github.com/urfave/cli/ Copyright (c) 2023
// urfave/cli maintainers and licensed under the MIT License.

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/benarmston/rpt/internal"
	"github.com/cpuguy83/go-md2man/v2/md2man"
	"github.com/urfave/cli/v2"
)

func main() {
	app := rpt.NewApp(rpt.DefaultVersion)
	license := "Licensed under the MIT License."
	file := "docs/usage.md"
	writeMarkdown(app, license, file)
	fmt.Printf("Wrote markdown docs to %s\n", file)

	file = "docs/rpt.1"
	writeMan(app, license, file)
	fmt.Printf("Wrote markdown docs to %s\n", file)
}

func writeMarkdown(app *cli.App, license, path string) {
	md, err := appToMarkdown(app, license)
	if err != nil {
		log.Fatalf("generating markdown: %s", err)
	}
	f, err := os.Create("docs/usage.md")
	if err != nil {
		log.Fatalf("opening output file: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString(md); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func writeMan(app *cli.App, license, path string) {
	man, err := appToMan(app, license)
	if err != nil {
		log.Fatalf("generating manpage: %s", err)
	}
	f, err := os.Create("docs/rpt.1")
	if err != nil {
		log.Fatalf("opening output file: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString(man); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

//go:embed usage.md.go.tmpl
var MarkdownDocTemplate string

// ToMarkdown creates a markdown string for the `*App`
// The function errors if either parsing or writing of the string fails.
func appToMarkdown(a *cli.App, license string) (string, error) {
	var w bytes.Buffer
	if err := writeDocTemplate(a, &w, license, 0); err != nil {
		return "", err
	}
	return w.String(), nil
}

// ToMarkdown creates a markdown string for the `*App`
// The function errors if either parsing or writing of the string fails.
func appToMan(a *cli.App, license string) (string, error) {
	var w bytes.Buffer
	if err := writeDocTemplate(a, &w, license, 1); err != nil {
		return "", err
	}
	man := md2man.Render(w.Bytes())
	return string(man), nil
}

type cliTemplate struct {
	App        *cli.App
	SectionNum int
	Commands   []string
	GlobalArgs []string
	License    string
}

func writeDocTemplate(a *cli.App, w io.Writer, license string, section int) error {
	const name = "cli"
	t, err := template.New(name).Parse(MarkdownDocTemplate)
	if err != nil {
		return err
	}
	return t.ExecuteTemplate(w, name, &cliTemplate{
		App:        a,
		SectionNum: section,
		Commands:   prepareCommands(a.Commands, 0),
		GlobalArgs: prepareArgsWithValues(a.VisibleFlags()),
		License:    license,
	})
}

func prepareCommands(commands []*cli.Command, level int) []string {
	var coms []string
	for _, command := range commands {
		if command.Hidden {
			continue
		}

		usageText := prepareUsageText(command)

		usage := prepareUsage(command, usageText)

		prepared := fmt.Sprintf("%s %s\n\n%s%s",
			strings.Repeat("#", level+2),
			strings.Join(command.Names(), ", "),
			usage,
			usageText,
		)

		flags := prepareArgsWithValues(command.VisibleFlags())
		if len(flags) > 0 {
			prepared += fmt.Sprintf("\n%s", strings.Join(flags, "\n"))
		}

		coms = append(coms, prepared)

		// recursively iterate subcommands
		if len(command.Subcommands) > 0 {
			coms = append(
				coms,
				prepareCommands(command.Subcommands, level+1)...,
			)
		}
	}

	return coms
}

func prepareArgsWithValues(flags []cli.Flag) []string {
	return prepareFlags(flags, ", ", "**", "**", true)
}

// Returns the placeholder, if any, and the unquoted usage string.
func unquoteUsage(usage string) (string, string) {
	for i := range len(usage) {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name := usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break
		}
	}
	return "", usage
}

func prepareFlags(
	flags []cli.Flag,
	sep, opener, closer string,
	addDetails bool,
) []string {
	args := []string{}
	for _, f := range flags {
		flag, ok := f.(cli.DocGenerationFlag)
		if !ok {
			continue
		}
		modifiedArg := opener

		names := flag.Names()
		slices.Reverse(names)
		for _, s := range names {
			trimmed := strings.TrimSpace(s)
			if len(modifiedArg) > len(opener) {
				modifiedArg += sep
			}
			if len(trimmed) > 1 {
				modifiedArg += fmt.Sprintf("--%s", trimmed)
			} else {
				modifiedArg += fmt.Sprintf("-%s", trimmed)
			}
			if flag.TakesValue() {
				placeholder, _ := unquoteUsage(flag.GetUsage())
				if placeholder == "" {
					placeholder = "value"
				}
				modifiedArg += fmt.Sprintf("=%s", placeholder)
			}
		}
		modifiedArg += closer

		if addDetails {
			modifiedArg += flagDetails(flag)
		}

		args = append(args, modifiedArg+"\n")

	}
	sort.Strings(args)
	return args
}

// flagDetails returns a string containing the flags metadata
func flagDetails(flag cli.DocGenerationFlag) string {
	description := flag.GetUsage()
	if flag.TakesValue() {
		defaultText := flag.GetDefaultText()
		if defaultText == "" {
			defaultText = flag.GetValue()
		}
		if defaultText != "" {
			description += " (default: " + defaultText + ")"
		}
	}
	return "\n: " + description
}

func prepareUsageText(command *cli.Command) string {
	if command.UsageText == "" {
		return ""
	}

	// Remove leading and trailing newlines
	preparedUsageText := strings.Trim(command.UsageText, "\n")

	var usageText string
	if strings.Contains(preparedUsageText, "\n") {
		// Format multi-line string as a code block using the 4 space schema to allow for embedded markdown such
		// that it will not break the continuous code block.
		for _, ln := range strings.Split(preparedUsageText, "\n") {
			usageText += fmt.Sprintf("    %s\n", ln)
		}
	} else {
		// Style a single line as a note
		usageText = fmt.Sprintf(">%s\n", preparedUsageText)
	}

	return usageText
}

func prepareUsage(command *cli.Command, usageText string) string {
	if command.Usage == "" {
		return ""
	}

	usage := command.Usage + "\n"
	// Add a newline to the Usage IFF there is a UsageText
	if usageText != "" {
		usage += "\n"
	}

	return usage
}
