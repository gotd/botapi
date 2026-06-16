package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendPaidMediaUpload(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaPhoto{
		Photo: &tg.Photo{ID: 7, AccessHash: 8, FileReference: []byte{0}, Sizes: []tg.PhotoSizeClass{
			&tg.PhotoSize{Type: "x", W: 640, H: 640, Size: 1000},
		}},
	})
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())

	b := newMockBot(inv)

	media := []InputPaidMedia{
		InputPaidMediaPhoto{Media: FileFromBytes("p.jpg", []byte("data"))},
	}
	if _, err := b.SendPaidMedia(context.Background(), userRef(10, 20), 50, media, "cap",
		WithPaidMediaPayload("internal")); err != nil {
		t.Fatalf("SendPaidMedia: %v", err)
	}

	if !inv.called(tg.MessagesUploadMediaRequestTypeID) {
		t.Fatal("uploaded photo should be finalized via uploadMedia")
	}

	var req tg.MessagesSendMediaRequest

	inv.decode(t, tg.MessagesSendMediaRequestTypeID, &req)

	paid, ok := req.Media.(*tg.InputMediaPaidMedia)
	if !ok {
		t.Fatalf("media = %#v, want paid media", req.Media)
	}

	if paid.StarsAmount != 50 {
		t.Fatalf("stars = %d, want 50", paid.StarsAmount)
	}

	if payload, ok := paid.GetPayload(); !ok || payload != "internal" {
		t.Fatalf("payload = %q ok=%v", payload, ok)
	}

	if len(paid.ExtendedMedia) != 1 {
		t.Fatalf("extended media = %#v", paid.ExtendedMedia)
	}

	if _, ok := paid.ExtendedMedia[0].(*tg.InputMediaPhoto); !ok {
		t.Fatalf("item 0 = %#v, want referenced photo", paid.ExtendedMedia[0])
	}
}

func TestSendPaidMediaURLVideo(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())

	b := newMockBot(inv)

	media := []InputPaidMedia{
		InputPaidMediaVideo{Media: FileURL("https://example.com/v.mp4"), SupportsStreaming: true},
	}
	if _, err := b.SendPaidMedia(context.Background(), userRef(10, 20), 10, media, ""); err != nil {
		t.Fatalf("SendPaidMedia: %v", err)
	}

	var req tg.MessagesSendMediaRequest

	inv.decode(t, tg.MessagesSendMediaRequestTypeID, &req)

	paid, ok := req.Media.(*tg.InputMediaPaidMedia)
	if !ok {
		t.Fatalf("media = %#v, want paid media", req.Media)
	}

	if _, ok := paid.ExtendedMedia[0].(*tg.InputMediaDocumentExternal); !ok {
		t.Fatalf("item 0 = %#v, want external document", paid.ExtendedMedia[0])
	}
}

func TestSendPaidMediaCountValidation(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).SendPaidMedia(context.Background(), userRef(10, 20), 5, nil, ""); err == nil {
		t.Fatal("expected error for empty media")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
