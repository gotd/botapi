package botdoc

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

type Field struct {
	// Can be one of:
	// 	* Array of T
	// 	* Array of Array of T
	//	* T1 or T2 or T3
	// Basic types:
	//	* String
	//	* Integer
	//	* Float[ number]
	Type        string
	Name        string
	Description string
	Optional    bool
}

type Definition struct {
	Name        string
	Description string
	Fields      []Field

	Method bool
}

func TestExtract(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("_testdata", "api.html"))
	require.NoError(t, err)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	require.NoError(t, err)

	var (
		def Definition
		d   []Definition
	)

	doc.Find("#dev_page_content").Children().Each(func(i int, s *goquery.Selection) {
		if s.Is("h4") {
			def.Name = s.Text()
		}
		if s.Is("table") {
			var head []string
			s.Find("th").Each(func(i int, s *goquery.Selection) {
				head = append(head, s.Text())
			})
			if len(head) == 0 {
				return
			}
			switch head[0] {
			case "Field":
				// Type
			case "Parameter":
				// Method
			default:
				return
			}
			s.Find("tr").Each(func(i int, s *goquery.Selection) {
				var definition []string
				s.Find("td").Each(func(j int, s *goquery.Selection) {
					definition = append(definition, s.Text())
				})
				if len(definition) == 0 {
					return
				}
				if len(definition) == 3 {
					def.Fields = append(def.Fields, Field{
						Name:        definition[0],
						Description: definition[2],
						Optional:    strings.HasPrefix(definition[2], "Optional"),
						Type:        definition[1],
					})
				} else if len(definition) == 4 {
					def.Fields = append(def.Fields, Field{
						Name:        definition[0],
						Description: definition[3],
						Optional:    strings.HasPrefix(definition[2], "Optional"),
						Type:        definition[1],
					})
				}
			})

			d = append(d, def)
			def = Definition{}
		}
	})

	for _, dd := range d {
		t.Log(dd.Name)
		for _, f := range dd.Fields {
			t.Logf(" %s %s", f.Name, f.Type)
		}
	}
}
