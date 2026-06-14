package botapi

import "github.com/go-faster/jx"

// Encode writes the message as a JSON object, omitting zero-value optional
// fields to match the Bot API wire format.
func (s *Message) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("message_id")
	e.Int(s.MessageID)

	if s.MessageThreadID != 0 {
		e.FieldStart("message_thread_id")
		e.Int(s.MessageThreadID)
	}

	if s.From != nil {
		e.FieldStart("from")
		s.From.Encode(e)
	}

	if s.SenderChat != nil {
		e.FieldStart("sender_chat")
		s.SenderChat.Encode(e)
	}

	e.FieldStart("date")
	e.Int(s.Date)
	e.FieldStart("chat")
	s.Chat.Encode(e)

	if s.ForwardOrigin != nil {
		e.FieldStart("forward_origin")
		s.ForwardOrigin.Encode(e)
	}

	if s.ReplyToMessage != nil {
		e.FieldStart("reply_to_message")
		s.ReplyToMessage.Encode(e)
	}

	if s.ViaBot != nil {
		e.FieldStart("via_bot")
		s.ViaBot.Encode(e)
	}

	if s.EditDate != 0 {
		e.FieldStart("edit_date")
		e.Int(s.EditDate)
	}

	if s.HasProtectedContent {
		e.FieldStart("has_protected_content")
		e.Bool(s.HasProtectedContent)
	}

	if s.MediaGroupID != "" {
		e.FieldStart("media_group_id")
		e.Str(s.MediaGroupID)
	}

	if s.AuthorSignature != "" {
		e.FieldStart("author_signature")
		e.Str(s.AuthorSignature)
	}

	if s.Text != "" {
		e.FieldStart("text")
		e.Str(s.Text)
	}

	encodeEntities(e, "entities", s.Entities)

	if s.Caption != "" {
		e.FieldStart("caption")
		e.Str(s.Caption)
	}

	encodeEntities(e, "caption_entities", s.CaptionEntities)

	if s.Animation != nil {
		e.FieldStart("animation")
		s.Animation.Encode(e)
	}

	if s.Audio != nil {
		e.FieldStart("audio")
		s.Audio.Encode(e)
	}

	if s.Document != nil {
		e.FieldStart("document")
		s.Document.Encode(e)
	}

	if len(s.Photo) > 0 {
		e.FieldStart("photo")
		e.ArrStart()

		for i := range s.Photo {
			s.Photo[i].Encode(e)
		}

		e.ArrEnd()
	}

	if s.Sticker != nil {
		e.FieldStart("sticker")
		s.Sticker.Encode(e)
	}

	if s.Video != nil {
		e.FieldStart("video")
		s.Video.Encode(e)
	}

	if s.VideoNote != nil {
		e.FieldStart("video_note")
		s.VideoNote.Encode(e)
	}

	if s.Voice != nil {
		e.FieldStart("voice")
		s.Voice.Encode(e)
	}

	if s.Contact != nil {
		e.FieldStart("contact")
		s.Contact.Encode(e)
	}

	if s.Dice != nil {
		e.FieldStart("dice")
		s.Dice.Encode(e)
	}

	if s.Poll != nil {
		e.FieldStart("poll")
		s.Poll.Encode(e)
	}

	if s.Venue != nil {
		e.FieldStart("venue")
		s.Venue.Encode(e)
	}

	if s.Location != nil {
		e.FieldStart("location")
		s.Location.Encode(e)
	}

	if len(s.NewChatMembers) > 0 {
		e.FieldStart("new_chat_members")
		e.ArrStart()

		for i := range s.NewChatMembers {
			s.NewChatMembers[i].Encode(e)
		}

		e.ArrEnd()
	}

	if s.LeftChatMember != nil {
		e.FieldStart("left_chat_member")
		s.LeftChatMember.Encode(e)
	}

	if s.NewChatTitle != "" {
		e.FieldStart("new_chat_title")
		e.Str(s.NewChatTitle)
	}

	if s.PinnedMessage != nil {
		e.FieldStart("pinned_message")
		s.PinnedMessage.Encode(e)
	}

	if s.ReplyMarkup != nil {
		e.FieldStart("reply_markup")
		s.ReplyMarkup.Encode(e)
	}

	e.ObjEnd()
}

// Decode parses the message from a JSON object, resolving the polymorphic
// forward_origin field to its concrete MessageOrigin variant.
func (s *Message) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "message_id":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.MessageID = v
		case "message_thread_id":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.MessageThreadID = v
		case "from":
			u := &User{}
			if err := u.Decode(d); err != nil {
				return err
			}

			s.From = u
		case "sender_chat":
			c := &Chat{}
			if err := c.Decode(d); err != nil {
				return err
			}

			s.SenderChat = c
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
		case "forward_origin":
			origin, err := decodeMessageOrigin(d)
			if err != nil {
				return err
			}

			s.ForwardOrigin = origin
		case "reply_to_message":
			m := &Message{}
			if err := m.Decode(d); err != nil {
				return err
			}

			s.ReplyToMessage = m
		case "via_bot":
			u := &User{}
			if err := u.Decode(d); err != nil {
				return err
			}

			s.ViaBot = u
		case "edit_date":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.EditDate = v
		case "has_protected_content":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.HasProtectedContent = v
		case "media_group_id":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.MediaGroupID = v
		case "author_signature":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.AuthorSignature = v
		case "text":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Text = v
		case "entities":
			if err := decodeEntities(d, &s.Entities); err != nil {
				return err
			}
		case "caption":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Caption = v
		case "caption_entities":
			if err := decodeEntities(d, &s.CaptionEntities); err != nil {
				return err
			}
		case "animation":
			v := &Animation{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Animation = v
		case "audio":
			v := &Audio{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Audio = v
		case "document":
			v := &Document{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Document = v
		case "photo":
			if err := d.Arr(func(d *jx.Decoder) error {
				var p PhotoSize

				if err := p.Decode(d); err != nil {
					return err
				}

				s.Photo = append(s.Photo, p)

				return nil
			}); err != nil {
				return err
			}
		case "sticker":
			v := &Sticker{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Sticker = v
		case "video":
			v := &Video{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Video = v
		case "video_note":
			v := &VideoNote{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.VideoNote = v
		case "voice":
			v := &Voice{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Voice = v
		case "contact":
			v := &Contact{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Contact = v
		case "dice":
			v := &Dice{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Dice = v
		case "poll":
			v := &Poll{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Poll = v
		case "venue":
			v := &Venue{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Venue = v
		case "location":
			v := &Location{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.Location = v
		case "new_chat_members":
			if err := d.Arr(func(d *jx.Decoder) error {
				var u User

				if err := u.Decode(d); err != nil {
					return err
				}

				s.NewChatMembers = append(s.NewChatMembers, u)

				return nil
			}); err != nil {
				return err
			}
		case "left_chat_member":
			u := &User{}
			if err := u.Decode(d); err != nil {
				return err
			}

			s.LeftChatMember = u
		case "new_chat_title":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.NewChatTitle = v
		case "pinned_message":
			m := &Message{}
			if err := m.Decode(d); err != nil {
				return err
			}

			s.PinnedMessage = m
		case "reply_markup":
			v := &InlineKeyboardMarkup{}
			if err := v.Decode(d); err != nil {
				return err
			}

			s.ReplyMarkup = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s *Message) MarshalJSON() ([]byte, error) { return marshalJX(s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Message) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
