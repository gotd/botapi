package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// errNonUploadInAlbum is returned for album items that reference a file_id or
// URL. The high-level send API can only place freshly uploaded files in a group.
var errNonUploadInAlbum = &Error{
	Code:        400,
	Description: "Bad Request: media group items must be uploaded files (file_id and URL are not supported yet)",
}

// photoMulti resolves a photo InputFile into an album item. Only uploads are
// supported in groups.
func (b *Bot) photoMulti(ctx context.Context, file InputFile, caption []styling.StyledTextOption) (message.MultiMediaOption, error) {
	up, ok := file.(*InputFileUpload)
	if !ok {
		return nil, errNonUploadInAlbum
	}

	upFile, err := b.uploadInputFile(ctx, up)
	if err != nil {
		return nil, err
	}

	return message.UploadedPhoto(upFile, caption...), nil
}

// docMulti resolves a document-family InputFile into an album item, applying the
// typed attribute builder. Only uploads are supported in groups.
func (b *Bot) docMulti(
	ctx context.Context,
	file InputFile,
	caption []styling.StyledTextOption,
	typed func(*message.UploadedDocumentBuilder) message.MultiMediaOption,
) (message.MultiMediaOption, error) {
	up, ok := file.(*InputFileUpload)
	if !ok {
		return nil, errNonUploadInAlbum
	}

	upFile, err := b.uploadInputFile(ctx, up)
	if err != nil {
		return nil, err
	}

	doc := message.UploadedDocument(upFile, caption...)
	if up.Name != "" {
		doc = doc.Filename(up.Name)
	}

	return typed(doc), nil
}

// inputMediaToMulti resolves a single InputMedia into an album item.
//
// The switch over the sealed InputMedia union is exhaustive.
func (b *Bot) inputMediaToMulti(ctx context.Context, m InputMedia) (message.MultiMediaOption, error) {
	resolver := b.peers.UserResolveHook(ctx)
	plain := func(d *message.UploadedDocumentBuilder) message.MultiMediaOption { return d }

	switch m := m.(type) {
	case *InputMediaPhoto:
		caption, err := styledText(m.Caption, m.ParseMode, resolver)
		if err != nil {
			return nil, err
		}

		return b.photoMulti(ctx, m.Media, caption)
	case *InputMediaVideo:
		caption, err := styledText(m.Caption, m.ParseMode, resolver)
		if err != nil {
			return nil, err
		}

		return b.docMulti(ctx, m.Media, caption, func(d *message.UploadedDocumentBuilder) message.MultiMediaOption { return d.Video() })
	case *InputMediaAnimation:
		caption, err := styledText(m.Caption, m.ParseMode, resolver)
		if err != nil {
			return nil, err
		}

		return b.docMulti(ctx, m.Media, caption, func(d *message.UploadedDocumentBuilder) message.MultiMediaOption { return d.GIF() })
	case *InputMediaAudio:
		caption, err := styledText(m.Caption, m.ParseMode, resolver)
		if err != nil {
			return nil, err
		}

		return b.docMulti(ctx, m.Media, caption, func(d *message.UploadedDocumentBuilder) message.MultiMediaOption { return d.Audio() })
	case *InputMediaDocument:
		caption, err := styledText(m.Caption, m.ParseMode, resolver)
		if err != nil {
			return nil, err
		}

		return b.docMulti(ctx, m.Media, caption, plain)
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid media"}
	}
}

// sentMessages extracts every message produced by a send into Bot API messages.
func (b *Bot) sentMessages(ctx context.Context, resp tg.UpdatesClass, sendErr error) ([]*Message, error) {
	if sendErr != nil {
		return nil, asAPIError(sendErr)
	}

	var updates []tg.UpdateClass

	switch u := resp.(type) {
	case *tg.Updates:
		updates = u.Updates
	case *tg.UpdatesCombined:
		updates = u.Updates
	default:
		return nil, nil
	}

	var out []*Message

	for _, upd := range updates {
		var msg tg.MessageClass

		switch u := upd.(type) {
		case *tg.UpdateNewMessage:
			msg = u.Message
		case *tg.UpdateNewChannelMessage:
			msg = u.Message
		default:
			continue
		}

		m, ok := msg.(*tg.Message)
		if !ok {
			continue
		}

		converted, err := b.convertMessage(ctx, m)
		if err != nil {
			return nil, err
		}

		out = append(out, converted)
	}

	return out, nil
}

// SendMediaGroup sends a group of 2-10 photos, videos, documents or audio files
// as an album. Items referencing an existing file_id are not yet supported.
func (b *Bot) SendMediaGroup(ctx context.Context, chat ChatID, media []InputMedia, opts ...SendOption) ([]*Message, error) {
	if len(media) < 2 || len(media) > 10 {
		return nil, &Error{Code: 400, Description: "Bad Request: media group must include 2-10 items"}
	}

	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	items := make([]message.MultiMediaOption, len(media))
	for i, m := range media {
		mm, err := b.inputMediaToMulti(ctx, m)
		if err != nil {
			return nil, err
		}

		items[i] = mm
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	builder := &b.sender.To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.Album(ctx, items[0], items[1:]...)

	return b.sentMessages(ctx, resp, err)
}
