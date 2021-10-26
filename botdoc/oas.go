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
	codeOK      = "200"
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

type stringBound struct {
	Min int64
	Max uint64
}

var boundRegex = regexp.MustCompile(`(\d+)-(\d+) characters`)

func stringBounds(s string) stringBound {
	// 1-32 characters
	matches := boundRegex.FindSubmatch([]byte(s))
	if len(matches) != 3 {
		return stringBound{}
	}
	start, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return stringBound{}
	}
	end, err := strconv.Atoi(string(matches[2]))
	if err != nil {
		return stringBound{}
	}
	return stringBound{Min: int64(start), Max: uint64(end)}
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
			p := ogen.Schema{
				Description: fixTypos(f.Description),
			}
			t := f.Type
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
				case Float:
					p.Type = "number"
				case Boolean:
					p.Type = "boolean"
				}
			case KindObject:
				p.Ref = "#/components/schemas/" + t.Name
			default:
				continue
			}

			if !f.Optional {
				s.Required = append(s.Required, f.Name)
			}

			s.Properties[f.Name] = p
		}
		c.Schemas[d.Name] = s
	}

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
					Schema: c.Schemas[name],
				},
			},
		}
	}

	for _, m := range a.Methods {
		s := ogen.Schema{
			Description: fmt.Sprintf("Input for %s", m.Name),
			Type:        "object",
			Properties:  map[string]ogen.Schema{},
		}
		for _, f := range m.Fields {
			p := ogen.Schema{
				Description: f.Description,
			}
			t := f.Type
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
				case Integer:
					p.Type = "integer"
					if strings.Contains(f.Name, "width") || strings.Contains(f.Name, "height") {
						v := int64(0)
						p.Minimum = &v
						p.ExclusiveMinimum = true
					}
				case Float:
					p.Type = "number"
				case Boolean:
					p.Type = "boolean"
				}
			case KindObject:
				p.Ref = "#/components/schemas/" + t.Name
			default:
				continue
			}

			if !f.Optional {
				s.Required = append(s.Required, f.Name)
			}

			s.Properties[f.Name] = p
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
					codeOK: ogen.Response{Ref: "#/components/responses/" + path.Base(response.Ref)},
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
