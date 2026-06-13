package botapi

import (
	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// encodeFileID encodes a parsed file location into the Bot API file_id and
// file_unique_id strings. The unique-id derivation lives in file.go.
func encodeFileID(f fileid.FileID) (fileID, uniqueID string) {
	fileID, _ = fileid.EncodeFileID(f)
	return fileID, fileUniqueID(f)
}

// convertMessageMedia maps an incoming tg media attachment onto the Bot API
// Message media fields. Only the attachment types a bot can receive are mapped;
// web pages and unsupported kinds are ignored.
func convertMessageMedia(media tg.MessageMediaClass, r *Message) {
	switch media := media.(type) {
	case *tg.MessageMediaPhoto:
		if p, ok := media.Photo.(*tg.Photo); ok {
			r.Photo = photoSizesFromTg(p)
		}
	case *tg.MessageMediaDocument:
		if d, ok := media.Document.(*tg.Document); ok {
			setDocumentMedia(d, r)
		}
	case *tg.MessageMediaContact:
		r.Contact = &Contact{
			PhoneNumber: media.PhoneNumber,
			FirstName:   media.FirstName,
			LastName:    media.LastName,
			UserID:      media.UserID,
			VCard:       media.Vcard,
		}
	case *tg.MessageMediaPoll:
		r.Poll = pollFromTg(&media.Poll, &media.Results)
	}
}

// photoSizesFromTg converts a tg.Photo into the Bot API list of photo sizes,
// each carrying a usable file_id.
func photoSizesFromTg(photo *tg.Photo) []PhotoSize {
	type sized interface {
		GetW() int
		GetH() int
		GetType() string
	}

	var out []PhotoSize
	for _, sz := range photo.Sizes {
		size, ok := sz.(sized)
		if !ok {
			continue
		}
		t := size.GetType()
		if t == "" {
			continue
		}

		var fileSize int
		switch s := sz.(type) {
		case *tg.PhotoSize:
			fileSize = s.Size
		case *tg.PhotoCachedSize:
			fileSize = len(s.Bytes)
		}

		fileID, uniqueID := encodeFileID(fileid.FromPhoto(photo, rune(t[0])))
		out = append(out, PhotoSize{
			FileID:       fileID,
			FileUniqueID: uniqueID,
			Width:        size.GetW(),
			Height:       size.GetH(),
			FileSize:     fileSize,
		})
	}
	return out
}

// setDocumentMedia inspects a document's attributes and fills the matching Bot
// API media field (sticker, video, video note, voice, audio, animation, or a
// generic document).
func setDocumentMedia(d *tg.Document, r *Message) {
	fileID, uniqueID := encodeFileID(fileid.FromDocument(d))

	var (
		fileName        string
		width, height   int
		duration        int
		animated, round bool
	)
	for _, attr := range d.Attributes {
		switch a := attr.(type) {
		case *tg.DocumentAttributeFilename:
			fileName = a.FileName
		case *tg.DocumentAttributeImageSize:
			width, height = a.W, a.H
		case *tg.DocumentAttributeVideo:
			width, height, duration, round = a.W, a.H, int(a.Duration), a.RoundMessage
		case *tg.DocumentAttributeAnimated:
			animated = true
		}
	}

	// The attribute set determines the concrete Bot API type. Order matters:
	// stickers, then animations, then video/voice/audio, else a plain document.
	for _, attr := range d.Attributes {
		switch a := attr.(type) {
		case *tg.DocumentAttributeSticker:
			r.Sticker = &Sticker{
				FileID:       fileID,
				FileUniqueID: uniqueID,
				Type:         StickerRegular,
				Width:        width,
				Height:       height,
				Emoji:        a.Alt,
				FileSize:     int(d.Size),
			}
			return
		case *tg.DocumentAttributeAudio:
			if a.Voice {
				r.Voice = &Voice{
					FileID:       fileID,
					FileUniqueID: uniqueID,
					Duration:     a.Duration,
					MIMEType:     d.MimeType,
					FileSize:     d.Size,
				}
			} else {
				r.Audio = &Audio{
					FileID:       fileID,
					FileUniqueID: uniqueID,
					Duration:     a.Duration,
					Performer:    a.Performer,
					Title:        a.Title,
					FileName:     fileName,
					MIMEType:     d.MimeType,
					FileSize:     d.Size,
				}
			}
			return
		}
	}

	switch {
	case animated:
		r.Animation = &Animation{
			FileID:       fileID,
			FileUniqueID: uniqueID,
			Width:        width,
			Height:       height,
			Duration:     duration,
			FileName:     fileName,
			MIMEType:     d.MimeType,
			FileSize:     d.Size,
		}
	case round:
		r.VideoNote = &VideoNote{
			FileID:       fileID,
			FileUniqueID: uniqueID,
			Length:       width,
			Duration:     duration,
			FileSize:     int(d.Size),
		}
	case width != 0 || height != 0:
		r.Video = &Video{
			FileID:       fileID,
			FileUniqueID: uniqueID,
			Width:        width,
			Height:       height,
			Duration:     duration,
			FileName:     fileName,
			MIMEType:     d.MimeType,
			FileSize:     d.Size,
		}
	default:
		r.Document = &Document{
			FileID:       fileID,
			FileUniqueID: uniqueID,
			FileName:     fileName,
			MIMEType:     d.MimeType,
			FileSize:     d.Size,
		}
	}
}
