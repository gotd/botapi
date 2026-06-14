package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestDeleteStickerSet(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.StickersDeleteStickerSetRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).DeleteStickerSet(context.Background(), "s"); err != nil {
		t.Fatalf("DeleteStickerSet: %v", err)
	}

	if !inv.called(tg.StickersDeleteStickerSetRequestTypeID) {
		t.Fatal("delete not called")
	}
}

func TestSetStickerSetTitle(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.StickersRenameStickerSetRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	if err := newMockBot(inv).SetStickerSetTitle(context.Background(), "s", "New Title"); err != nil {
		t.Fatalf("SetStickerSetTitle: %v", err)
	}

	var req tg.StickersRenameStickerSetRequest

	inv.decode(t, tg.StickersRenameStickerSetRequestTypeID, &req)

	if req.Title != "New Title" {
		t.Fatalf("title = %q", req.Title)
	}
}

func TestSetStickerEmojiList(t *testing.T) {
	fid := documentFileID(t, 0x55)

	inv := newMockInvoker()
	inv.reply(tg.StickersChangeStickerRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	if err := newMockBot(inv).SetStickerEmojiList(context.Background(), fid, []string{"😀", "😎"}); err != nil {
		t.Fatalf("SetStickerEmojiList: %v", err)
	}

	var req tg.StickersChangeStickerRequest

	inv.decode(t, tg.StickersChangeStickerRequestTypeID, &req)

	if got, ok := req.GetEmoji(); !ok || got != "😀😎" {
		t.Fatalf("emoji = %q, ok=%v", got, ok)
	}
}

func TestSetStickerKeywords(t *testing.T) {
	fid := documentFileID(t, 0x56)

	inv := newMockInvoker()
	inv.reply(tg.StickersChangeStickerRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	if err := newMockBot(inv).SetStickerKeywords(context.Background(), fid, []string{"a", "b"}); err != nil {
		t.Fatalf("SetStickerKeywords: %v", err)
	}

	var req tg.StickersChangeStickerRequest

	inv.decode(t, tg.StickersChangeStickerRequestTypeID, &req)

	if got, ok := req.GetKeywords(); !ok || got != "a,b" {
		t.Fatalf("keywords = %q, ok=%v", got, ok)
	}
}

func TestSetStickerMaskPosition(t *testing.T) {
	fid := documentFileID(t, 0x57)

	inv := newMockInvoker()
	inv.reply(tg.StickersChangeStickerRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	pos := &MaskPosition{Point: "eyes", XShift: 0.1, YShift: -0.2, Scale: 2.0}
	if err := newMockBot(inv).SetStickerMaskPosition(context.Background(), fid, pos); err != nil {
		t.Fatalf("SetStickerMaskPosition: %v", err)
	}

	var req tg.StickersChangeStickerRequest

	inv.decode(t, tg.StickersChangeStickerRequestTypeID, &req)

	coords, ok := req.GetMaskCoords()
	if !ok || coords.N != 1 || coords.Zoom != 2.0 {
		t.Fatalf("coords = %#v, ok=%v", coords, ok)
	}
}

func TestSetStickerMaskPositionInvalidPoint(t *testing.T) {
	fid := documentFileID(t, 0x58)

	err := newMockBot(newMockInvoker()).SetStickerMaskPosition(context.Background(), fid, &MaskPosition{Point: "nose"})
	if err == nil {
		t.Fatal("expected error for invalid point")
	}
}

func TestReplaceStickerInSet(t *testing.T) {
	fid := documentFileID(t, 0x59)
	newFid := documentFileID(t, 0x5a)

	inv := newMockInvoker()
	inv.reply(tg.StickersReplaceStickerRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

	sticker := InputSticker{Sticker: FileID(newFid), Format: StickerFormatStatic, EmojiList: []string{"😀"}}
	if err := newMockBot(inv).ReplaceStickerInSet(context.Background(), fid, sticker); err != nil {
		t.Fatalf("ReplaceStickerInSet: %v", err)
	}

	var req tg.StickersReplaceStickerRequest

	inv.decode(t, tg.StickersReplaceStickerRequestTypeID, &req)

	if req.NewSticker.Emoji != "😀" {
		t.Fatalf("new sticker = %#v", req.NewSticker)
	}
}

func TestGetCustomEmojiStickers(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesGetCustomEmojiDocumentsRequestTypeID, &tg.DocumentClassVector{
		Elems: []tg.DocumentClass{
			&tg.Document{
				ID:         1,
				AccessHash: 2,
				MimeType:   "image/webp",
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeSticker{Alt: "😀", Stickerset: &tg.InputStickerSetEmpty{}},
				},
			},
		},
	})

	stickers, err := newMockBot(inv).GetCustomEmojiStickers(context.Background(), []string{"12345"})
	if err != nil {
		t.Fatalf("GetCustomEmojiStickers: %v", err)
	}

	if len(stickers) != 1 || stickers[0].Type != StickerCustomEmoji || stickers[0].Emoji != "😀" {
		t.Fatalf("stickers = %#v", stickers)
	}

	var req tg.MessagesGetCustomEmojiDocumentsRequest

	inv.decode(t, tg.MessagesGetCustomEmojiDocumentsRequestTypeID, &req)

	if len(req.DocumentID) != 1 || req.DocumentID[0] != 12345 {
		t.Fatalf("ids = %v", req.DocumentID)
	}
}
