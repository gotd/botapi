package botdoc

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// Extract API definition from goquery document.
func Extract(doc *goquery.Document) (a API) {
	var (
		d   Definition
		sec section
	)
	doc.Find("#dev_page_content").Children().Each(func(i int, s *goquery.Selection) {
		// Replace emoji images with alts.
		s.Find(".emoji").Each(func(i int, s *goquery.Selection) {
			s.ReplaceWithHtml(s.AttrOr("alt", ""))
		})

		if text := strings.TrimPrefix(s.Text(), "Bot API "); s.Is("p") &&
			text != s.Text() &&
			(a.Version == "" || text > a.Version) {
			a.Version = text
		}
		if s.Is("h3") {
			switch strings.TrimSpace(s.Text()) {
			case "Available types":
				sec = sectionTypes
			case "Available methods":
				sec = sectionMethods
			}
		}
		appendDefinition := func() {
			if d.Name == "" {
				d = Definition{}
				return
			}
			for i, c := range d.Name {
				if i != 0 {
					break
				}
				if c == unicode.ToUpper(c) {
					sec = sectionTypes
				} else {
					sec = sectionMethods
				}
			}
			const canBeString = "String can be used instead of this object"
			if strings.Contains(d.Description, canBeString) {
				newName := d.Name + "Object"
				a.Types = append(a.Types, Definition{
					Name:        d.Name,
					Description: d.Description,
					Ret: &Type{
						Name: d.Name,
						Kind: KindSum,
						Sum: []Type{
							newPrimitive(String),
							{Name: newName, Kind: KindObject},
						},
					},
				})
				d.Name = newName
			}
			switch sec {
			case sectionMethods:
				a.Methods = append(a.Methods, d)
			case sectionTypes:
				a.Types = append(a.Types, d)
			}
			d = Definition{}
		}
		if s.Is("h4") {
			d.Name = strings.TrimSpace(s.Text())
			return
		}
		if s.Is("p") && d.Name != "" {
			d.Description = selDescription(s)
			if strings.Contains(strings.ToLower(d.Description), `currently holds no information`) {
				appendDefinition()
			}
			if strings.Contains(d.Description, `Returns basic information about the bot`) {
				d.Ret = &Type{
					Kind: KindObject,
					Name: "User",
				}
				appendDefinition()
			}
		}
		switch desc := d.Description; {
		case strings.Contains(desc, `as String on success`):
			t := newPrimitive(String)
			d.Ret = &t
		case strings.Contains(desc, `Returns True on success`):
			t := newPrimitive(Boolean)
			d.Ret = &t
		case strings.Contains(desc, `Returns Int on success`):
			t := newPrimitive(Integer)
			d.Ret = &t
		}

		// Detect definition of sum-types.
		var (
			probablySum    bool
			probablyMarker string
		)
		for _, sumMarker := range []string{
			`It should be one of`,
			`Telegram clients currently support the following`,
			`Currently, the following`,
			`This object represents one result of`,
		} {
			if strings.Contains(d.Description, sumMarker) {
				probablySum = true
				probablyMarker = sumMarker
			}
		}
		if s.Is("ul") && probablySum {
			t := &Type{
				Kind: KindSum,
			}
			d.Description = strings.TrimSpace(strings.ReplaceAll(d.Description, probablyMarker, ""))
			s.Find("li").Each(func(i int, s *goquery.Selection) {
				t.Sum = append(t.Sum, ParseType(s.Text()))
			})
			d.Ret = t
			appendDefinition()
			return
		}

		if d.Ret == nil {
			var links []string
			s.Find("a").Each(func(i int, selection *goquery.Selection) {
				if href, _ := selection.Attr("href"); strings.HasPrefix(href, "#") {
					links = append(links, selection.Text())
				}
			})
			const (
				retPrefix       = `on success, the`
				retPrefix2      = `on success, a`
				retPrefix3      = `returns a`
				retPrefix4      = `returns the`
				retArrayPrefix  = `an array of`
				retArrayPrefix2 = `returns array of`
				retSuffix       = ` is returned`
				retSuffix2      = ` object`
				retSuffix3      = ` objects`
				retSuffix4      = ` on success`
			)
			var (
				start, end int
				prefix     string
			)
			loweredDesc := strings.TrimSuffix(strings.ToLower(d.Description), ".")
			start, prefix = IndexOneOf(loweredDesc,
				retArrayPrefix,
				retArrayPrefix2,
				retPrefix,
				retPrefix2,
				retPrefix3,
				retPrefix4,
			)
			if prefix == retArrayPrefix || prefix == retArrayPrefix2 {
				// Do not cut prefix, if we do ParseType will be unable to detect an array clause.
				prefix = ""
			}

			end, _ = IndexOneOf(loweredDesc, retSuffix, retSuffix2, retSuffix3, retSuffix4)
			if start > 0 && end > start {
				ret := strings.TrimSpace(d.Description[start+len(prefix) : end])
				ret = strings.TrimSuffix(ret, ".")
				ret = strings.TrimSuffix(ret, "object")
				ret = strings.TrimSuffix(ret, "objects")

				var found bool
				for _, link := range links {
					if strings.Contains(ret, link) {
						ret = link
						found = true
						break
					}
				}
				// HACK: replace Array of Messages with Array of Message.
				if ret == "Messages" && prefix == "" {
					ret = "Message"
				}
				// HACK: if prefix is Array of, add it manually, so ParseType can detect it.
				if found && prefix == "" {
					ret = "Array of " + ret
				}
				if idxSpace := strings.LastIndex(ret, " "); !found && idxSpace > 0 {
					// Skipping verb like "sent".
					ret = ret[idxSpace+1:]
				}
				t := ParseType(ret)
				d.Ret = &t
			} else if strings.Contains(loweredDesc,
				`the edited message is returned, otherwise true is returned`) {
				// HACK: sum type result for editing methods.
				t := ParseType("Message or True")
				d.Ret = &t
			}
		}

		if !s.Is("table") {
			if strings.Contains(d.Description, "Requires no parameters") {
				appendDefinition()
			}
			return
		}

		var head []string
		s.Find("th").Each(func(i int, s *goquery.Selection) {
			head = append(head, strings.TrimSpace(s.Text()))
		})
		if len(head) == 0 {
			return
		}
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			sel := s.Find("td")
			const (
				fName     = 0
				fType     = 1
				fOptional = 2
				optPrefix = "Optional. "
			)
			var fDescription int
			switch sel.Length() {
			case 3:
				fDescription = 2
			case 4:
				fDescription = 3
			default:
				return
			}
			name := sel.Eq(fName).Text()
			rawText := sel.Eq(fDescription).Text()

			optional := strings.HasPrefix(rawText, optPrefix)
			if sel.Eq(fOptional).Text() == "Optional" {
				optional = true
			}
			typ := ParseType(sel.Eq(fType).Text())

			description := selDescription(sel.Eq(fDescription))
			d.Fields = append(d.Fields, Field{
				Name:        name,
				Description: strings.TrimSuffix(strings.TrimPrefix(description, optPrefix), "."),
				Optional:    optional,
				Enum:        collectEnum(typ, name, rawText),
				Type:        typ,
			})
		})
		appendDefinition()
	})
	return a
}

func collectEnum(typ Type, name, description string) []string {
	if typ.Primitive != String || !isDiscriminatorField(name) {
		return nil
	}

	const (
		enumClause  = "can be"
		oneOfClause = "one of"
	)
	idx, _ := IndexOneOf(strings.ToLower(description), enumClause, oneOfClause)
	if idx < 0 {
		return nil
	}
	return collectEnumValues(description[idx:])
}

func collectEnumValues(s string) (r []string) {
	const (
		start = '“'
		end   = '”'
	)

	var (
		i   = 0
		idx = -1
	)
	for i < len(s) {
		c, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case c == start:
			idx = i + size
		case c == end:
			r = append(r, s[idx:i])
			idx = -1
		}
		i += size
	}

	return r
}
