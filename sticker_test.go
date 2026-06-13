package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestStickerFormatMime(t *testing.T) {
	cases := map[StickerFormat]string{
		StickerFormatStatic:   mimeStickerStatic,
		StickerFormatAnimated: "application/x-tgsticker",
		StickerFormatVideo:    "video/webm",
		StickerFormat("xxx"):  mimeStickerStatic,
	}
	for f, want := range cases {
		if got := f.mimeType(); got != want {
			t.Fatalf("%s: got %q, want %q", f, got, want)
		}
	}
}

func TestResolveStickerItemFromFileID(t *testing.T) {
	b := &Bot{}
	fid := documentFileID(t, 0x99)
	item, err := b.resolveStickerItem(context.Background(), InputSticker{
		Sticker:   InputFileID(fid),
		EmojiList: []string{"😀", "🎉"},
		Keywords:  []string{"party", "fun"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if d, ok := item.Document.(*tg.InputDocument); !ok || d.ID != 0x99 {
		t.Fatalf("document: %#v", item.Document)
	}
	if item.Emoji != "😀🎉" {
		t.Fatalf("emoji: %q", item.Emoji)
	}
	if kw, ok := item.GetKeywords(); !ok || kw != "party,fun" {
		t.Fatalf("keywords: %q %v", kw, ok)
	}
}

func TestResolveStickerItemRejectsURL(t *testing.T) {
	b := &Bot{}
	_, err := b.resolveStickerItem(context.Background(), InputSticker{Sticker: InputFileURL("https://x/y.png")})
	if err == nil {
		t.Fatal("expected error for URL sticker")
	}
}
