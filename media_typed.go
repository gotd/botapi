package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
)

// typedMedia resolves an InputFile into a typed document media option. For local
// uploads it applies the typed attribute builder (video/audio/...); for file_id
// and URL inputs it falls back to a plain document, since the stored document
// already carries its type.
func (b *Bot) typedMedia(
	ctx context.Context,
	file InputFile,
	caption []styling.StyledTextOption,
	typed func(*message.UploadedDocumentBuilder) message.MediaOption,
) (message.MediaOption, error) {
	up, ok := file.(*InputFileUpload)
	if !ok {
		return b.documentMedia(ctx, file, caption)
	}
	f, err := b.uploadInputFile(ctx, up)
	if err != nil {
		return nil, err
	}
	doc := message.UploadedDocument(f, caption...)
	if up.Name != "" {
		doc = doc.Filename(up.Name)
	}
	return typed(doc), nil
}

// SendVideo sends a video from a file_id, URL or local upload.
func (b *Bot) SendVideo(ctx context.Context, chat ChatID, video InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, video, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.Video() })
	}, opts...)
}

// SendAnimation sends an animation (GIF or silent video) from a file_id, URL or
// local upload.
func (b *Bot) SendAnimation(ctx context.Context, chat ChatID, animation InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, animation, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.GIF() })
	}, opts...)
}

// SendAudio sends an audio file (music) from a file_id, URL or local upload.
func (b *Bot) SendAudio(ctx context.Context, chat ChatID, audio InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, audio, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.Audio() })
	}, opts...)
}

// SendVoice sends a voice note from a file_id, URL or local upload.
func (b *Bot) SendVoice(ctx context.Context, chat ChatID, voice InputFile, caption string, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, voice, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.Voice() })
	}, opts...)
}

// SendVideoNote sends a rounded square video message from a file_id or local
// upload.
func (b *Bot) SendVideoNote(ctx context.Context, chat ChatID, videoNote InputFile, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, "", func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, videoNote, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.RoundVideo() })
	}, opts...)
}

// SendSticker sends a sticker from a file_id, URL or local upload.
func (b *Bot) SendSticker(ctx context.Context, chat ChatID, sticker InputFile, opts ...SendOption) (*Message, error) {
	return b.sendResolvedMedia(ctx, chat, "", func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		return b.typedMedia(ctx, sticker, c, func(d *message.UploadedDocumentBuilder) message.MediaOption { return d.UploadedSticker() })
	}, opts...)
}
