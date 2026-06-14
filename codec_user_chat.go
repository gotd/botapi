package botapi

import "github.com/go-faster/jx"

// Encode writes the user as a JSON object.
func (s *User) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("id")
	e.Int64(s.ID)

	if s.IsBot {
		e.FieldStart("is_bot")
		e.Bool(s.IsBot)
	}

	e.FieldStart("first_name")
	e.Str(s.FirstName)

	if s.LastName != "" {
		e.FieldStart("last_name")
		e.Str(s.LastName)
	}

	if s.Username != "" {
		e.FieldStart("username")
		e.Str(s.Username)
	}

	if s.LanguageCode != "" {
		e.FieldStart("language_code")
		e.Str(s.LanguageCode)
	}

	if s.IsPremium {
		e.FieldStart("is_premium")
		e.Bool(s.IsPremium)
	}

	if s.AddedToAttachmentMenu {
		e.FieldStart("added_to_attachment_menu")
		e.Bool(s.AddedToAttachmentMenu)
	}

	if s.CanJoinGroups {
		e.FieldStart("can_join_groups")
		e.Bool(s.CanJoinGroups)
	}

	if s.CanReadAllGroupMessages {
		e.FieldStart("can_read_all_group_messages")
		e.Bool(s.CanReadAllGroupMessages)
	}

	if s.SupportsInlineQueries {
		e.FieldStart("supports_inline_queries")
		e.Bool(s.SupportsInlineQueries)
	}

	e.ObjEnd()
}

// Decode parses the user from a JSON object.
func (s *User) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "id":
			v, err := d.Int64()
			if err != nil {
				return err
			}

			s.ID = v
		case "is_bot":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.IsBot = v
		case "first_name":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.FirstName = v
		case "last_name":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.LastName = v
		case "username":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Username = v
		case "language_code":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.LanguageCode = v
		case "is_premium":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.IsPremium = v
		case "added_to_attachment_menu":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.AddedToAttachmentMenu = v
		case "can_join_groups":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.CanJoinGroups = v
		case "can_read_all_group_messages":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.CanReadAllGroupMessages = v
		case "supports_inline_queries":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.SupportsInlineQueries = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s User) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *User) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the chat as a JSON object.
func (s *Chat) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("id")
	e.Int64(s.ID)
	e.FieldStart("type")
	e.Str(string(s.Type))

	if s.Title != "" {
		e.FieldStart("title")
		e.Str(s.Title)
	}

	if s.Username != "" {
		e.FieldStart("username")
		e.Str(s.Username)
	}

	if s.FirstName != "" {
		e.FieldStart("first_name")
		e.Str(s.FirstName)
	}

	if s.LastName != "" {
		e.FieldStart("last_name")
		e.Str(s.LastName)
	}

	if s.IsForum {
		e.FieldStart("is_forum")
		e.Bool(s.IsForum)
	}

	e.ObjEnd()
}

// Decode parses the chat from a JSON object.
func (s *Chat) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "id":
			v, err := d.Int64()
			if err != nil {
				return err
			}

			s.ID = v
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = ChatType(v)
		case "title":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Title = v
		case "username":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Username = v
		case "first_name":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.FirstName = v
		case "last_name":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.LastName = v
		case "is_forum":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.IsForum = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Chat) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Chat) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
