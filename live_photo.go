package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// liveVideoDocument resolves the video component of a live photo into an MTProto
// input document. An upload is finalized through messages.uploadMedia; a file_id
// references an existing document. URLs are not supported.
func (b *Bot) liveVideoDocument(ctx context.Context, file InputFile) (tg.InputDocumentClass, error) {
	switch f := file.(type) {
	case InputFileID:
		return inputDocumentFromFileID(string(f))
	case *InputFileUpload:
		uploaded, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		media, err := b.uploadMediaRef(ctx, &tg.InputMediaUploadedDocument{
			File:       uploaded,
			MimeType:   mimeVideoMP4,
			Attributes: []tg.DocumentAttributeClass{&tg.DocumentAttributeVideo{}},
		})
		if err != nil {
			return nil, err
		}

		doc, ok := media.(*tg.InputMediaDocument)
		if !ok {
			return nil, &Error{Code: 500, Description: "Internal Server Error: unexpected live photo video media"}
		}

		return doc.ID, nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: live photo video must be an uploaded file or file_id"}
	}
}

// SendLivePhoto sends a live photo (an Apple-style still image paired with a short
// video). The still photo must be an uploaded file; livePhoto is its video, an
// uploaded file or a file_id.
func (b *Bot) SendLivePhoto(
	ctx context.Context, chat ChatID, photo, livePhoto InputFile, caption string, opts ...SendOption,
) (*Message, error) {
	up, ok := photo.(*InputFileUpload)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: live photo still must be an uploaded file"}
	}

	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		photoFile, err := b.uploadInputFile(ctx, up)
		if err != nil {
			return nil, err
		}

		video, err := b.liveVideoDocument(ctx, livePhoto)
		if err != nil {
			return nil, err
		}

		media := &tg.InputMediaUploadedPhoto{
			File:      photoFile,
			Video:     video,
			LivePhoto: true,
		}

		return message.Media(media, c...), nil
	}, opts...)
}
