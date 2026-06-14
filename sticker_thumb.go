package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// SetStickerSetThumbnail sets the thumbnail of a regular or mask sticker set
// owned by the bot. The thumbnail is referenced by file_id or uploaded; format
// selects the uploaded file's encoding. This is the format-aware counterpart of
// SetStickerSetThumb.
func (b *Bot) SetStickerSetThumbnail(ctx context.Context, name string, thumb InputFile, format StickerFormat) error {
	var doc tg.InputDocumentClass

	switch f := thumb.(type) {
	case InputFileID:
		d, err := inputDocumentFromFileID(string(f))
		if err != nil {
			return err
		}

		doc = d
	case *InputFileUpload:
		uploaded, err := b.uploadStickerDocument(ctx, f, format)
		if err != nil {
			return err
		}

		doc = &tg.InputDocument{ID: uploaded.ID, AccessHash: uploaded.AccessHash, FileReference: uploaded.FileReference}
	default:
		return &Error{Code: 400, Description: "Bad Request: thumbnail must be a file_id or an uploaded file"}
	}

	req := &tg.StickersSetStickerSetThumbRequest{Stickerset: &tg.InputStickerSetShortName{ShortName: name}}
	req.SetThumb(doc)

	if _, err := b.raw.StickersSetStickerSetThumb(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetCustomEmojiStickerSetThumbnail sets the thumbnail of a custom emoji sticker
// set owned by the bot to one of its stickers, referenced by custom_emoji_id. An
// empty customEmojiID drops the set thumbnail.
func (b *Bot) SetCustomEmojiStickerSetThumbnail(ctx context.Context, name, customEmojiID string) error {
	req := &tg.StickersSetStickerSetThumbRequest{Stickerset: &tg.InputStickerSetShortName{ShortName: name}}

	if customEmojiID != "" {
		id, err := strconv.ParseInt(customEmojiID, 10, 64)
		if err != nil {
			return errInvalidCustomEmojiID()
		}

		req.SetThumbDocumentID(id)
	}

	if _, err := b.raw.StickersSetStickerSetThumb(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}
