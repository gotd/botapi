package botdoc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ogen-go/ogen"
)

const (
	contentJSON = "application/json"
	codeOK      = "200"
)

// OAS generates OpenAPI v3 Specification from API definition.
func (a API) OAS() *ogen.Spec {
	c := &ogen.Components{
		Schemas: map[string]ogen.Schema{},
	}
	p := ogen.Paths{
		"/getMe": ogen.PathItem{
			Description: "A simple method for testing your bot's auth token. " +
				"Requires no parameters. " +
				"Returns basic information about the bot in form of a User object.",
			Post: &ogen.Operation{
				OperationID: "getMe",
				Responses: ogen.Responses{
					codeOK: ogen.Response{
						Description: "Basic information about the bot",
						Content: map[string]ogen.Media{
							contentJSON: {
								Schema: ogen.Schema{
									Ref: "#/components/schemas/User",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, d := range a.Types {
		s := ogen.Schema{
			Description: d.Description,
			Type:        "object",
			Properties:  map[string]ogen.Schema{},
		}
		for _, f := range d.Fields {
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
				case Float:
					p.Type = "number"
				case Boolean:
					p.Type = "boolean"
				}
			case KindObject:
				if _, ok := c.Schemas[t.Name]; !ok {
					continue
				}
				p.Ref = "#/components/schemas/" + t.Name
			default:
				continue
			}

			if !f.Optional {
				s.Required = append(s.Required, f.Name)
			}

			s.Properties[f.Name] = p
		}
		c.Schemas[m.Name] = s
		p["/"+m.Name] = ogen.PathItem{
			Description: m.Description,
			Post: &ogen.Operation{
				OperationID: m.Name,
				Responses: ogen.Responses{
					codeOK: ogen.Response{},
				},
				RequestBody: &ogen.RequestBody{
					Content: map[string]ogen.Media{
						contentJSON: {
							Schema: ogen.Schema{
								Ref: "#/components/schemas/" + m.Name,
							},
						},
					},
					Required: true,
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
