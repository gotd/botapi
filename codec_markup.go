package botapi

import "github.com/go-faster/jx"

// Encode writes the web app info as a JSON object.
func (s *WebAppInfo) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("url")
	e.Str(s.URL)
	e.ObjEnd()
}

// Decode parses the web app info from a JSON object.
func (s *WebAppInfo) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "url":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.URL = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s WebAppInfo) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *WebAppInfo) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the inline keyboard button as a JSON object.
func (s *InlineKeyboardButton) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("text")
	e.Str(s.Text)
	if s.URL != "" {
		e.FieldStart("url")
		e.Str(s.URL)
	}
	if s.CallbackData != "" {
		e.FieldStart("callback_data")
		e.Str(s.CallbackData)
	}
	if s.WebApp != nil {
		e.FieldStart("web_app")
		s.WebApp.Encode(e)
	}
	if s.SwitchInlineQuery != nil {
		e.FieldStart("switch_inline_query")
		e.Str(*s.SwitchInlineQuery)
	}
	if s.SwitchInlineQueryCurrentChat != nil {
		e.FieldStart("switch_inline_query_current_chat")
		e.Str(*s.SwitchInlineQueryCurrentChat)
	}
	if s.Pay {
		e.FieldStart("pay")
		e.Bool(s.Pay)
	}
	e.ObjEnd()
}

// Decode parses the inline keyboard button from a JSON object.
func (s *InlineKeyboardButton) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "text":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.Text = v
		case "url":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.URL = v
		case "callback_data":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.CallbackData = v
		case "web_app":
			w := &WebAppInfo{}
			if err := w.Decode(d); err != nil {
				return err
			}
			s.WebApp = w
		case "switch_inline_query":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.SwitchInlineQuery = &v
		case "switch_inline_query_current_chat":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.SwitchInlineQueryCurrentChat = &v
		case "pay":
			v, err := d.Bool()
			if err != nil {
				return err
			}
			s.Pay = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s InlineKeyboardButton) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *InlineKeyboardButton) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the inline keyboard markup as a JSON object.
func (s *InlineKeyboardMarkup) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("inline_keyboard")
	e.ArrStart()
	for _, row := range s.InlineKeyboard {
		e.ArrStart()
		for i := range row {
			row[i].Encode(e)
		}
		e.ArrEnd()
	}
	e.ArrEnd()
	e.ObjEnd()
}

// Decode parses the inline keyboard markup from a JSON object.
func (s *InlineKeyboardMarkup) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "inline_keyboard":
			return d.Arr(func(d *jx.Decoder) error {
				var row []InlineKeyboardButton
				if err := d.Arr(func(d *jx.Decoder) error {
					var btn InlineKeyboardButton
					if err := btn.Decode(d); err != nil {
						return err
					}
					row = append(row, btn)
					return nil
				}); err != nil {
					return err
				}
				s.InlineKeyboard = append(s.InlineKeyboard, row)
				return nil
			})
		default:
			return d.Skip()
		}
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s InlineKeyboardMarkup) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *InlineKeyboardMarkup) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
