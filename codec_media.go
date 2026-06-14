package botapi

import "github.com/go-faster/jx"

// encodeThumbnail writes an optional "thumbnail" field when set.
func encodeThumbnail(e *jx.Encoder, thumb *PhotoSize) {
	if thumb == nil {
		return
	}
	e.FieldStart("thumbnail")
	thumb.Encode(e)
}

// decodePhotoSizePtr decodes a PhotoSize into a freshly allocated pointer.
func decodePhotoSizePtr(d *jx.Decoder, dst **PhotoSize) error {
	v := &PhotoSize{}
	if err := v.Decode(d); err != nil {
		return err
	}
	*dst = v
	return nil
}

// Encode writes the photo size as a JSON object.
func (s *PhotoSize) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("width")
	e.Int(s.Width)
	e.FieldStart("height")
	e.Int(s.Height)
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the photo size from a JSON object.
func (s *PhotoSize) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "width":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Width = v
		case "height":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Height = v
		case "file_size":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s PhotoSize) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *PhotoSize) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the animation as a JSON object.
func (s *Animation) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("width")
	e.Int(s.Width)
	e.FieldStart("height")
	e.Int(s.Height)
	e.FieldStart("duration")
	e.Int(s.Duration)
	encodeThumbnail(e, s.Thumbnail)
	if s.FileName != "" {
		e.FieldStart("file_name")
		e.Str(s.FileName)
	}
	if s.MIMEType != "" {
		e.FieldStart("mime_type")
		e.Str(s.MIMEType)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int64(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the animation from a JSON object.
func (s *Animation) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "width":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Width = v
		case "height":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Height = v
		case "duration":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Duration = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		case "file_name":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileName = v
		case "mime_type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.MIMEType = v
		case "file_size":
			v, err := d.Int64()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Animation) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Animation) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the audio as a JSON object.
func (s *Audio) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("duration")
	e.Int(s.Duration)
	if s.Performer != "" {
		e.FieldStart("performer")
		e.Str(s.Performer)
	}
	if s.Title != "" {
		e.FieldStart("title")
		e.Str(s.Title)
	}
	if s.FileName != "" {
		e.FieldStart("file_name")
		e.Str(s.FileName)
	}
	if s.MIMEType != "" {
		e.FieldStart("mime_type")
		e.Str(s.MIMEType)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int64(s.FileSize)
	}
	encodeThumbnail(e, s.Thumbnail)
	e.ObjEnd()
}

// Decode parses the audio from a JSON object.
func (s *Audio) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "duration":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Duration = v
		case "performer":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.Performer = v
		case "title":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.Title = v
		case "file_name":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileName = v
		case "mime_type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.MIMEType = v
		case "file_size":
			v, err := d.Int64()
			if err != nil {
				return err
			}
			s.FileSize = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Audio) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Audio) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the document as a JSON object.
func (s *Document) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	encodeThumbnail(e, s.Thumbnail)
	if s.FileName != "" {
		e.FieldStart("file_name")
		e.Str(s.FileName)
	}
	if s.MIMEType != "" {
		e.FieldStart("mime_type")
		e.Str(s.MIMEType)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int64(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the document from a JSON object.
func (s *Document) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		case "file_name":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileName = v
		case "mime_type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.MIMEType = v
		case "file_size":
			v, err := d.Int64()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Document) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Document) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the video as a JSON object.
func (s *Video) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("width")
	e.Int(s.Width)
	e.FieldStart("height")
	e.Int(s.Height)
	e.FieldStart("duration")
	e.Int(s.Duration)
	encodeThumbnail(e, s.Thumbnail)
	if s.FileName != "" {
		e.FieldStart("file_name")
		e.Str(s.FileName)
	}
	if s.MIMEType != "" {
		e.FieldStart("mime_type")
		e.Str(s.MIMEType)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int64(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the video from a JSON object.
func (s *Video) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "width":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Width = v
		case "height":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Height = v
		case "duration":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Duration = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		case "file_name":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileName = v
		case "mime_type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.MIMEType = v
		case "file_size":
			v, err := d.Int64()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Video) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Video) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the video note as a JSON object.
func (s *VideoNote) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("length")
	e.Int(s.Length)
	e.FieldStart("duration")
	e.Int(s.Duration)
	encodeThumbnail(e, s.Thumbnail)
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the video note from a JSON object.
func (s *VideoNote) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "length":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Length = v
		case "duration":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Duration = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		case "file_size":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s VideoNote) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *VideoNote) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the voice as a JSON object.
func (s *Voice) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("duration")
	e.Int(s.Duration)
	if s.MIMEType != "" {
		e.FieldStart("mime_type")
		e.Str(s.MIMEType)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int64(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the voice from a JSON object.
func (s *Voice) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "duration":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Duration = v
		case "mime_type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.MIMEType = v
		case "file_size":
			v, err := d.Int64()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Voice) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Voice) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }

// Encode writes the sticker as a JSON object.
func (s *Sticker) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.FieldStart("file_id")
	e.Str(s.FileID)
	e.FieldStart("file_unique_id")
	e.Str(s.FileUniqueID)
	e.FieldStart("type")
	e.Str(string(s.Type))
	e.FieldStart("width")
	e.Int(s.Width)
	e.FieldStart("height")
	e.Int(s.Height)
	if s.IsAnimated {
		e.FieldStart("is_animated")
		e.Bool(s.IsAnimated)
	}
	if s.IsVideo {
		e.FieldStart("is_video")
		e.Bool(s.IsVideo)
	}
	encodeThumbnail(e, s.Thumbnail)
	if s.Emoji != "" {
		e.FieldStart("emoji")
		e.Str(s.Emoji)
	}
	if s.SetName != "" {
		e.FieldStart("set_name")
		e.Str(s.SetName)
	}
	if s.FileSize != 0 {
		e.FieldStart("file_size")
		e.Int(s.FileSize)
	}
	e.ObjEnd()
}

// Decode parses the sticker from a JSON object.
func (s *Sticker) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "file_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileID = v
		case "file_unique_id":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.FileUniqueID = v
		case "type":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.Type = StickerType(v)
		case "width":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Width = v
		case "height":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.Height = v
		case "is_animated":
			v, err := d.Bool()
			if err != nil {
				return err
			}
			s.IsAnimated = v
		case "is_video":
			v, err := d.Bool()
			if err != nil {
				return err
			}
			s.IsVideo = v
		case "thumbnail":
			if err := decodePhotoSizePtr(d, &s.Thumbnail); err != nil {
				return err
			}
		case "emoji":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.Emoji = v
		case "set_name":
			v, err := d.Str()
			if err != nil {
				return err
			}
			s.SetName = v
		case "file_size":
			v, err := d.Int()
			if err != nil {
				return err
			}
			s.FileSize = v
		default:
			return d.Skip()
		}
		return nil
	})
}

// MarshalJSON implements json.Marshaler via jx.
func (s Sticker) MarshalJSON() ([]byte, error) { return marshalJX(&s) }

// UnmarshalJSON implements json.Unmarshaler via jx.
func (s *Sticker) UnmarshalJSON(data []byte) error { return unmarshalJX(data, s) }
