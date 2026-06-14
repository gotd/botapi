package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// EditMessageMedia replaces the media (and caption) of a message. The new media
// may reference an existing file_id, an HTTP URL Telegram fetches, or a local
// upload.
func (b *Bot) EditMessageMedia(ctx context.Context, chat ChatID, messageID int, media InputMedia, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	tgMedia, capt, err := b.inputMediaToTg(ctx, media)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesEditMessageRequest{Peer: peer, ID: messageID}
	req.SetMedia(tgMedia)

	msg, entities, err := b.caption(ctx, capt)
	if err != nil {
		return nil, err
	}

	if msg != "" {
		req.Message = msg
	}

	if len(entities) > 0 {
		req.Entities = entities
	}

	if cfg.markup != nil {
		mkp, err := replyMarkupToTg(cfg.markup)
		if err != nil {
			return nil, err
		}

		req.SetReplyMarkup(mkp)
	}

	resp, err := b.raw.MessagesEditMessage(ctx, req)

	return b.sentMessage(ctx, peer, resp, err)
}

// captionData is the caption text and how to format it, extracted from an
// InputMedia.
type captionData struct {
	text      string
	parseMode ParseMode
	entities  []MessageEntity
}

// caption resolves a caption into a message string and MTProto entities,
// honoring explicit entities or a parse mode.
func (b *Bot) caption(ctx context.Context, c captionData) (string, []tg.MessageEntityClass, error) {
	if len(c.entities) > 0 {
		return c.text, entitiesToTg(c.entities), nil
	}

	if c.parseMode != ParseModeNone && c.text != "" {
		return b.styledMessage(ctx, c.text, c.parseMode)
	}

	return c.text, nil, nil
}

// inputMediaToTg converts a Bot API InputMedia into the MTProto media and its
// caption. The switch over the sealed InputMedia union is exhaustive.
func (b *Bot) inputMediaToTg(ctx context.Context, m InputMedia) (tg.InputMediaClass, captionData, error) {
	switch m := m.(type) {
	case *InputMediaPhoto:
		media, err := b.photoInputMedia(ctx, m.Media)
		return media, captionData{m.Caption, m.ParseMode, m.CaptionEntities}, err
	case *InputMediaVideo:
		media, err := b.documentInputMedia(ctx, m.Media, "video/mp4")
		return media, captionData{m.Caption, m.ParseMode, m.CaptionEntities}, err
	case *InputMediaAnimation:
		media, err := b.documentInputMedia(ctx, m.Media, "video/mp4")
		return media, captionData{m.Caption, m.ParseMode, m.CaptionEntities}, err
	case *InputMediaAudio:
		media, err := b.documentInputMedia(ctx, m.Media, "audio/mpeg")
		return media, captionData{m.Caption, m.ParseMode, m.CaptionEntities}, err
	case *InputMediaDocument:
		media, err := b.documentInputMedia(ctx, m.Media, "application/octet-stream")
		return media, captionData{m.Caption, m.ParseMode, m.CaptionEntities}, err
	default:
		return nil, captionData{}, &Error{Code: 400, Description: "Bad Request: unsupported input media type"}
	}
}

// photoInputMedia resolves an InputFile into MTProto photo media.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) photoInputMedia(ctx context.Context, file InputFile) (tg.InputMediaClass, error) {
	switch f := file.(type) {
	case InputFileID:
		photo, err := inputPhotoFromFileID(string(f))
		if err != nil {
			return nil, err
		}

		return &tg.InputMediaPhoto{ID: photo}, nil
	case InputFileURL:
		return &tg.InputMediaPhotoExternal{URL: string(f)}, nil
	case *InputFileUpload:
		uploaded, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		return &tg.InputMediaUploadedPhoto{File: uploaded}, nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// documentInputMedia resolves an InputFile into MTProto document media.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) documentInputMedia(ctx context.Context, file InputFile, mimeType string) (tg.InputMediaClass, error) {
	switch f := file.(type) {
	case InputFileID:
		doc, err := inputDocumentFromFileID(string(f))
		if err != nil {
			return nil, err
		}

		return &tg.InputMediaDocument{ID: doc}, nil
	case InputFileURL:
		return &tg.InputMediaDocumentExternal{URL: string(f)}, nil
	case *InputFileUpload:
		uploaded, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		media := &tg.InputMediaUploadedDocument{File: uploaded, MimeType: mimeType}
		if f.Name != "" {
			media.Attributes = []tg.DocumentAttributeClass{&tg.DocumentAttributeFilename{FileName: f.Name}}
		}

		return media, nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}
