// Package botdoc implement types definition extraction from documentation.
package botdoc

import (
	"fmt"
	"strings"
	"unicode"

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
	Types   []Definition
	Methods []Definition
}

type section string

const (
	sectionTypes   = "types"
	sectionMethods = "methods"
)

// Extract API definition from goquery document.
func Extract(doc *goquery.Document) (a API) {
	var (
		d   Definition
		sec section
	)
	doc.Find("#dev_page_content").Children().Each(func(i int, s *goquery.Selection) {
		if s.Is("h3") {
			switch strings.TrimSpace(s.Text()) {
			case "Available types":
				sec = sectionTypes
			case "Available methods":
				sec = sectionMethods
			}
		}
		appendDefinition := func() {
			if len(d.Name) == 0 {
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
		if s.Is("p") && d.Name != "" && d.Description == "" {
			d.Description = strings.TrimSpace(s.Text())
			if strings.Contains(strings.ToLower(d.Description), `currently holds no information`) {
				appendDefinition()
			}
			return
		}
		if strings.Contains(d.Description, `Returns True on success`) {
			t := newPrimitive(Boolean)
			d.Ret = &t
		}

		// Detect definition of sum-types.
		var probablySum bool
		for _, sumMarker := range []string{
			`It should be one of`,
			`Telegram clients currently support the following`,
			`Currently, the following`,
		} {
			if strings.Contains(d.Description, sumMarker) {
				probablySum = true
			}
		}
		if s.Is("ul") && probablySum {
			t := &Type{
				Kind: KindSum,
			}
			s.Find("li").Each(func(i int, s *goquery.Selection) {
				t.Sum = append(t.Sum, ParseType(s.Text()))
			})
			d.Ret = t
			appendDefinition()
			return
		}

		if d.Ret == nil {
			const (
				retPrefix = `On success, the`
				retSuffix = ` is returned.`
			)
			var (
				start = strings.Index(d.Description, retPrefix)
				end   = strings.Index(d.Description, retSuffix)
			)
			if start > 0 && end > start {
				ret := strings.TrimSpace(d.Description[start+len(retPrefix) : end])
				if idxSpace := strings.LastIndex(ret, " "); idxSpace > 0 {
					// Skipping verb like "sent".
					ret = ret[idxSpace+1:]
				}
				t := ParseType(ret)
				d.Ret = &t
			}
		}

		if !s.Is("table") {
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
			d.Fields = append(d.Fields, Field{
				Name:        definition[fName],
				Description: strings.TrimSuffix(strings.TrimPrefix(definition[fDescription], optPrefix), "."),
				Optional:    strings.HasPrefix(definition[fDescription], optPrefix),
				Type:        ParseType(definition[fType]),
			})
		})
		appendDefinition()
	})
	return a
}
