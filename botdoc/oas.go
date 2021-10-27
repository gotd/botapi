package botdoc

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/ogen-go/ogen"
)

const (
	contentJSON = "application/json"
)

func resultFor(s ogen.Schema) ogen.Schema {
	return ogen.Schema{
		Type:     "object",
		Required: []string{"ok"},
		Properties: map[string]ogen.Schema{
			"result": s,
			"ok": {
				Type:    "boolean",
				Default: []byte(`true`),
			},
		},
	}
}

type bound struct {
	Min int64
	Max uint64
}

var (
	charBoundRegex = regexp.MustCompile(`(\d+)-(\d+) characters`)
	intBoundRegex  = regexp.MustCompile(`Values between (\d+)-(\d+) are accepted`)
)

func matchBounds(matches [][]byte) (a, b int) {
	start, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return a, b
	}
	end, err := strconv.Atoi(string(matches[2]))
	if err != nil {
		return a, b
	}
	return start, end
}

func regexBounds(r *regexp.Regexp, s string) (a, b int) {
	matches := r.FindSubmatch([]byte(s))
	if len(matches) != 3 {
		return a, b
	}
	return matchBounds(matches)
}

func stringBounds(s string) bound {
	start, end := regexBounds(charBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}

func intBounds(s string) bound {
	start, end := regexBounds(intBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}

func (a API) typeOAS(f Field) *ogen.Schema {
	t := f.Type
	p := &ogen.Schema{
		Description: fixTypos(f.Description),
	}
	switch t.Kind {
	case KindPrimitive:
		switch t.Primitive {
		case String:
			p.Type = "string"

			const defaultMarker = `, must be `
			if idx := strings.LastIndex(p.Description, defaultMarker); idx > 0 {
				// Handle possible default value.
				v := p.Description[idx+len(defaultMarker):]
				if !strings.Contains(p.Description, `one of `) {
					data, err := json.Marshal(v)
					if err != nil {
						panic(err)
					}
					p.Default = data
				}
			}
			b := stringBounds(f.Description)
			if b.Max > 0 {
				p.MaxLength = &b.Max
			}
			if b.Min > 0 {
				p.MinLength = &b.Min
			}

			if strings.Contains(f.Name, "url") {
				p.Format = "uri"
			}
		case Integer:
			p.Type = "integer"
			for _, n := range []string{
				"width",
				"height",
				"duration",
			} {
				if strings.Contains(f.Name, n) {
					v := int64(0)
					p.Minimum = &v
					p.ExclusiveMinimum = true
				}
			}
			if f.Name == "offset" {
				p.Default = []byte(`0`)
			}
			b := intBounds(p.Description)
			if b.Max > 0 {
				v := int64(b.Max)
				p.Maximum = &v
			}
			if b.Min > 0 {
				p.Minimum = &b.Min
			}
		case Float:
			p.Type = "number"
		case Boolean:
			p.Type = "boolean"
		}
	case KindObject:
		p.Ref = "#/components/schemas/" + t.Name
	case KindArray:
		p.Type = "array"
		p.Items = a.typeOAS(Field{Type: *t.Item})
	default:
		if f.Type.String() == "Integer or String" {
			p.Ref = "#/components/schemas/ID"
		} else if f.Type.String() == "String or String" {
			// TODO(ernado): Hack for FileInput, should be removed.
			p.Type = "string"
		} else if len(t.Sum) > 0 {
			for _, s := range t.Sum {
				if one := a.typeOAS(Field{Type: s}); one != nil {
					p.OneOf = append(p.OneOf, *one)
				}
			}
			// TODO(ernado): Implement
			return nil
		} else {
			fmt.Println("unknown", t.Item)
			return nil
		}
	}
	if p.Ref != "" {
		p.Description = ""
	}
	return p
}

func (a API) fieldOAS(parent *ogen.Schema, f Field) *ogen.Schema {
	p := a.typeOAS(f)
	if !f.Optional {
		parent.Required = append(parent.Required, f.Name)
	}
	return p
}

// OAS generates OpenAPI v3 Specification from API definition.
//
//nolint:dupl // TODO(ernado): refactor
func (a API) OAS() *ogen.Spec {
	c := &ogen.Components{
		Schemas:   map[string]ogen.Schema{},
		Responses: map[string]ogen.Response{},
	}
	p := ogen.Paths{}

	for _, d := range a.Types {
		s := ogen.Schema{
			Description: fixTypos(d.Description),
			Type:        "object",
			Properties:  map[string]ogen.Schema{},
		}
		for _, f := range d.Fields {
			if p := a.fieldOAS(&s, f); p != nil {
				s.Properties[f.Name] = *p
			}
		}
		c.Schemas[d.Name] = s
	}

	c.Schemas["InlineKeyboardMarkup"] = ogen.Schema{
		Description: "Hack",
		Type:        "string",
	}
	c.Schemas["InlineQueryResult"] = ogen.Schema{
		Description: "Hack",
		Type:        "string",
	}
	c.Schemas["ID"] = resultFor(ogen.Schema{
		OneOf: []ogen.Schema{
			{Type: "string"},
			{Type: "integer"},
		},
	})
	c.Schemas["Result"] = resultFor(ogen.Schema{
		Type: "boolean",
	})
	c.Schemas["ResultStr"] = resultFor(ogen.Schema{
		Type: "string",
	})
	c.Schemas["ResultMsg"] = resultFor(ogen.Schema{
		Ref: "#/components/schemas/Message",
	})
	c.Schemas["ResultUsr"] = resultFor(ogen.Schema{
		Ref: "#/components/schemas/User",
	})
	c.Schemas["Response"] = ogen.Schema{
		Description: "Contains information about why a request was unsuccessful.",
		Type:        "object",
		Properties: map[string]ogen.Schema{
			"migrate_to_chat_id": {
				Description: "The group has been migrated to a supergroup with the specified identifier. " +
					"This number may be greater than 32 bits and some programming languages may have " +
					"difficulty/silent defects in interpreting it. But it is smaller than 52 bits, " +
					"so a signed 64 bit integer or double-precision float type are safe for storing " +
					"this identifier.",
				Type:   "integer",
				Format: "int64",
			},
			"retry_after": {
				Description: "In case of exceeding flood control, the number of seconds left to wait before the request can be repeated",
				Type:        "integer",
			},
		},
	}
	c.Schemas["Error"] = ogen.Schema{
		Type: "object",
		Required: []string{
			"ok", "error_code", "description",
		},
		Properties: map[string]ogen.Schema{
			"ok": {
				Default: []byte(`false`),
				Type:    "boolean",
			},
			"error_code": {
				Type: "integer",
			},
			"description": {
				Type: "string",
			},
			"parameters": {
				Ref: "#/components/schemas/Response",
			},
		},
	}
	for _, name := range []string{
		"Result",
		"ResultStr",
		"ResultMsg",
		"ResultUsr",
	} {
		c.Responses[name] = ogen.Response{
			Description: "Result of method invocation",
			Content: map[string]ogen.Media{
				contentJSON: {
					Schema: ogen.Schema{
						Ref: "#/components/schemas/" + name,
					},
				},
			},
		}
	}
	c.Responses["Error"] = ogen.Response{
		Description: "Method invocation error",
		Content: map[string]ogen.Media{
			contentJSON: {
				Schema: ogen.Schema{
					Ref: "#/components/schemas/Error",
				},
			},
		},
	}

	for _, m := range a.Methods {
		s := ogen.Schema{
			Description: fmt.Sprintf("Input for %s", m.Name),
			Type:        "object",
			Properties:  map[string]ogen.Schema{},
		}
		for _, f := range m.Fields {
			if p := a.fieldOAS(&s, f); p != nil {
				s.Properties[f.Name] = *p
			}
		}

		schemaName := m.Name
		if len(m.Fields) > 0 {
			c.Schemas[schemaName] = s
		}

		response := ogen.Schema{
			Ref: "#/components/schemas/Result",
		}
		if t := m.Ret; t != nil {
			switch t.Primitive {
			case String:
				response.Ref = "#/components/schemas/ResultStr"
			case Boolean:
				response.Ref = "#/components/schemas/Result"
			}
			switch t.Name {
			case "Message":
				response.Ref = "#/components/schemas/ResultMsg"
			case "User":
				response.Ref = "#/components/schemas/ResultUsr"
			}
			if response.Ref == "" {
				panic("Unable to infer result type")
			}
		}

		var reqBody *ogen.RequestBody
		if len(m.Fields) > 0 {
			reqBody = &ogen.RequestBody{
				Content: map[string]ogen.Media{
					contentJSON: {
						Schema: ogen.Schema{
							Ref: "#/components/schemas/" + schemaName,
						},
					},
				},
				Required: true,
			}
		}
		p["/"+m.Name] = ogen.PathItem{
			Description: m.Description,
			Post: &ogen.Operation{
				OperationID: m.Name,
				RequestBody: reqBody,
				Responses: ogen.Responses{
					"200":     ogen.Response{Ref: "#/components/responses/" + path.Base(response.Ref)},
					"default": ogen.Response{Ref: "#/components/responses/Error"},
				},
			},
		}
	}
	return &ogen.Spec{
		OpenAPI: "3.0.3",
		Info: ogen.Info{
			Title:          "Telegram Bot API",
			TermsOfService: "https://telegram.org/tos",
			Description:    "API for Telegram bots",
			Version:        "5.3",
		},
		Servers: []ogen.Server{
			{
				Description: "production",
				URL:         "https://api.telegram.org/",
			},
		},
		Paths:      p,
		Components: c,
	}
}
