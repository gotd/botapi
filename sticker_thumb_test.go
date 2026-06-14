package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetForumTopicIconStickers(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesGetStickerSetRequestTypeID, &tg.MessagesStickerSet{
		Set: tg.StickerSet{ShortName: "topics", Emojis: true},
		Documents: []tg.DocumentClass{
			&tg.Document{
				ID:       1,
				MimeType: "image/webp",
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeSticker{Alt: "🎯", Stickerset: &tg.InputStickerSetEmpty{}},
				},
			},
		},
	})

	stickers, err := newMockBot(inv).GetForumTopicIconStickers(context.Background())
	if err != nil {
		t.Fatalf("GetForumTopicIconStickers: %v", err)
	}

	if len(stickers) != 1 || stickers[0].Type != StickerCustomEmoji {
		t.Fatalf("stickers = %#v", stickers)
	}
}

func TestSetStickerSetThumbnail(t *testing.T) {
	fid := documentFileID(t, 0x61)

	inv := newMockInvoker()
	inv.reply(tg.StickersSetStickerSetThumbRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	if err := newMockBot(inv).SetStickerSetThumbnail(context.Background(), "s", FileID(fid), StickerFormatStatic); err != nil {
		t.Fatalf("SetStickerSetThumbnail: %v", err)
	}

	var req tg.StickersSetStickerSetThumbRequest

	inv.decode(t, tg.StickersSetStickerSetThumbRequestTypeID, &req)

	if _, ok := req.GetThumb(); !ok {
		t.Fatal("thumb not set")
	}
}

func TestSetCustomEmojiStickerSetThumbnail(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.StickersSetStickerSetThumbRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	if err := newMockBot(inv).SetCustomEmojiStickerSetThumbnail(context.Background(), "s", "777"); err != nil {
		t.Fatalf("SetCustomEmojiStickerSetThumbnail: %v", err)
	}

	var req tg.StickersSetStickerSetThumbRequest

	inv.decode(t, tg.StickersSetStickerSetThumbRequestTypeID, &req)

	if id, ok := req.GetThumbDocumentID(); !ok || id != 777 {
		t.Fatalf("thumb document id = %d, ok=%v", id, ok)
	}
}
