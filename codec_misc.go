package botapi

import "github.com/go-faster/jx"

// Encode writes the message entity as a JSON object.
func (s *MessageEntity) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("offset")
	e.Int(s.Offset)
	e.FieldStart("length")
	e.Int(s.Length)

	if s.URL != "" {
		e.FieldStart("url")
		e.Str(s.URL)
	}

	if s.User != nil {
		e.FieldStart("user")
		s.User.Encode(e)
	}

	if s.Language != "" {
		e.FieldStart("language")
		e.Str(s.Language)
	}

	if s.CustomEmojiID != "" {
		e.FieldStart("custom_emoji_id")
		e.Str(s.CustomEmojiID)
	}

	e.ObjEnd()
}

// Decode parses the message entity from a JSON object.
func (s *MessageEntity) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = MessageEntityType(v)
		case "offset":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Offset = v
		case "length":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Length = v
		case "url":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.URL = v
		case "user":
			u := &User{}
			if err := u.Decode(d); err != nil {
				return err
			}

			s.User = u
		case "language":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Language = v
		case "custom_emoji_id":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.CustomEmojiID = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s MessageEntity) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *MessageEntity) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// encodeEntities writes a "field" of message entities when the slice is
// non-empty.
func encodeEntities(e *jx.Encoder, field string, entities []MessageEntity) {
	if len(entities) == 0 {
		return
	}

	e.FieldStart(field)
	e.ArrStart()

	for i := range entities {
		entities[i].Encode(e)
	}

	e.ArrEnd()
}

// decodeEntities decodes an array of message entities.
func decodeEntities(d *jx.Decoder, dst *[]MessageEntity) error {
	return d.Arr(func(d *jx.Decoder) error {
		var ent MessageEntity

		if err := ent.Decode(d); err != nil {
			return err
		}

		*dst = append(*dst, ent)

		return nil
	})
}

// Encode writes the contact as a JSON object.
func (s *Contact) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("phone_number")
	e.Str(s.PhoneNumber)
	e.FieldStart("first_name")
	e.Str(s.FirstName)

	if s.LastName != "" {
		e.FieldStart("last_name")
		e.Str(s.LastName)
	}

	if s.UserID != 0 {
		e.FieldStart("user_id")
		e.Int64(s.UserID)
	}

	if s.VCard != "" {
		e.FieldStart("vcard")
		e.Str(s.VCard)
	}

	e.ObjEnd()
}

// Decode parses the contact from a JSON object.
func (s *Contact) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "phone_number":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.PhoneNumber = v
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
		case "user_id":
			v, err := d.Int64()
			if err != nil {
				return err
			}

			s.UserID = v
		case "vcard":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.VCard = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Contact) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Contact) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the dice as a JSON object.
func (s *Dice) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("emoji")
	e.Str(string(s.Emoji))
	e.FieldStart("value")
	e.Int(s.Value)
	e.ObjEnd()
}

// Decode parses the dice from a JSON object.
func (s *Dice) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "emoji":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Emoji = DiceEmoji(v)
		case "value":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Value = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Dice) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Dice) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the location as a JSON object.
func (s *Location) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("longitude")
	e.Float64(s.Longitude)
	e.FieldStart("latitude")
	e.Float64(s.Latitude)

	if s.HorizontalAccuracy != 0 {
		e.FieldStart("horizontal_accuracy")
		e.Float64(s.HorizontalAccuracy)
	}

	if s.LivePeriod != 0 {
		e.FieldStart("live_period")
		e.Int(s.LivePeriod)
	}

	if s.Heading != 0 {
		e.FieldStart("heading")
		e.Int(s.Heading)
	}

	if s.ProximityAlertRadius != 0 {
		e.FieldStart("proximity_alert_radius")
		e.Int(s.ProximityAlertRadius)
	}

	e.ObjEnd()
}

// Decode parses the location from a JSON object.
func (s *Location) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "longitude":
			v, err := d.Float64()
			if err != nil {
				return err
			}

			s.Longitude = v
		case "latitude":
			v, err := d.Float64()
			if err != nil {
				return err
			}

			s.Latitude = v
		case "horizontal_accuracy":
			v, err := d.Float64()
			if err != nil {
				return err
			}

			s.HorizontalAccuracy = v
		case "live_period":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.LivePeriod = v
		case "heading":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.Heading = v
		case "proximity_alert_radius":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.ProximityAlertRadius = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Location) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Location) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the venue as a JSON object.
