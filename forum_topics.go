package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// GetForumTopicIconStickers returns custom emoji stickers that can be used as a
// forum topic icon by any user.
func (b *Bot) GetForumTopicIconStickers(ctx context.Context) ([]Sticker, error) {
	res, err := b.raw.MessagesGetStickerSet(ctx, &tg.MessagesGetStickerSetRequest{
		Stickerset: &tg.InputStickerSetEmojiDefaultTopicIcons{},
		Hash:       0,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	set, ok := res.(*tg.MessagesStickerSet)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: sticker set not found"}
	}

	out := make([]Sticker, 0, len(set.Documents))

	for _, d := range set.Documents {
		if doc, ok := d.(*tg.Document); ok {
			out = append(out, stickerFromDocument(doc, set.Set.ShortName, StickerCustomEmoji))
		}
	}

	return out, nil
}
