package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendLivePhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaDocument{
		Document: &tg.Document{
			ID: 5, AccessHash: 6, FileReference: []byte{0}, MimeType: "video/mp4", Size: 2000, DCID: 2,
			Attributes: []tg.DocumentAttributeClass{&tg.DocumentAttributeVideo{}},
		},
	})
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())

	b := newMockBot(inv)

	_, err := b.SendLivePhoto(context.Background(), userRef(10, 20),
		FileFromBytes("still.jpg", []byte("img")),
		FileFromBytes("live.mp4", []byte("vid")),
		"caption")
	if err != nil {
		t.Fatalf("SendLivePhoto: %v", err)
	}

	if !inv.called(tg.MessagesUploadMediaRequestTypeID) {
		t.Fatal("live photo video should be finalized via uploadMedia")
	}

	var req tg.MessagesSendMediaRequest

	inv.decode(t, tg.MessagesSendMediaRequestTypeID, &req)

	media, ok := req.Media.(*tg.InputMediaUploadedPhoto)
	if !ok {
		t.Fatalf("media = %#v, want uploaded photo", req.Media)
	}

	if !media.LivePhoto {
		t.Fatal("LivePhoto flag should be set")
	}

	if _, ok := media.GetVideo(); !ok {
		t.Fatal("live photo should attach a video document")
	}
}

func TestSendLivePhotoStillMustUpload(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).SendLivePhoto(context.Background(), userRef(10, 20),
		FileID("somefileid"), FileFromBytes("live.mp4", []byte("vid")), ""); err == nil {
		t.Fatal("expected error when still photo is not an upload")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
