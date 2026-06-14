package botapi

import (
	"context"
	"strconv"
	"strings"

	"github.com/gotd/td/tg"
)

// MaskPosition describes the position on faces where a mask sticker is placed.
type MaskPosition struct {
	// Point is the part of the face relative to which the mask is placed: one of
	// "forehead", "eyes", "mouth" or "chin".
	Point string `json:"point"`
	// XShift is the horizontal shift, in widths of the mask scaled to the face
	// size, left to right. Negative shifts the mask left.
	XShift float64 `json:"x_shift"`
	// YShift is the vertical shift, in heights of the mask scaled to the face
	// size, top to bottom. Negative shifts the mask up.
	YShift float64 `json:"y_shift"`
	// Scale is the mask scaling coefficient.
	Scale float64 `json:"scale"`
}

// maskPoints maps a Bot API mask point to the MTProto MaskCoords.n index.
var maskPoints = map[string]int{
	"forehead": 0,
	"eyes":     1,
	"mouth":    2,
	"chin":     3,
}

// toMaskCoords converts the Bot API mask position to MTProto mask coordinates.
func (m MaskPosition) toMaskCoords() (tg.MaskCoords, error) {
	n, ok := maskPoints[m.Point]
	if !ok {
		return tg.MaskCoords{}, &Error{Code: 400, Description: "Bad Request: invalid mask position point"}
	}

	return tg.MaskCoords{N: n, X: m.XShift, Y: m.YShift, Zoom: m.Scale}, nil
}

// DeleteStickerSet deletes a sticker set created by the bot, referenced by its
// short name.
func (b *Bot) DeleteStickerSet(ctx context.Context, name string) error {
	if _, err := b.raw.StickersDeleteStickerSet(ctx, &tg.InputStickerSetShortName{ShortName: name}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetStickerSetTitle sets the title of a sticker set created by the bot.
func (b *Bot) SetStickerSetTitle(ctx context.Context, name, title string) error {
	if _, err := b.raw.StickersRenameStickerSet(ctx, &tg.StickersRenameStickerSetRequest{
		Stickerset: &tg.InputStickerSetShortName{ShortName: name},
		Title:      title,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetStickerEmojiList changes the list of emoji assigned to a sticker created by
// the bot, referenced by file_id.
func (b *Bot) SetStickerEmojiList(ctx context.Context, sticker string, emojiList []string) error {
	doc, err := inputDocumentFromFileID(sticker)
	if err != nil {
		return err
	}

	req := &tg.StickersChangeStickerRequest{Sticker: doc}
	req.SetEmoji(strings.Join(emojiList, ""))

	if _, err := b.raw.StickersChangeSticker(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetStickerKeywords changes the search keywords assigned to a sticker created
// by the bot, referenced by file_id.
func (b *Bot) SetStickerKeywords(ctx context.Context, sticker string, keywords []string) error {
	doc, err := inputDocumentFromFileID(sticker)
	if err != nil {
		return err
	}

	req := &tg.StickersChangeStickerRequest{Sticker: doc}
	req.SetKeywords(strings.Join(keywords, ","))

	if _, err := b.raw.StickersChangeSticker(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SetStickerMaskPosition changes the mask position of a mask sticker created by
// the bot, referenced by file_id. A nil position clears the mask position.
func (b *Bot) SetStickerMaskPosition(ctx context.Context, sticker string, position *MaskPosition) error {
	doc, err := inputDocumentFromFileID(sticker)
	if err != nil {
		return err
	}

	req := &tg.StickersChangeStickerRequest{Sticker: doc}

	if position != nil {
		coords, err := position.toMaskCoords()
		if err != nil {
			return err
		}

		req.SetMaskCoords(coords)
	}

	if _, err := b.raw.StickersChangeSticker(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// ReplaceStickerInSet replaces an existing sticker (referenced by file_id) in a
// set created by the bot with a new one, keeping its position.
func (b *Bot) ReplaceStickerInSet(ctx context.Context, oldSticker string, sticker InputSticker) error {
	doc, err := inputDocumentFromFileID(oldSticker)
	if err != nil {
		return err
	}

	item, err := b.resolveStickerItem(ctx, sticker)
	if err != nil {
		return err
	}

	if _, err := b.raw.StickersReplaceSticker(ctx, &tg.StickersReplaceStickerRequest{
		Sticker:    doc,
		NewSticker: item,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// GetCustomEmojiStickers returns information about custom emoji stickers by their
// identifiers (at most 200).
func (b *Bot) GetCustomEmojiStickers(ctx context.Context, customEmojiIDs []string) ([]Sticker, error) {
	if len(customEmojiIDs) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(customEmojiIDs))

	for _, s := range customEmojiIDs {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errInvalidCustomEmojiID()
		}

		ids = append(ids, id)
	}

	docs, err := b.raw.MessagesGetCustomEmojiDocuments(ctx, ids)
	if err != nil {
		return nil, asAPIError(err)
	}

	out := make([]Sticker, 0, len(docs))

	for _, d := range docs {
		if doc, ok := d.(*tg.Document); ok {
			out = append(out, stickerFromDocument(doc, "", StickerCustomEmoji))
		}
	}

	return out, nil
}
