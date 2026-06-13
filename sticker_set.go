package botapi

import (
	"context"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// GetStickerSet returns a sticker set by its short name.
func (b *Bot) GetStickerSet(ctx context.Context, name string) (*StickerSet, error) {
	res, err := b.raw.MessagesGetStickerSet(ctx, &tg.MessagesGetStickerSetRequest{
		Stickerset: &tg.InputStickerSetShortName{ShortName: name},
		Hash:       0,
	})
	if err != nil {
		return nil, asAPIError(err)
	}
	set, ok := res.(*tg.MessagesStickerSet)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: sticker set not found"}
	}

	typ := stickerSetType(set.Set)
	out := &StickerSet{
		Name:        set.Set.ShortName,
		Title:       set.Set.Title,
		StickerType: typ,
	}
	out.Stickers = make([]Sticker, 0, len(set.Documents))
	for _, d := range set.Documents {
		if doc, ok := d.(*tg.Document); ok {
			out.Stickers = append(out.Stickers, stickerFromDocument(doc, set.Set.ShortName, typ))
		}
	}
	if thumbs := set.Set.Thumbs; len(thumbs) > 0 {
		if ps := photoSizeFromThumb(thumbs[0]); ps != nil {
			out.Thumbnail = ps
		}
	}
	return out, nil
}

// SetStickerSetThumb sets the thumbnail of a sticker set owned by the bot. The
// thumbnail is referenced by file_id or uploaded.
func (b *Bot) SetStickerSetThumb(ctx context.Context, name string, thumb InputFile) error {
	var doc tg.InputDocumentClass
	switch f := thumb.(type) {
	case InputFileID:
		d, err := inputDocumentFromFileID(string(f))
		if err != nil {
			return err
		}
		doc = d
	case *InputFileUpload:
		uploaded, err := b.uploadStickerDocument(ctx, f, StickerFormatStatic)
		if err != nil {
			return err
		}
		doc = &tg.InputDocument{ID: uploaded.ID, AccessHash: uploaded.AccessHash, FileReference: uploaded.FileReference}
	default:
		return &Error{Code: 400, Description: "Bad Request: thumbnail must be a file_id or an uploaded file"}
	}

	if _, err := b.raw.StickersSetStickerSetThumb(ctx, &tg.StickersSetStickerSetThumbRequest{
		Stickerset: &tg.InputStickerSetShortName{ShortName: name},
		Thumb:      doc,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// stickerSetType reports the Bot API sticker type of a set.
func stickerSetType(set tg.StickerSet) StickerType {
	switch {
	case set.Masks:
		return StickerMask
	case set.Emojis:
		return StickerCustomEmoji
	default:
		return StickerRegular
	}
}

// stickerFromDocument converts a sticker document into the Bot API Sticker. typ
// is the type of the enclosing set.
func stickerFromDocument(d *tg.Document, setName string, typ StickerType) Sticker {
	fileID, uniqueID := encodeFileID(fileid.FromDocument(d))
	s := Sticker{
		FileID:       fileID,
		FileUniqueID: uniqueID,
		Type:         typ,
		SetName:      setName,
		FileSize:     int(d.Size),
		IsAnimated:   d.MimeType == mimeStickerAnimated,
		IsVideo:      d.MimeType == mimeStickerVideo,
	}
	for _, attr := range d.Attributes {
		switch a := attr.(type) {
		case *tg.DocumentAttributeImageSize:
			s.Width, s.Height = a.W, a.H
		case *tg.DocumentAttributeVideo:
			s.Width, s.Height = a.W, a.H
		case *tg.DocumentAttributeSticker:
			s.Emoji = a.Alt
		}
	}
	if len(d.Thumbs) > 0 {
		s.Thumbnail = photoSizeFromThumb(d.Thumbs[0])
	}
	return s
}

// photoSizeFromThumb converts a document/sticker thumbnail into a PhotoSize.
// File ids for thumbnails are not derivable from the document alone, so only the
// dimensions and size are filled.
func photoSizeFromThumb(t tg.PhotoSizeClass) *PhotoSize {
	type sized interface {
		GetW() int
		GetH() int
	}
	s, ok := t.(sized)
	if !ok {
		return nil
	}
	ps := &PhotoSize{Width: s.GetW(), Height: s.GetH()}
	if base, ok := t.(*tg.PhotoSize); ok {
		ps.FileSize = base.Size
	}
	return ps
}
