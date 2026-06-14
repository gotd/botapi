package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// uploadDocInvoker replies to messages.uploadMedia with a saved document, the
// shape the sticker upload path expects.
func uploadDocInvoker() *mockInvoker {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaDocument{
		Document: &tg.Document{ID: 9, AccessHash: 8, FileReference: []byte{1}, DCID: 2, MimeType: "image/png", Size: 42},
	})
	return inv
}

func TestUploadStickerFileBranches(t *testing.T) {
	ctx := context.Background()

	// A file_id (non-upload) is rejected.
	if _, err := newMockBot(newMockInvoker()).UploadStickerFile(ctx, 1, FileID("x"), StickerFormatStatic); err == nil {
		t.Fatal("non-upload sticker file should fail")
	}

	// A successful upload yields a file_id.
	b := newMockBot(uploadDocInvoker())
	f, err := b.UploadStickerFile(ctx, 1, FileFromBytes("s.png", []byte("img")), StickerFormatStatic)
	if err != nil {
		t.Fatalf("UploadStickerFile: %v", err)
	}
	if f.FileID == "" || f.FileSize != 42 {
		t.Fatalf("file = %#v", f)
	}
}

func TestResolveStickerItemBranches(t *testing.T) {
	ctx := context.Background()
	b := newMockBot(uploadDocInvoker())

	// file_id source with keywords.
	item, err := b.resolveStickerItem(ctx, InputSticker{
		Sticker:   FileID(documentFileID(t, 0x10)),
		Format:    StickerFormatStatic,
		EmojiList: []string{"🙂"},
		Keywords:  []string{"smile", "happy"},
	})
	if err != nil {
		t.Fatalf("file_id item: %v", err)
	}
	if kw, ok := item.GetKeywords(); !ok || kw != "smile,happy" {
		t.Fatalf("keywords = %q ok=%v", kw, ok)
	}

	// Upload source.
	if _, err := b.resolveStickerItem(ctx, InputSticker{
		Sticker: FileFromBytes("s.png", []byte("img")), Format: StickerFormatStatic, EmojiList: []string{"x"},
	}); err != nil {
		t.Fatalf("upload item: %v", err)
	}

	// Bad file_id and unsupported URL sources are rejected.
	if _, err := b.resolveStickerItem(ctx, InputSticker{Sticker: FileID("bad")}); err == nil {
		t.Fatal("bad file_id should fail")
	}
	if _, err := b.resolveStickerItem(ctx, InputSticker{Sticker: FileURL("https://e/s.png")}); err == nil {
		t.Fatal("URL sticker should fail")
	}
}

// TestUploadStickerDocumentBadResponse covers the unexpected-response branch of
// uploadStickerDocument.
func TestUploadStickerDocumentBadResponse(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaEmpty{})
	b := newMockBot(inv)
	if _, err := b.uploadStickerDocument(context.Background(), &InputFileUpload{Name: "s", Bytes: []byte("x")}, StickerFormatStatic); err == nil {
		t.Fatal("expected error for non-document upload response")
	}
}
