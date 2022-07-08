package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-faster/errors"
	"github.com/go-json-experiment/json"

	"github.com/gotd/botapi/botdoc"
)

func marshalDeterministic(src interface{}) ([]byte, error) {
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	val := json.RawValue(data)
	if err := val.Canonicalize(); err != nil {
		return nil, errors.Wrap(err, "canonicalize")
	}
	if err := val.Indent("", "\t"); err != nil {
		return nil, errors.Wrap(err, "indent")
	}
	return val, nil
}

func run(ctx context.Context) error {
	var arg struct {
		URL    string
		Target string
	}
	flag.StringVar(&arg.URL, "url", "https://core.telegram.org/bots/api", "bot url")
	flag.StringVar(&arg.Target, "target", filepath.Join("_oas", "openapi.json"), "output file")
	flag.Parse()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, arg.URL, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "req")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do")
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("code %d: %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return errors.Wrap(err, "parse")
	}

	spec, err := botdoc.Extract(doc).OAS()
	if err != nil {
		return errors.Wrap(err, "generate")
	}

	data, err := marshalDeterministic(spec)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	return os.WriteFile(arg.Target, data, 0o600)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed: %+v\n", err)
		os.Exit(2)
	}
}
