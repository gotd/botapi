// Package botdoc implement types definition extraction from documentation.
package botdoc

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// Kind of Type.
type Kind string

// Possible types.
const (
	KindPrimitive Kind = "primitive"
	KindArray     Kind = "array"
	KindObject    Kind = "object"
	KindSum       Kind = "sum"
)

// Primitive type.
type Primitive string

// Possible primitives.
const (
	String  Primitive = "String"
	Integer Primitive = "Integer"
	Float   Primitive = "Float"
	Boolean Primitive = "Boolean"
)

func (p Primitive) String() string {
	return string(p)
}

// Type of field or parameter.
type Type struct {
	Name      string
	Kind      Kind
	Primitive Primitive
	Item      *Type
	Sum       []Type
}

func newPrimitive(p Primitive) Type {
	return Type{
		Kind:      KindPrimitive,
		Primitive: p,
	}
}

// ParseType parses telegram documentation Type from string
func ParseType(s string) Type {
	switch p := Primitive(s); p {
	case String, Integer, Float, Boolean:
		return newPrimitive(p)
	case "Float number":
		return newPrimitive(Float)
	case "True":
		return newPrimitive(Boolean)
	case "InputFile":
		// TODO(ernado): Implement file upload
		return newPrimitive(String)
	}

	if t := strings.TrimPrefix(s, "Array of "); t != s {
		item := ParseType(t)
		return Type{
			Kind: KindArray,
			Item: &item,
		}
	}

	const sumDelim = " or "
	s = strings.ReplaceAll(s, " and ", sumDelim)
	s = strings.ReplaceAll(s, ", ", sumDelim)
	if strings.Contains(s, sumDelim) {
		t := Type{
			Kind: KindSum,
		}
		for _, e := range strings.Split(s, sumDelim) {
			t.Sum = append(t.Sum, ParseType(e))
		}
		return t
	}

	if strings.Contains(s, " ") {
		// Unknown or invalid type.
		return Type{
			Name: s,
		}
	}
	return Type{
		Kind: KindObject,
		Name: s,
	}
}

func (t Type) String() string {
	switch t.Kind {
	case KindPrimitive:
		return t.Primitive.String()
	case KindObject:
		return t.Name
	case KindArray:
		return fmt.Sprintf("Array of %s", t.Item)
	case KindSum:
		var sum []string
		for _, s := range t.Sum {
			sum = append(sum, s.String())
		}
		return strings.Join(sum, " or ")
	default:
		if t.Name == "" {
			return "unknown"
		}
		return fmt.Sprintf("unknown (%s)", t.Name)
	}
}

// Field of object or argument of function.
type Field struct {
	Type        Type
	Name        string
	Description string
	Enum        []string
	Optional    bool
}

// Definition of structure (method or object).
type Definition struct {
	Name        string
	Description string
	Fields      []Field
	Ret         *Type
}

// API definition.
type API struct {
	Version string
	Types   []Definition
	Methods []Definition
}

type section string

const (
	sectionTypes   = "types"
	sectionMethods = "methods"
)

var typosReplacer = strings.NewReplacer(
	`unpriviledged`, `unprivileged`,
	`Url`, `URL`,
)

func fixTypos(s string) string {
	return typosReplacer.Replace(s)
}

// Extract API definition from goquery document.
func Extract(doc *goquery.Document) (a API) {
	var (
		d   Definition
		sec section
	)
	doc.Find("#dev_page_content").Children().Each(func(i int, s *goquery.Selection) {
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
			d.Description = fixTypos(strings.TrimSpace(s.Text()))
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
			if strings.Contains(d.Description, `It should be one of`) {
				d.Description = strings.TrimSpace(
					strings.ReplaceAll(d.Description, `It should be one of`, ``),
				)
			}
			d.Description = strings.ReplaceAll(d.Description, probablyMarker, "")
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
				retPrefix2      = `returns a`
				retPrefix3      = `returns the`
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
			var definition []string
			s.Find("td").Each(func(j int, s *goquery.Selection) {
				definition = append(definition, strings.TrimSpace(s.Text()))
			})
			const (
				fName     = 0
				fType     = 1
				fOptional = 2
				optPrefix = "Optional. "
			)
			var fDescription int
			switch len(definition) {
			case 3:
				fDescription = 2
			case 4:
				fDescription = 3
			default:
				return
			}
			name := definition[fName]
			description := definition[fDescription]

			optional := strings.HasPrefix(description, optPrefix)
			if definition[fOptional] == "Optional" {
				optional = true
			}
			typ := ParseType(definition[fType])

			d.Fields = append(d.Fields, Field{
				Name:        name,
				Description: strings.TrimSuffix(strings.TrimPrefix(description, optPrefix), "."),
				Optional:    optional,
				Enum:        collectEnum(typ, name, description),
				Type:        typ,
			})
		})
		appendDefinition()
	})
	return a
}

var discriminatorFields = []string{
	"type",
	"source",
	"status",
}

func isDiscriminatorField(n string) bool {
	for _, name := range discriminatorFields {
		if name == n {
			return true
		}
	}
	return false
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
