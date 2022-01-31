package botdoc

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen"
)

const (
	contentJSON = "application/json"
)

func resultFor(s *ogen.Schema) *ogen.Schema {
	return &ogen.Schema{
		Type:     "object",
		Required: []string{"ok"},
		Properties: ogen.Properties{
			{
				Name:   "result",
				Schema: s,
			},
			{
				Name: "ok",
				Schema: &ogen.Schema{
					Type:    "boolean",
					Default: []byte(`true`),
				},
			},
		},
	}
}

func (a API) typeOAS(f Field) *ogen.Schema {
	t := f.Type
	p := &ogen.Schema{
		Description: fixTypos(f.Description),
	}
	for _, value := range f.Enum {
		p.Enum = append(p.Enum, strconv.AppendQuoteToASCII(nil, value))
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
		case Integer:
			p.Type = "integer"
			// Telegram uses int64 (int53, really) for IDs.
			if isIDLikeName(f.Name) || isIDLikeDesc(f.Description) {
				p.Format = "int64"
			}

			if isExclusiveMinimum(f.Name) {
				v := int64(0)
				p.Minimum = &v
				p.ExclusiveMinimum = true
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
		switch t.String() {
		case "Integer or String":
			p.Ref = "#/components/schemas/ID"
		case "String or String":
			p.Type = "string"
		case "InlineKeyboardMarkup or ReplyKeyboardMarkup or ReplyKeyboardRemove or ForceReply":
			p.Ref = "#/components/schemas/SendReplyMarkup"
		default:
			for _, s := range t.Sum {
				p.OneOf = append(p.OneOf, a.typeOAS(Field{Type: s}))
			}
			return p
		}
	}
	if p.Ref != "" {
		p.Description = ""
	}
	return p
}

func (a API) fieldOAS(parent *ogen.Schema, f Field) *ogen.Schema {
	p := a.typeOAS(f)
	if parent != nil && !f.Optional {
		parent.Required = append(parent.Required, f.Name)
	}
	return p
}

func prop(s ogen.Properties, k string) (*ogen.Schema, bool) {
	for _, p := range s {
		if p.Name == k {
			return p.Schema, true
		}
	}
	return nil, false
}

// OAS generates OpenAPI v3 Specification from API definition.
func (a API) OAS() (*ogen.Spec, error) {
	c := &ogen.Components{
		Schemas:   map[string]*ogen.Schema{},
		Responses: map[string]*ogen.Response{},
	}
	p := ogen.Paths{}

	for _, d := range a.Types {
		s := &ogen.Schema{
			Description: fixTypos(d.Description),
			Type:        "object",
		}
		if d.Ret != nil && d.Ret.Kind == KindSum {
			s.Properties = nil
			p := a.typeOAS(Field{Type: *d.Ret})
			s.OneOf = p.OneOf
			s.Type = ""
			c.Schemas[d.Name] = s
			continue
		}
		for _, f := range d.Fields {
			p := a.fieldOAS(s, f)
			if p == nil {
				return nil, errors.Errorf("unable to generate type for %s", f.Type)
			}
			s.Properties = append(s.Properties, ogen.Property{
				Name:   f.Name,
				Schema: p,
			})
		}
		c.Schemas[d.Name] = s
	}

	// Second pass for sum types.
	discriminator := map[string]*ogen.Discriminator{}
Schemas:
	for k, s := range c.Schemas {
		if len(s.OneOf) == 0 {
			continue
		}
		for _, o := range s.OneOf {
			if o.Ref == "" {
				continue Schemas
			}
			target := path.Base(o.Ref)
			one, ok := c.Schemas[target]
			if !ok {
				return nil, errors.Errorf("failed to find %s of %s in schemas", target, k)
			}
			var def []byte

			for _, name := range discriminatorFields {
				p, ok := prop(one.Properties, name)
				if !ok {
					continue
				}

				if s.Discriminator == nil {
					s.Discriminator = &ogen.Discriminator{
						PropertyName: name,
						Mapping:      map[string]string{},
					}
				}

				if len(p.Default) == 0 {
					continue
				}
				def = p.Default

				break
			}
			if len(def) == 0 {
				continue
			}
			discriminator[o.Ref] = s.Discriminator
			v, err := jx.DecodeBytes(def).Str()
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode default")
			}
			s.Discriminator.Mapping[v] = path.Base(o.Ref)
		}
		c.Schemas[k] = s
	}
	c.Schemas["ID"] = &ogen.Schema{
		OneOf: []*ogen.Schema{
			{Type: "string"},
			{Type: "integer", Format: "int64"},
		},
	}
	c.Schemas["Result"] = resultFor(&ogen.Schema{
		Type: "boolean",
	})
	c.Schemas["ResultString"] = resultFor(&ogen.Schema{
		Type: "string",
	})
	c.Schemas["ResultInt"] = resultFor(&ogen.Schema{
		Type: "integer",
	})
	c.Schemas["SendReplyMarkup"] = &ogen.Schema{
		OneOf: []*ogen.Schema{
			{Ref: "#/components/schemas/InlineKeyboardMarkup"},
			{Ref: "#/components/schemas/ReplyKeyboardMarkup"},
			{Ref: "#/components/schemas/ReplyKeyboardRemove"},
			{Ref: "#/components/schemas/ForceReply"},
		},
	}
	addResponse := func(name, ref, description string) {
		c.Responses[name] = &ogen.Response{
			Description: description,
			Content: map[string]ogen.Media{
				contentJSON: {
					Schema: &ogen.Schema{
						Ref: ref,
					},
				},
			},
		}
	}

	for _, t := range wellKnownTypes {
		resultName := "Result" + t
		c.Schemas[resultName] = resultFor(&ogen.Schema{
			Ref: "#/components/schemas/" + t,
		})
		addResponse(resultName, "#/components/schemas/"+resultName, "Result of method invocation")
	}
	c.Schemas["Response"] = &ogen.Schema{
		Description: "Contains information about why a request was unsuccessful.",
		Type:        "object",
		Properties: ogen.Properties{
			{
				Name: "migrate_to_chat_id",
				Schema: &ogen.Schema{
					Description: "The group has been migrated to a supergroup with the specified identifier. " +
						"This number may be greater than 32 bits and some programming languages may have " +
						"difficulty/silent defects in interpreting it. But it is smaller than 52 bits, " +
						"so a signed 64 bit integer or double-precision float type are safe for storing " +
						"this identifier.",
					Type:   "integer",
					Format: "int64",
				},
			},
			{
				Name: "retry_after",
				Schema: &ogen.Schema{
					Description: "In case of exceeding flood control, the number of seconds left to wait before the request can be repeated",
					Type:        "integer",
				},
			},
		},
	}
	c.Schemas["Error"] = &ogen.Schema{
		Type: "object",
		Required: []string{
			"ok", "error_code", "description",
		},
		Properties: ogen.Properties{
			{
				Name: "ok",
				Schema: &ogen.Schema{
					Default: []byte(`false`),
					Type:    "boolean",
				},
			},
			{
				Name: "error_code",
				Schema: &ogen.Schema{
					Type: "integer",
				},
			},
			{
				Name: "description",
				Schema: &ogen.Schema{
					Type: "string",
				},
			},
			{
				Name: "parameters",
				Schema: &ogen.Schema{
					Ref: "#/components/schemas/Response",
				},
			},
		},
	}
	for _, name := range []string{
		"Result",
		"ResultString",
		"ResultInt",
	} {
		addResponse(name, "#/components/schemas/"+name, "Result of method invocation")
	}
	c.Responses["Error"] = &ogen.Response{
		Description: "Method invocation error",
		Content: map[string]ogen.Media{
			contentJSON: {
				Schema: &ogen.Schema{
					Ref: "#/components/schemas/Error",
				},
			},
		},
	}

	for _, m := range a.Methods {
		s := &ogen.Schema{
			Description: fmt.Sprintf("Input for %s", m.Name),
			Type:        "object",
		}
		var hasConditionalRequired bool
		for _, f := range m.Fields {
			if strings.Contains(f.Description, "Required if") {
				hasConditionalRequired = true
			}
			p := a.fieldOAS(s, f)
			oneOf := p.OneOf
			if p.Items != nil {
				oneOf = p.Items.OneOf
			}
			var d *ogen.Discriminator
			for _, s := range oneOf {
				if d != nil {
					break
				}
				d = discriminator[s.Ref]
			}
			if d != nil {
				df := &ogen.Discriminator{
					PropertyName: d.PropertyName,
					Mapping:      map[string]string{},
				}
				// Copy only existing variants of oneOf.
				for _, o := range oneOf {
					for k, v := range d.Mapping {
						if v == path.Base(o.Ref) {
							df.Mapping[k] = v
						}
					}
				}
				if p.Items != nil {
					p.Items.Discriminator = df
				} else {
					p.Discriminator = df
				}
			}
			s.Properties = append(s.Properties, ogen.Property{
				Name:   f.Name,
				Schema: p,
			})
		}

		schemaName := m.Name
		if len(m.Fields) > 0 {
			c.Schemas[schemaName] = s
		}

		response := &ogen.Schema{
			Ref: "#/components/schemas/Result",
		}
		if t := m.Ret; t != nil {
			getRef := func(t *Type) string {
				switch t.Kind {
				case KindPrimitive:
					switch t.Primitive {
					case String:
						return "#/components/schemas/ResultString"
					case Integer:
						return "#/components/schemas/ResultInt"
					case Boolean:
						return "#/components/schemas/Result"
					}
				case KindObject:
					for _, typ := range wellKnownTypes {
						if typ == t.Name {
							return "#/components/schemas/Result" + typ
						}
					}
				}

				return "#/components/schemas/Result"
			}
			switch t.Kind {
			case KindPrimitive, KindObject:
				response.Ref = getRef(t)
			case KindArray:
				itemRef := getRef(t.Item)
				resultName := "ResultArrayOf" + t.Item.Name
				itemName := strings.ReplaceAll(itemRef,
					`#/components/schemas/Result`,
					`#/components/schemas/`,
				)
				c.Schemas[resultName] = resultFor(&ogen.Schema{
					Type: "array",
					Items: &ogen.Schema{
						Ref: itemName,
					},
				})
				addResponse(resultName, "#/components/schemas/"+resultName, "Result of method invocation")
				response.Ref = "#/components/responses/" + resultName
			case KindSum:
				var name strings.Builder
				for i, v := range t.Sum {
					variantName := v.Name
					if variantName == "" {
						variantName = v.Primitive.String()
					}
					if variantName == "" {
						return nil, errors.Errorf("unsupported type %v", v.Kind)
					}
					name.WriteString(variantName)
					if i != len(t.Sum)-1 {
						name.WriteString("Or")
					}
				}
				n := name.String()
				resultName := "Result" + n
				c.Schemas[resultName] = resultFor(a.typeOAS(Field{
					Type: *t,
					Name: n,
				}))
				addResponse(resultName, "#/components/schemas/"+resultName, "Result of method invocation")
				response.Ref = "#/components/responses/" + resultName
			}

			if response.Ref == "" {
				panic("Unable to infer result type")
			}
		}

		requestSchema := ogen.Schema{Ref: "#/components/schemas/" + schemaName}
		var reqBody *ogen.RequestBody
		if len(m.Fields) > 0 {
			reqBody = &ogen.RequestBody{
				Content: map[string]ogen.Media{
					contentJSON: {Schema: &requestSchema},
				},
				Required: len(s.Required) > 0 || hasConditionalRequired,
			}
		}
		responses := ogen.Responses{
			"200":     &ogen.Response{Ref: "#/components/responses/" + path.Base(response.Ref)},
			"default": &ogen.Response{Ref: "#/components/responses/Error"},
		}
		item := &ogen.PathItem{
			Description: m.Description,
			Post: &ogen.Operation{
				OperationID: m.Name,
				RequestBody: reqBody,
				Responses:   responses,
			},
		}
		p["/"+m.Name] = item
	}
	return patchSchema(&ogen.Spec{
		OpenAPI: "3.0.3",
		Info: ogen.Info{
			Title:          "Telegram Bot API",
			TermsOfService: "https://telegram.org/tos",
			Description:    "API for Telegram bots",
			Version:        a.Version,
		},
		Servers: []ogen.Server{
			{
				Description: "production",
				URL:         "https://api.telegram.org/",
			},
		},
		Paths:      p,
		Components: c,
	}), nil
}