func (s *Venue) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("location")
	s.Location.Encode(e)
	e.FieldStart("title")
	e.Str(s.Title)
	e.FieldStart("address")
	e.Str(s.Address)

	if s.FoursquareID != "" {
		e.FieldStart("foursquare_id")
		e.Str(s.FoursquareID)
	}

	if s.FoursquareType != "" {
		e.FieldStart("foursquare_type")
		e.Str(s.FoursquareType)
	}

	if s.GooglePlaceID != "" {
		e.FieldStart("google_place_id")
		e.Str(s.GooglePlaceID)
	}

	if s.GooglePlaceType != "" {
		e.FieldStart("google_place_type")
		e.Str(s.GooglePlaceType)
	}

	e.ObjEnd()
}

// Decode parses the venue from a JSON object.
func (s *Venue) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "location":
			if err := s.Location.Decode(d); err != nil {
				return err
			}
		case "title":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Title = v
		case "address":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Address = v
		case "foursquare_id":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.FoursquareID = v
		case "foursquare_type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.FoursquareType = v
		case "google_place_id":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.GooglePlaceID = v
		case "google_place_type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.GooglePlaceType = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Venue) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Venue) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the poll option as a JSON object.
func (s *PollOption) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("text")
	e.Str(s.Text)
	e.FieldStart("voter_count")
	e.Int(s.VoterCount)
	e.ObjEnd()
}

// Decode parses the poll option from a JSON object.
func (s *PollOption) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "text":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Text = v
		case "voter_count":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.VoterCount = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s PollOption) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *PollOption) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the poll as a JSON object.
func (s *Poll) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("id")
	e.Str(s.ID)
	e.FieldStart("question")
	e.Str(s.Question)
	e.FieldStart("options")
	e.ArrStart()

	for i := range s.Options {
		s.Options[i].Encode(e)
	}

	e.ArrEnd()
	e.FieldStart("total_voter_count")
	e.Int(s.TotalVoterCount)

	if s.IsClosed {
		e.FieldStart("is_closed")
		e.Bool(s.IsClosed)
	}

	if s.IsAnonymous {
		e.FieldStart("is_anonymous")
		e.Bool(s.IsAnonymous)
	}

	e.FieldStart("type")
	e.Str(string(s.Type))

	if s.AllowsMultipleAnswers {
		e.FieldStart("allows_multiple_answers")
		e.Bool(s.AllowsMultipleAnswers)
	}

	if s.CorrectOptionID != 0 {
		e.FieldStart("correct_option_id")
		e.Int(s.CorrectOptionID)
	}

	if s.Explanation != "" {
		e.FieldStart("explanation")
		e.Str(s.Explanation)
	}

	encodeEntities(e, "explanation_entities", s.ExplanationEntities)

	if s.OpenPeriod != 0 {
		e.FieldStart("open_period")
		e.Int(s.OpenPeriod)
	}

	e.ObjEnd()
}

// Decode parses the poll from a JSON object.
func (s *Poll) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "id":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.ID = v
		case "question":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Question = v
		case "options":
			if err := d.Arr(func(d *jx.Decoder) error {
				var o PollOption

				if err := o.Decode(d); err != nil {
					return err
				}

				s.Options = append(s.Options, o)

				return nil
			}); err != nil {
				return err
			}
		case "total_voter_count":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.TotalVoterCount = v
		case "is_closed":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.IsClosed = v
		case "is_anonymous":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.IsAnonymous = v
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Type = PollType(v)
		case "allows_multiple_answers":
			v, err := d.Bool()
			if err != nil {
				return err
			}

			s.AllowsMultipleAnswers = v
		case "correct_option_id":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.CorrectOptionID = v
		case "explanation":
			v, err := d.Str()
			if err != nil {
				return err
			}

			s.Explanation = v
		case "explanation_entities":
			if err := decodeEntities(d, &s.ExplanationEntities); err != nil {
				return err
			}
		case "open_period":
			v, err := d.Int()
			if err != nil {
				return err
			}

			s.OpenPeriod = v
		default:
			return d.Skip()
		}

		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Poll) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Poll) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
