package botapi

import (
	"bytes"
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func sendMediaOK() *tg.Updates {
	return messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}})
}

func TestSendMediaVariantsURL(t *testing.T) {
	url := FileURL("https://e/f.bin")
	sends := map[string]func(b *Bot) (*Message, error){
		"photo":    func(b *Bot) (*Message, error) { return b.SendPhoto(context.Background(), userRef(10, 20), url, "c") },
		"document": func(b *Bot) (*Message, error) { return b.SendDocument(context.Background(), userRef(10, 20), url, "c") },
		"video":    func(b *Bot) (*Message, error) { return b.SendVideo(context.Background(), userRef(10, 20), url, "c") },
		"animation": func(b *Bot) (*Message, error) {
			return b.SendAnimation(context.Background(), userRef(10, 20), url, "c")
		},
		"audio":     func(b *Bot) (*Message, error) { return b.SendAudio(context.Background(), userRef(10, 20), url, "c") },
		"voice":     func(b *Bot) (*Message, error) { return b.SendVoice(context.Background(), userRef(10, 20), url, "c") },
		"videonote": func(b *Bot) (*Message, error) { return b.SendVideoNote(context.Background(), userRef(10, 20), url) },
		"sticker": func(b *Bot) (*Message, error) {
			return b.SendSticker(context.Background(), userRef(10, 20), FileID(documentFileID(t, 1)))
		},
	}
	for name, send := range sends {
		t.Run(name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
			b := newMockBot(inv)
			if _, err := send(b); err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if !inv.called(tg.MessagesSendMediaRequestTypeID) {
				t.Fatalf("%s did not call messages.sendMedia", name)
			}
		})
	}
}

func TestSendPhotoUpload(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	if _, err := b.SendPhoto(context.Background(), userRef(10, 20), FileFromBytes("p.jpg", []byte("data")), "cap"); err != nil {
		t.Fatalf("SendPhoto upload: %v", err)
	}
	if !inv.called(tg.UploadSaveFilePartRequestTypeID) {
		t.Fatal("upload should save a file part")
	}
}

func TestSendDocumentReader(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	r := bytes.NewReader([]byte("streamed content"))
	if _, err := b.SendDocument(context.Background(), userRef(10, 20), FileFromReader("f.txt", r), ""); err != nil {
		t.Fatalf("SendDocument reader: %v", err)
	}
}

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

func TestBuilders(t *testing.T) {
	if kb := Keyboard(Row(Button("a"), ButtonContact("c")), Row(Button("b"))); len(kb.Keyboard) != 2 {
		t.Fatalf("keyboard = %#v", kb)
	}
	if ik := InlineKeyboard([]InlineKeyboardButton{InlineButtonData("x", "d"), InlineButtonURL("u", "https://e")}); len(ik.InlineKeyboard) != 1 {
		t.Fatalf("inline keyboard = %#v", ik)
	}
	for _, f := range []InputFile{
		FileID("id"), FileURL("https://e"), FileFromPath("/tmp/x"),
		FileFromBytes("n", []byte("d")), FileFromReader("n", bytes.NewReader(nil)),
	} {
		if f == nil {
			t.Fatal("nil input file")
		}
	}
	if e := Emoji("👍"); e == nil {
		t.Fatal("nil emoji reaction")
	}
	if e := CustomEmoji("123"); e == nil {
		t.Fatal("nil custom emoji reaction")
	}
}
