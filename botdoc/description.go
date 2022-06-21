package botdoc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var typosReplacer = strings.NewReplacer(
	`unpriviledged`, `unprivileged`,
	`Url`, `URL`,
	"»", "",
	// Replace Unicode quotes.
	"“", `"`,
	"”", `"`,
	// Replace apostrophe to single quote.
	"\u2019", `'`,
	// Replace ellipsis to 3 dots.
	"…", `...`,
	// Replace Unicode dashes to ASCII dash.
	"\u2013", `-`,
	"\u2014", `-`,
)

func fixTypos(s string) string {
	return typosReplacer.Replace(s)
}

func cleanDescription(s string) string {
	return strings.TrimSpace(fixTypos(s))
}

func selDescription(sel *goquery.Selection) string {
	sel = sel.Clone()
	sel.Find("a").Each(func(i int, s *goquery.Selection) {
		ref, ok := s.Attr("href")
		if !ok {
			return
		}

		u, err := url.Parse(ref)
		if err != nil {
			return
		}

		link := rootURL.ResolveReference(u).String()
		text := strings.ReplaceAll(s.Text(), "»", "")
		text = strings.TrimSpace(text)

		s.SetText(fmt.Sprintf("[%s](%s)", text, link))
	})
	return cleanDescription(sel.Text())
}
