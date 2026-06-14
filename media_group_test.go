package botapi

import (
	"context"
	"errors"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendMediaGroup(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaPhoto{
		Photo: &tg.Photo{ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2, Sizes: []tg.PhotoSizeClass{&tg.PhotoSize{Type: "x", W: 1, H: 1, Size: 1}}},
	})
	inv.reply(tg.MessagesSendMultiMediaRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}},
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 2, PeerID: &tg.PeerUser{UserID: 10}}},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	})

	b := newMockBot(inv)

	media := []InputMedia{
		&InputMediaPhoto{Media: FileFromBytes("a.jpg", []byte("a"))},
		&InputMediaPhoto{Media: FileFromBytes("b.jpg", []byte("b"))},
	}

	msgs, err := b.SendMediaGroup(context.Background(), userRef(10, 20), media)
	if err != nil {
		t.Fatalf("SendMediaGroup: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("messages = %d", len(msgs))
	}
}

func TestSendMediaGroupDocuments(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesUploadMediaRequestTypeID, &tg.MessageMediaDocument{
		Document: &tg.Document{ID: 1, AccessHash: 2, FileReference: []byte{1}, DCID: 2, MimeType: "application/pdf", Size: 10},
	})
	inv.reply(tg.MessagesSendMultiMediaRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}},
			&tg.UpdateNewMessage{Message: &tg.Message{ID: 2, PeerID: &tg.PeerUser{UserID: 10}}},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	})

	b := newMockBot(inv)

	media := []InputMedia{
		&InputMediaDocument{Media: FileFromBytes("a.pdf", []byte("a")), Caption: "first"},
		&InputMediaVideo{Media: FileFromBytes("b.mp4", []byte("b"))},
	}

	msgs, err := b.SendMediaGroup(context.Background(), userRef(10, 20), media)
	if err != nil {
		t.Fatalf("SendMediaGroup docs: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("messages = %d", len(msgs))
	}
}

func TestSendMediaGroupCountValidation(t *testing.T) {
	b := newTestBot(t)
	_, err := b.SendMediaGroup(context.Background(), ID(1), []InputMedia{
		&InputMediaPhoto{Type: InputMediaPhotoType, Media: FileFromBytes("a", []byte("x"))},
	})

	var apiErr *Error

	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("single-item group should be a 400, got %v", err)
	}
}

func TestMediaGroupRejectsNonUpload(t *testing.T) {
	b := newTestBot(t)
	_, err := b.inputMediaToMulti(context.Background(), &InputMediaPhoto{
		Type:  InputMediaPhotoType,
		Media: FileURL("https://example.com/a.jpg"),
	})

	if !errors.Is(err, error(errNonUploadInAlbum)) {
		t.Fatalf("URL item in album should be rejected, got %v", err)
	}
}
