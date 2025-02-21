//go:build docs
// +build docs

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	app := newApp()

	md, err := app.ToMarkdown()
	if err != nil {
		log.Fatalf("generating markdown: %s", err)
	}
	md = strings.ReplaceAll(md, "GLOBAL OPTIONS", "OPTIONS")
	f, err := os.Create("docs/cli.md")
	if err != nil {
		log.Fatalf("opening output file: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString("# CLI\n\n" + md); err != nil {
		log.Fatalf("writing output: %s", err)
	}
	fmt.Printf("Wrote markdown docs to %s\n", f.Name())

	man, err := app.ToManWithSection(1)
	if err != nil {
		log.Fatalf("generating manpage: %s", err)
	}
	man = strings.ReplaceAll(man, "GLOBAL OPTIONS", "OPTIONS")

	f, err = os.Create(fmt.Sprintf("docs/%s.1", app.Name))
	if err != nil {
		log.Fatalf("opening output file: %s", err)
	}
	defer f.Close()
	if _, err := f.WriteString(man); err != nil {
		log.Fatalf("writing output: %s", err)
	}
	fmt.Printf("Wrote manpage to %s\n", f.Name())
}
