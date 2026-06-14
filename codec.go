package botapi

import (
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

// This file wires the receivable entity types to github.com/go-faster/jx so
// they round-trip through JSON: every type in Message's transitive closure
// implements jx-based Encode/Decode, and exposes them to encoding/json through
// MarshalJSON/UnmarshalJSON. The json struct tags on the types remain as
// documentation of the wire field names.

// jsonEncoder is implemented by entity types that serialize themselves through
// the jx streaming encoder.
type jsonEncoder interface {
	Encode(e *jx.Encoder)
}

// jsonDecoder is implemented by entity types that parse themselves from a jx
// streaming decoder.
type jsonDecoder interface {
	Decode(d *jx.Decoder) error
}

// marshalJX renders v through its jx Encode method into a freshly allocated
// JSON document. It is the shared implementation behind every entity's
// MarshalJSON.
func marshalJX(v jsonEncoder) ([]byte, error) {
	e := jx.GetEncoder()
	defer jx.PutEncoder(e)

	v.Encode(e)

	// e.Bytes() aliases the pooled buffer; copy before the buffer is returned
	// to the pool by the deferred PutEncoder.
	return append([]byte(nil), e.Bytes()...), nil
}

// unmarshalJX parses data into v through its jx Decode method. It is the shared
// implementation behind every entity's UnmarshalJSON.
func unmarshalJX(data []byte, v jsonDecoder) error {
	d := jx.GetDecoder()
	defer jx.PutDecoder(d)

	d.ResetBytes(data)

	return v.Decode(d)
}

// decodeMessageOrigin reads a MessageOrigin object, dispatching on its "type"
// discriminator to the matching concrete variant. The decoder must be
// byte-backed (it is, via unmarshalJX's ResetBytes) so Capture can peek at the
// discriminator before the variant is decoded.
func decodeMessageOrigin(d *jx.Decoder) (MessageOrigin, error) {
	var kind string

	if err := d.Capture(func(d *jx.Decoder) error {
		return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
			if string(key) == "type" {
				v, err := d.StrBytes()
				if err != nil {
					return err
				}

				kind = string(v) // copy: StrBytes aliases the read buffer

				return nil
			}

			return d.Skip()
		})
	}); err != nil {
		return nil, err
	}

	var origin MessageOrigin

	switch MessageOriginType(kind) {
	case OriginUser:
		origin = &MessageOriginUser{}
	case OriginHiddenUser:
		origin = &MessageOriginHiddenUser{}
	case OriginChat:
		origin = &MessageOriginChat{}
	case OriginChannel:
		origin = &MessageOriginChannel{}
	default:
		return nil, errors.Errorf("unknown message origin type %q", kind)
	}

	if err := origin.(jsonDecoder).Decode(d); err != nil {
		return nil, err
	}

	return origin, nil
}