func patchSchema(spec *ogen.Spec) *ogen.Spec {
	c := spec.Components
	addMissedProperties(c.Schemas)
	updateProperty := func(typeName, propName string, cb func(p *ogen.Property)) {
		schema := c.Schemas[typeName]
		props := schema.Properties

		for i := range props {
			if props[i].Name == propName {
				cb(&props[i])
				return
			}
		}
		panic(fmt.Sprintf("property %q of %q not found", propName, typeName))
	}
	setItemsFormat := func(typeName, propName, format string) {
		updateProperty(typeName, propName, func(p *ogen.Property) {
			p.Schema.Items.Format = format
		})
	}
	// MTProto uses int64, use it in BotAPI too to reduce copying.
	setItemsFormat("sendInvoice", "suggested_tip_amounts", "int64")
	setItemsFormat("InputInvoiceMessageContent", "suggested_tip_amounts", "int64")
	return spec
}

func addMissedProperties(schemas map[string]*ogen.Schema) {
	add := func(name string, props ...ogen.Property) {
		schemas[name].Properties = append(schemas[name].Properties, props...)
	}

	add("Chat", ogen.Property{
		Name: "all_members_are_administrators",
		Schema: &ogen.Schema{
			Type: "boolean",
		},
	})
	// Seems like legacy fields.
	add("Message", ogen.Property{
		Name: "new_chat_member",
		Schema: &ogen.Schema{
			Ref: "#/components/schemas/User",
		},
	}, ogen.Property{
		Name: "new_chat_participant",
		Schema: &ogen.Schema{
			Ref: "#/components/schemas/User",
		},
	}, ogen.Property{
		Name: "left_chat_participant",
		Schema: &ogen.Schema{
			Ref: "#/components/schemas/User",
		},
	})
}
