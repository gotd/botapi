package botapi

import "github.com/go-faster/jx"

// Encode writes the user origin as a JSON object.
func (s *MessageOriginUser) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("date")
	e.Int(s.Date)
	e.FieldStart("sender_user")
	s.SenderUser.Encode(e)
	e.ObjEnd()
}

// Decode parses the user origin from a JSON object.
func (s *MessageOriginUser) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = MessageOriginType(v)
		case "date":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Date = v
		case "sender_user":
			if err := s.SenderUser.Decode(d); err != nil {
				return err
			}
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s *MessageOriginUser) MarshalJSON() ([]byte, error) { return marshalJX(s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *MessageOriginUser) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the hidden-user origin as a JSON object.
func (s *MessageOriginHiddenUser) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("date")
	e.Int(s.Date)
	e.FieldStart("sender_user_name")
	e.Str(s.SenderUserName)
	e.ObjEnd()
}

// Decode parses the hidden-user origin from a JSON object.
func (s *MessageOriginHiddenUser) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = MessageOriginType(v)
		case "date":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Date = v
		case "sender_user_name":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.SenderUserName = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s *MessageOriginHiddenUser) MarshalJSON() ([]byte, error) { return marshalJX(s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *MessageOriginHiddenUser) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the chat origin as a JSON object.
func (s *MessageOriginChat) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("date")
	e.Int(s.Date)
	e.FieldStart("sender_chat")
	s.SenderChat.Encode(e)

	if s.AuthorSignature != "" {
		e.FieldStart("author_signature")
		e.Str(s.AuthorSignature)
	}

	e.ObjEnd()
}

// Decode parses the chat origin from a JSON object.
func (s *MessageOriginChat) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = MessageOriginType(v)
		case "date":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Date = v
		case "sender_chat":
			if err := s.SenderChat.Decode(d); err != nil {
				return err
			}
		case "author_signature":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.AuthorSignature = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s *MessageOriginChat) MarshalJSON() ([]byte, error) { return marshalJX(s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *MessageOriginChat) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the channel origin as a JSON object.
func (s *MessageOriginChannel) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("date")
	e.Int(s.Date)
	e.FieldStart("chat")
	s.Chat.Encode(e)
	e.FieldStart("message_id")
	e.Int(s.MessageID)

	if s.AuthorSignature != "" {
		e.FieldStart("author_signature")
		e.Str(s.AuthorSignature)
	}

	e.ObjEnd()
}

// Decode parses the channel origin from a JSON object.
func (s *MessageOriginChannel) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = MessageOriginType(v)
		case "date":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Date = v
		case "chat":
			if err := s.Chat.Decode(d); err != nil {
				return err
			}
		case "message_id":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.MessageID = v
		case "author_signature":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.AuthorSignature = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s *MessageOriginChannel) MarshalJSON() ([]byte, error) { return marshalJX(s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *MessageOriginChannel) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
