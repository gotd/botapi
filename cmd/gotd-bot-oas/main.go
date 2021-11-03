package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/ogen-go/errors"

	"github.com/gotd/botapi/botdoc"
)

func run(ctx context.Context) error {
	var arg struct {
		URL    string
		Target string
	}
	flag.StringVar(&arg.URL, "url", "https://core.telegram.org/bots/api", "bot url")
	flag.StringVar(&arg.Target, "target", filepath.Join("_oas", "openapi.json"), "output file")
	flag.Parse()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, arg.URL, nil)
	if err != nil {
		return errors.Wrap(err, "req")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do")
	}

	defer func() { _ = res.Body.Close() }()
	switch res.StatusCode {
	case http.StatusOK: // ok
	default:
		return errors.Errorf("code %d: %s", res.StatusCode, res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "read")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	api := botdoc.Extract(doc)
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)
	e.SetIndent("", "  ")
	s, err := api.OAS()
	if err != nil {
		return errors.Wrap(err, "generate")
	}
	if err := e.Encode(s); err != nil {
		return errors.Wrap(err, "encode")
	}

	if err := os.WriteFile(arg.Target, buf.Bytes(), 0600); err != nil {
		return errors.Wrap(err, "write")
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %+v\n", err)
		os.Exit(2)
	}
}
