package botapi

import (
	"context"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
)

// InputPaidMedia is a sealed union describing one item of paid media sent with
// SendPaidMedia.
//
// Concrete variants: InputPaidMediaPhoto, InputPaidMediaVideo.
type InputPaidMedia interface {
	isInputPaidMedia()
}

// InputPaidMediaPhoto is a photo paid-media item.
type InputPaidMediaPhoto struct {
	// Media is the photo, from a file_id, URL or local upload.
	Media InputFile
}

// InputPaidMediaVideo is a video paid-media item.
type InputPaidMediaVideo struct {
	// Media is the video, from a file_id, URL or local upload.
	Media InputFile
	// Thumbnail is an optional cover thumbnail, an uploaded file only.
	Thumbnail InputFile
	// Width, Height and Duration describe the video. Optional.
	Width    int
	Height   int
	Duration int
	// SupportsStreaming marks the video as suitable for streaming.
	SupportsStreaming bool
	// StartTimestamp is the timestamp in seconds from which the video plays.
	StartTimestamp int
}

func (InputPaidMediaPhoto) isInputPaidMedia() {}
func (InputPaidMediaVideo) isInputPaidMedia() {}

// uploadMediaRef finalizes a freshly built uploaded media via
// messages.uploadMedia and returns a reference to it usable as paid media.
func (b *Bot) uploadMediaRef(ctx context.Context, m tg.InputMediaClass) (tg.InputMediaClass, error) {
	res, err := b.raw.MessagesUploadMedia(ctx, &tg.MessagesUploadMediaRequest{
		Peer:  &tg.InputPeerSelf{},
		Media: m,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	switch r := res.(type) {
	case *tg.MessageMediaPhoto:
		photo, ok := r.Photo.(*tg.Photo)
		if !ok {
			return nil, &Error{Code: 500, Description: "Internal Server Error: uploaded media is not a photo"}
		}

		return &tg.InputMediaPhoto{ID: &tg.InputPhoto{
			ID:            photo.ID,
			AccessHash:    photo.AccessHash,
			FileReference: photo.FileReference,
		}}, nil
	case *tg.MessageMediaDocument:
		doc, ok := r.Document.(*tg.Document)
		if !ok {
			return nil, &Error{Code: 500, Description: "Internal Server Error: uploaded media is not a document"}
		}

		return &tg.InputMediaDocument{ID: &tg.InputDocument{
			ID:            doc.ID,
			AccessHash:    doc.AccessHash,
			FileReference: doc.FileReference,
		}}, nil
	default:
		return nil, &Error{Code: 500, Description: "Internal Server Error: unexpected upload media response"}
	}
}

// paidPhotoItem resolves a paid photo into an MTProto input media.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) paidPhotoItem(ctx context.Context, file InputFile) (tg.InputMediaClass, error) {
	switch f := file.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		return &tg.InputMediaPhoto{ID: &tg.InputPhoto{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}, nil
	case InputFileURL:
		return &tg.InputMediaPhotoExternal{URL: string(f)}, nil
	case *InputFileUpload:
		uploaded, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		return b.uploadMediaRef(ctx, &tg.InputMediaUploadedPhoto{File: uploaded})
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// paidVideoItem resolves a paid video into an MTProto input media.
//
// The switch over the sealed InputFile union is exhaustive.
func (b *Bot) paidVideoItem(ctx context.Context, v InputPaidMediaVideo) (tg.InputMediaClass, error) {
	switch f := v.Media.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		return &tg.InputMediaDocument{ID: &tg.InputDocument{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}, nil
	case InputFileURL:
		return &tg.InputMediaDocumentExternal{URL: string(f)}, nil
	case *InputFileUpload:
		return b.paidUploadedVideo(ctx, f, v)
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// paidUploadedVideo uploads a local video file and references it for paid media.
func (b *Bot) paidUploadedVideo(ctx context.Context, f *InputFileUpload, v InputPaidMediaVideo) (tg.InputMediaClass, error) {
	uploaded, err := b.uploadInputFile(ctx, f)
	if err != nil {
		return nil, err
	}

	video := &tg.DocumentAttributeVideo{
		SupportsStreaming: v.SupportsStreaming,
		Duration:          float64(v.Duration),
		W:                 v.Width,
		H:                 v.Height,
	}
	if v.StartTimestamp != 0 {
		video.VideoStartTs = float64(v.StartTimestamp)
	}

	doc := &tg.InputMediaUploadedDocument{
		File:       uploaded,
		MimeType:   mimeVideoMP4,
		Attributes: []tg.DocumentAttributeClass{video},
	}
	if f.Name != "" {
		doc.Attributes = append(doc.Attributes, &tg.DocumentAttributeFilename{FileName: f.Name})
	}

	if up, ok := v.Thumbnail.(*InputFileUpload); ok {
		thumb, err := b.uploadInputFile(ctx, up)
		if err != nil {
			return nil, err
		}

		doc.SetThumb(thumb)
	}

	return b.uploadMediaRef(ctx, doc)
}

// paidMediaItem resolves a single paid media item.
//
// The switch over the sealed InputPaidMedia union is exhaustive.
func (b *Bot) paidMediaItem(ctx context.Context, m InputPaidMedia) (tg.InputMediaClass, error) {
	switch m := m.(type) {
	case InputPaidMediaPhoto:
		return b.paidPhotoItem(ctx, m.Media)
	case InputPaidMediaVideo:
		return b.paidVideoItem(ctx, m)
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid paid media"}
	}
}

// SendPaidMedia sends paid media to a channel chat. Users pay starCount Telegram
// Stars to unlock the media. The media slice must contain 1-10 items.
func (b *Bot) SendPaidMedia(
	ctx context.Context, chat ChatID, starCount int, media []InputPaidMedia, caption string, opts ...SendOption,
) (*Message, error) {
	if len(media) < 1 || len(media) > 10 {
		return nil, &Error{Code: 400, Description: "Bad Request: paid media must include 1-10 items"}
	}

	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	return b.sendResolvedMedia(ctx, chat, caption, func(ctx context.Context, c []styling.StyledTextOption) (message.MediaOption, error) {
		extended := make([]tg.InputMediaClass, len(media))

		for i, m := range media {
			item, err := b.paidMediaItem(ctx, m)
			if err != nil {
				return nil, err
			}

			extended[i] = item
		}

		paid := &tg.InputMediaPaidMedia{
			StarsAmount:   int64(starCount),
			ExtendedMedia: extended,
		}
		if cfg.paidMediaPayload != "" {
			paid.SetPayload(cfg.paidMediaPayload)
		}

		return message.Media(paid, c...), nil
	}, opts...)
}
