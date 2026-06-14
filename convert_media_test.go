package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestConvertMessageMedia_Photo(t *testing.T) {
	photo := &tg.Photo{
		ID:            10,
		AccessHash:    20,
		FileReference: []byte{1, 2, 3},
		DCID:          2,
		Sizes: []tg.PhotoSizeClass{
			&tg.PhotoSize{Type: "x", W: 800, H: 600, Size: 1234},
		},
	}

	var r Message

	convertMessageMedia(&tg.MessageMediaPhoto{Photo: photo}, &r)

	if len(r.Photo) != 1 {
		t.Fatalf("want 1 photo size, got %d", len(r.Photo))
	}

	got := r.Photo[0]
	if got.Width != 800 || got.Height != 600 || got.FileSize != 1234 {
		t.Fatalf("dimensions: %#v", got)
	}

	if got.FileID == "" || got.FileUniqueID == "" {
		t.Fatalf("file ids not set: %#v", got)
	}
}

func docWith(attrs ...tg.DocumentAttributeClass) *tg.MessageMediaDocument {
	return &tg.MessageMediaDocument{Document: &tg.Document{
		ID:            1,
		AccessHash:    2,
		FileReference: []byte{9},
		DCID:          2,
		Size:          5000,
		MimeType:      "application/octet-stream",
		Attributes:    attrs,
	}}
}

func TestConvertMessageMedia_Document(t *testing.T) {
	var r Message

	convertMessageMedia(docWith(&tg.DocumentAttributeFilename{FileName: "report.pdf"}), &r)

	if r.Document == nil {
		t.Fatal("document not set")
	}

	if r.Document.FileName != "report.pdf" || r.Document.FileSize != 5000 || r.Document.FileID == "" {
		t.Fatalf("document: %#v", r.Document)
	}
}

func TestConvertMessageMedia_Video(t *testing.T) {
	var r Message

	convertMessageMedia(docWith(&tg.DocumentAttributeVideo{W: 1280, H: 720, Duration: 30}), &r)

	if r.Video == nil {
		t.Fatalf("video not set: %#v", r)
	}

	if r.Video.Width != 1280 || r.Video.Height != 720 || r.Video.Duration != 30 {
		t.Fatalf("video: %#v", r.Video)
	}
}

func TestConvertMessageMedia_Voice(t *testing.T) {
	var r Message

	convertMessageMedia(docWith(&tg.DocumentAttributeAudio{Voice: true, Duration: 7}), &r)

	if r.Voice == nil || r.Voice.Duration != 7 {
		t.Fatalf("voice: %#v", r.Voice)
	}

	if r.Audio != nil {
		t.Fatal("voice must not also set Audio")
	}
}

func TestConvertMessageMedia_Sticker(t *testing.T) {
	var r Message

	convertMessageMedia(docWith(
		&tg.DocumentAttributeSticker{Alt: "😀"},
		&tg.DocumentAttributeImageSize{W: 512, H: 512},
	), &r)

	if r.Sticker == nil || r.Sticker.Emoji != "😀" || r.Sticker.Width != 512 {
		t.Fatalf("sticker: %#v", r.Sticker)
	}
}

func TestConvertMessageMedia_Contact(t *testing.T) {
	var r Message

	convertMessageMedia(&tg.MessageMediaContact{
		PhoneNumber: "+1",
		FirstName:   "Ada",
		Vcard:       "BEGIN:VCARD",
	}, &r)

	if r.Contact == nil || r.Contact.PhoneNumber != "+1" || r.Contact.VCard != "BEGIN:VCARD" {
		t.Fatalf("contact: %#v", r.Contact)
	}
}
