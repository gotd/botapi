// Package botdoc implement types definition extraction from documentation.
package botdoc

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

type API struct {
}

func Extract(doc *goquery.Document) (*API, error) {
	return nil, errors.New("not implemented")
}
