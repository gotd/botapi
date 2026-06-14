package botapi

import (
	"context"
	"io"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
)

// uploadInputFile uploads a local InputFileUpload and returns the MTProto input
// file handle.
func (b *Bot) uploadInputFile(ctx context.Context, f *InputFileUpload) (tg.InputFileClass, error) {
	up := uploader.NewUploader(b.raw)
	name := f.Name

	if name == "" {
		name = "file"
	}

	switch {
	case f.Path != "":
		file, err := up.FromPath(ctx, f.Path)
		return file, asAPIError(err)
	case f.Bytes != nil:
		file, err := up.FromBytes(ctx, name, f.Bytes)
		return file, asAPIError(err)
	case f.Reader != nil:
		file, err := up.FromReader(ctx, name, io.Reader(f.Reader))
		return file, asAPIError(err)
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: empty file"}
	}
}

// photoMedia resolves an InputFile into a photo media option.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) photoMedia(ctx context.Context, file InputFile, caption []styling.StyledTextOption) (message.MediaOption, error) {
	switch f := file.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		photo := &tg.InputMediaPhoto{ID: &tg.InputPhoto{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}

		return message.Media(photo, caption...), nil
	case InputFileURL:
		return message.PhotoExternal(string(f), caption...), nil
	case *InputFileUpload:
		upFile, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		return message.UploadedPhoto(upFile, caption...), nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// documentMedia resolves an InputFile into a general document media option.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) documentMedia(ctx context.Context, file InputFile, caption []styling.StyledTextOption) (message.MediaOption, error) {
	switch f := file.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		doc := &tg.InputMediaDocument{ID: &tg.InputDocument{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}

		return message.Media(doc, caption...), nil
	case InputFileURL:
		return message.DocumentExternal(string(f), caption...), nil
	case *InputFileUpload:
		upFile, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		b := message.UploadedDocument(upFile, caption...)
		if f.Name != "" {
			b = b.Filename(f.Name)
		}

		return b, nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// sendResolvedMedia styles the caption, resolves the media and sends it.
func (b *Bot) sendResolvedMedia(
	ctx context.Context,
	chat ChatID,
	caption string,
	build func(ctx context.Context, caption []styling.StyledTextOption) (message.MediaOption, error),
	opts ...SendOption,
) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	styled, err := styledText(caption, cfg.parseMode, b.peers.UserResolveHook(ctx))
	if err != nil {
		return nil, err
	}

	media, err := build(ctx, styled)
	if err != nil {
		return nil, err
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

	resp, err := builder.Media(ctx, media)

	return b.sentMessage(ctx, peer, resp, err)
}

// SendPhoto sends a photo from a file_id, URL or local upload.
func (b *Bot) SendPhoto(ctx context.Context, chat ChatID, photo InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.photoMedia(ctx, photo, c)
	}, opts...)
}

// SendDocument sends a general file from a file_id, URL or local upload.
func (b *Bot) SendDocument(ctx context.Context, chat ChatID, document InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.documentMedia(ctx, document, c)
	}, opts...)
}
