package botapi

import (
	"context"
	"strings"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// StickerFormat is the format of a sticker file.
type StickerFormat string

const (
	// StickerFormatStatic is a static .WEBP or .PNG sticker.
	StickerFormatStatic StickerFormat = "static"
	// StickerFormatAnimated is an animated .TGS sticker.
	StickerFormatAnimated StickerFormat = "animated"
	// StickerFormatVideo is a video .WEBM sticker.
	StickerFormatVideo StickerFormat = "video"
)

// mimeStickerStatic is the MIME type used for static (and unknown) stickers.
const mimeStickerStatic = "image/png"

// mimeType reports the MTProto MIME type for the sticker format. Static (and
// any unknown) formats default to PNG.
func (f StickerFormat) mimeType() string {
	switch f {
	case StickerFormatAnimated:
		return "application/x-tgsticker"
	case StickerFormatVideo:
		return "video/webm"
	case StickerFormatStatic:
		return mimeStickerStatic
	default:
		return mimeStickerStatic
	}
}

// InputSticker describes a sticker to be added to a set.
type InputSticker struct {
	// Sticker is the sticker file: a file_id from UploadStickerFile, or a local
	// upload. URLs are not accepted.
	Sticker InputFile
	// Format is the sticker file format.
	Format StickerFormat
	// EmojiList are the emoji associated with the sticker (1-20).
	EmojiList []string
	// Keywords are optional search keywords, comma-joined on the wire.
	Keywords []string
}

// UploadStickerFile uploads a file for later use in a sticker set and returns a
// File whose file_id can be passed to CreateNewStickerSet or AddStickerToSet.
func (b *Bot) UploadStickerFile(ctx context.Context, userID int64, sticker InputFile, format StickerFormat) (*File, error) {
	up, ok := sticker.(*InputFileUpload)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: sticker file must be an uploaded file"}
	}
	doc, err := b.uploadStickerDocument(ctx, up, format)
	if err != nil {
		return nil, err
	}
	encoded, err := fileid.EncodeFileID(fileid.FromDocument(doc))
	if err != nil {
		return nil, asAPIError(err)
	}
	return &File{
		FileID:       encoded,
		FileUniqueID: fileUniqueID(fileid.FromDocument(doc)),
		FileSize:     doc.Size,
	}, nil
}

// uploadStickerDocument uploads a local file and turns it into a saved document
// via messages.uploadMedia, returning the resulting document.
func (b *Bot) uploadStickerDocument(ctx context.Context, up *InputFileUpload, format StickerFormat) (*tg.Document, error) {
	uploaded, err := b.uploadInputFile(ctx, up)
	if err != nil {
		return nil, err
	}
	name := up.Name
	if name == "" {
		name = "sticker"
	}
	media, err := b.raw.MessagesUploadMedia(ctx, &tg.MessagesUploadMediaRequest{
		Peer: &tg.InputPeerSelf{},
		Media: &tg.InputMediaUploadedDocument{
			File:     uploaded,
			MimeType: format.mimeType(),
			Attributes: []tg.DocumentAttributeClass{
				&tg.DocumentAttributeFilename{FileName: name},
			},
		},
	})
	if err != nil {
		return nil, asAPIError(err)
	}
	mediaDoc, ok := media.(*tg.MessageMediaDocument)
	if !ok {
		return nil, &Error{Code: 500, Description: "Internal Server Error: unexpected upload media response"}
	}
	doc, ok := mediaDoc.Document.(*tg.Document)
	if !ok {
		return nil, &Error{Code: 500, Description: "Internal Server Error: uploaded media is not a document"}
	}
	return doc, nil
}

// resolveStickerItem converts an InputSticker into the MTProto set item,
// uploading the file when it is a local upload.
func (b *Bot) resolveStickerItem(ctx context.Context, sticker InputSticker) (tg.InputStickerSetItem, error) {
	var doc tg.InputDocumentClass
	switch f := sticker.Sticker.(type) {
	case InputFileID:
		d, err := inputDocumentFromFileID(string(f))
		if err != nil {
			return tg.InputStickerSetItem{}, err
		}
		doc = d
	case *InputFileUpload:
		uploaded, err := b.uploadStickerDocument(ctx, f, sticker.Format)
		if err != nil {
			return tg.InputStickerSetItem{}, err
		}
		doc = &tg.InputDocument{ID: uploaded.ID, AccessHash: uploaded.AccessHash, FileReference: uploaded.FileReference}
	default:
		return tg.InputStickerSetItem{}, &Error{Code: 400, Description: "Bad Request: sticker must be a file_id or an uploaded file"}
	}

	item := tg.InputStickerSetItem{
		Document: doc,
		Emoji:    strings.Join(sticker.EmojiList, ""),
	}
	if len(sticker.Keywords) > 0 {
		item.SetKeywords(strings.Join(sticker.Keywords, ","))
	}
	return item, nil
}

// CreateNewStickerSet creates a new sticker set owned by the given user. name is
// the set short name (used in t.me/addstickers/<name>).
func (b *Bot) CreateNewStickerSet(
	ctx context.Context, userID int64, name, title string, stickers []InputSticker, opts ...StickerSetOption,
) error {
	var cfg stickerSetConfig
	for _, o := range opts {
		o(&cfg)
	}
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}
	items := make([]tg.InputStickerSetItem, 0, len(stickers))
	for _, s := range stickers {
		item, err := b.resolveStickerItem(ctx, s)
		if err != nil {
			return err
		}
		items = append(items, item)
	}
	if _, err := b.raw.StickersCreateStickerSet(ctx, &tg.StickersCreateStickerSetRequest{
		Masks:     cfg.masks,
		UserID:    user,
		Title:     title,
		ShortName: name,
		Stickers:  items,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// AddStickerToSet adds a sticker to a set created by the bot.
func (b *Bot) AddStickerToSet(ctx context.Context, name string, sticker InputSticker) error {
	item, err := b.resolveStickerItem(ctx, sticker)
	if err != nil {
		return err
	}
	if _, err := b.raw.StickersAddStickerToSet(ctx, &tg.StickersAddStickerToSetRequest{
		Stickerset: &tg.InputStickerSetShortName{ShortName: name},
		Sticker:    item,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// DeleteStickerFromSet removes a sticker, referenced by file_id, from the set it
// belongs to.
func (b *Bot) DeleteStickerFromSet(ctx context.Context, sticker string) error {
	doc, err := inputDocumentFromFileID(sticker)
	if err != nil {
		return err
	}
	if _, err := b.raw.StickersRemoveStickerFromSet(ctx, doc); err != nil {
		return asAPIError(err)
	}
	return nil
}

// SetStickerPositionInSet moves a sticker, referenced by file_id, to the given
// zero-based position in its set.
func (b *Bot) SetStickerPositionInSet(ctx context.Context, sticker string, position int) error {
	doc, err := inputDocumentFromFileID(sticker)
	if err != nil {
		return err
	}
	if _, err := b.raw.StickersChangeStickerPosition(ctx, &tg.StickersChangeStickerPositionRequest{
		Sticker:  doc,
		Position: position,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// StickerSetOption configures CreateNewStickerSet.
type StickerSetOption func(*stickerSetConfig)

type stickerSetConfig struct {
	masks bool
}

// WithMaskStickers marks the new set as a set of mask stickers.
func WithMaskStickers() StickerSetOption {
	return func(c *stickerSetConfig) { c.masks = true }
}
