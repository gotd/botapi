package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestStickerFromDocument(t *testing.T) {
	d := &tg.Document{
		ID:            5,
		AccessHash:    6,
		FileReference: []byte{1},
		DCID:          2,
		Size:          4096,
		MimeType:      mimeStickerVideo,
		Attributes: []tg.DocumentAttributeClass{
			&tg.DocumentAttributeImageSize{W: 512, H: 512},
			&tg.DocumentAttributeSticker{Alt: "🎉"},
		},
	}
	s := stickerFromDocument(d, "myset", StickerRegular)

	if s.FileID == "" || s.FileUniqueID == "" {
		t.Fatalf("file ids: %#v", s)
	}

	if s.Width != 512 || s.Height != 512 {
		t.Fatalf("dimensions: %#v", s)
	}

	if s.Emoji != "🎉" || s.SetName != "myset" || s.Type != StickerRegular {
		t.Fatalf("meta: %#v", s)
	}

	if !s.IsVideo || s.IsAnimated {
		t.Fatalf("video flags: %#v", s)
	}

	if s.FileSize != 4096 {
		t.Fatalf("size: %d", s.FileSize)
	}
}

func TestStickerSetType(t *testing.T) {
	if stickerSetType(tg.StickerSet{Masks: true}) != StickerMask {
		t.Fatal("masks -> mask")
	}

	if stickerSetType(tg.StickerSet{Emojis: true}) != StickerCustomEmoji {
		t.Fatal("emojis -> custom_emoji")
	}

	if stickerSetType(tg.StickerSet{}) != StickerRegular {
		t.Fatal("default -> regular")
	}
}
