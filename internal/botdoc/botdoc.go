// Package botdoc implement types definition extraction from documentation.
package botdoc

import (
	"fmt"
	"strings"
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
	Type              Type
	Name              string
	RawText           string
	PrettyDescription string
	Enum              []string
	Optional          bool
}

// Definition of structure (method or object).
type Definition struct {
	Name              string
	RawText           string
	PrettyDescription string
	Fields            []Field
	Ret               *Type
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
