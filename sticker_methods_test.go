package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestStickerSetMethods(t *testing.T) {
	fid := documentFileID(t, 0x77)

	t.Run("create", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.StickersCreateStickerSetRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

		b := newMockBot(inv)
		stickers := []InputSticker{{Sticker: FileID(fid), Format: StickerFormatStatic, EmojiList: []string{"😀"}, Keywords: []string{"smile"}}}

		if err := b.CreateNewStickerSet(context.Background(), 99, "s", "T", stickers, WithMaskStickers()); err != nil {
			t.Fatalf("CreateNewStickerSet: %v", err)
		}

		var req tg.StickersCreateStickerSetRequest

		inv.decode(t, tg.StickersCreateStickerSetRequestTypeID, &req)

		if !req.Masks || req.ShortName != "s" || len(req.Stickers) != 1 {
			t.Fatalf("req = %#v", req)
		}
	})

	t.Run("add", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.StickersAddStickerToSetRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

		b := newMockBot(inv)
		if err := b.AddStickerToSet(context.Background(), "s", InputSticker{Sticker: FileID(fid), Format: StickerFormatStatic, EmojiList: []string{"😀"}}); err != nil {
			t.Fatalf("AddStickerToSet: %v", err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.StickersRemoveStickerFromSetRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

		b := newMockBot(inv)
		if err := b.DeleteStickerFromSet(context.Background(), fid); err != nil {
			t.Fatalf("DeleteStickerFromSet: %v", err)
		}
	})

	t.Run("position", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.StickersChangeStickerPositionRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

		b := newMockBot(inv)
		if err := b.SetStickerPositionInSet(context.Background(), fid, 3); err != nil {
			t.Fatalf("SetStickerPositionInSet: %v", err)
		}
	})

	t.Run("thumb", func(t *testing.T) {
		inv := newMockInvoker()
		inv.reply(tg.StickersSetStickerSetThumbRequestTypeID, &tg.MessagesStickerSet{Set: tg.StickerSet{ShortName: "s"}})

		b := newMockBot(inv)
		if err := b.SetStickerSetThumb(context.Background(), "s", FileID(fid)); err != nil {
			t.Fatalf("SetStickerSetThumb: %v", err)
		}
	})
}

func TestUploadStickerFile(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaDocument{
		Document: &tg.Document{ID: 5, AccessHash: 6, FileReference: []byte{1}, DCID: 2, MimeType: "image/png", Size: 100},
	})

	b := newMockBot(inv)

	file, err := b.UploadStickerFile(context.Background(), 99, FileFromBytes("s.png", []byte("img")), StickerFormatStatic)
	if err != nil {
		t.Fatalf("UploadStickerFile: %v", err)
	}

	if file.FileID == "" || file.FileSize != 100 {
		t.Fatalf("file = %#v", file)
	}
}
