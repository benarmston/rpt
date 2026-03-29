//go:build docs
// +build docs

// Based on code taken from https://github.com/urfave/cli/ Copyright (c) 2023
// urfave/cli maintainers and licensed under the MIT License.

package main

import (
	_ "embed"
	"fmt"
	"os"

	rpt "github.com/benarmston/rpt/internal"
	docs "github.com/urfave/cli-docs/v3"
)

//go:embed usage.md.go.tmpl
var MarkdownDocTemplate string

func main() {
	app := rpt.NewApp(rpt.DefaultVersion)
	docs.MarkdownDocTemplate = MarkdownDocTemplate
	md, err := docs.ToMarkdown(app)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("docs/usage.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString(md); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote markdown docs to %s\n", file.Name())

	man, err := docs.ToManWithSection(app, 1)
	if err != nil {
		panic(err)
	}
	file, err = os.Create("docs/rpt.1")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString(man); err != nil {
		panic(err)
	}
	fmt.Printf("Wrote manpage to %s\n", file.Name())
}
