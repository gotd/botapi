package botdoc

import (
	"github.com/ogen-go/ogen"
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
			Get: &ogen.Operation{
				OperationID: "getMe",
				Responses: ogen.Responses{
					"200": ogen.Response{
						Description: "Basic information about the bot",
						Content: map[string]ogen.Media{
							"application/json": {
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
				case Integer:
					p.Type = "int"
				case Float:
					p.Type = "float64"
				case Boolean:
					p.Type = "bool"
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
