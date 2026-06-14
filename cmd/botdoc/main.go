// Command botdoc fetches the published Telegram Bot API documentation and
// inspects the structured types and methods extracted from it. It is a
// developer tool: a reference oracle while hand-writing the library and a way
// to spot API drift.
//
// Examples:
//
//	# Summary (version, counts) from the live docs:
//	go run ./cmd/botdoc
//
//	# List every method name:
//	go run ./cmd/botdoc -methods
//
//	# Show the fields and return type of one definition:
//	go run ./cmd/botdoc -name SendMessage
//	go run ./cmd/botdoc -name Message
//
//	# Dump the whole extracted API as JSON (e.g. to diff across versions):
//	go run ./cmd/botdoc -json > api.json
//
//	# Work offline from a saved page:
//	go run ./cmd/botdoc -file api.html -types
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/gotd/botapi/internal/botdoc"
)

const defaultURL = "https://core.telegram.org/bots/api"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "botdoc:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		url     = flag.String("url", defaultURL, "Bot API documentation URL to fetch")
		file    = flag.String("file", "", "read the documentation HTML from a local file instead of fetching")
		asJSON  = flag.Bool("json", false, "print the full extracted API as JSON")
		types   = flag.Bool("types", false, "list all type names")
		methods = flag.Bool("methods", false, "list all method names")
		name    = flag.String("name", "", "show details for a specific type or method")
	)

	flag.Parse()

	doc, err := loadDocument(*url, *file)
	if err != nil {
		return err
	}

	api := botdoc.Extract(doc)

	switch {
	case *asJSON:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		return enc.Encode(api)
	case *name != "":
		return printDefinition(api, *name)
	case *types:
		printNames(api.Types)

		return nil
	case *methods:
		printNames(api.Methods)

		return nil
	default:
		printSummary(api)

		return nil
	}
}

// loadDocument reads the documentation HTML from a file or fetches it.
func loadDocument(url, file string) (*goquery.Document, error) {
	if file != "" {
		f, err := os.Open(file) //nolint:gosec // a CLI reading a user-specified path is the intended behavior
		if err != nil {
			return nil, err
		}

		defer func() { _ = f.Close() }()

		return goquery.NewDocumentFromReader(f)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch %s: unexpected status %s", url, resp.Status)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func printSummary(api botdoc.API) {
	fmt.Printf("Bot API version: %s\n", api.Version)
	fmt.Printf("Types:   %d\n", len(api.Types))
	fmt.Printf("Methods: %d\n", len(api.Methods))
}

func printNames(defs []botdoc.Definition) {
	names := make([]string, len(defs))
	for i, d := range defs {
		names[i] = d.Name
	}

	sort.Strings(names)

	for _, n := range names {
		fmt.Println(n)
	}
}

func printDefinition(api botdoc.API, name string) error {
	if d, ok := findDefinition(api.Types, name); ok {
		printDefinitionDetail("type", d)

		return nil
	}

	if d, ok := findDefinition(api.Methods, name); ok {
		printDefinitionDetail("method", d)

		return nil
	}

	return fmt.Errorf("no type or method named %q", name)
}

// findDefinition matches a definition by name, case-insensitively.
func findDefinition(defs []botdoc.Definition, name string) (botdoc.Definition, bool) {
	for _, d := range defs {
		if strings.EqualFold(d.Name, name) {
			return d, true
		}
	}

	return botdoc.Definition{}, false
}

func printDefinitionDetail(kind string, d botdoc.Definition) {
	fmt.Printf("%s %s\n", kind, d.Name)

	if d.PrettyDescription != "" {
		fmt.Printf("\n%s\n", d.PrettyDescription)
	}

	if d.Ret != nil {
		fmt.Printf("\nreturns: %s\n", d.Ret)
	}

	if len(d.Fields) == 0 {
		return
	}

	label := "Fields"
	if kind == "method" {
		label = "Parameters"
	}

	fmt.Printf("\n%s:\n", label)

	for _, f := range d.Fields {
		opt := ""
		if f.Optional {
			opt = " (optional)"
		}

		fmt.Printf("  %-28s %s%s\n", f.Name, f.Type.String(), opt)

		if len(f.Enum) > 0 {
			fmt.Printf("  %-28s   one of: %s\n", "", strings.Join(f.Enum, ", "))
		}
	}
}
