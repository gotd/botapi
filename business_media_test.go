package botapi

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

// businessAlbumHandler answers the wrapped RPCs an album send issues on behalf of
// a business account: uploadMedia per item (returning a usable photo) and the
// final sendMultiMedia.
func businessAlbumHandler(t *testing.T, connSeen *string) func(*bin.Buffer) (bin.Encoder, error) {
	t.Helper()

	return func(buf *bin.Buffer) (bin.Encoder, error) {
		if err := buf.ConsumeID(tg.InvokeWithBusinessConnectionRequestTypeID); err != nil {
			return nil, err
		}

		id, err := buf.String()
		if err != nil {
			return nil, err
		}

		*connSeen = id

		inner, err := buf.PeekID()
		if err != nil {
			return nil, err
		}

		switch inner {
		case tg.MessagesUploadMediaRequestTypeID:
			m := &tg.MessageMediaPhoto{}
			m.SetPhoto(&tg.Photo{
				ID: 1, AccessHash: 2, FileReference: []byte{0}, DCID: 2,
				Sizes: []tg.PhotoSizeClass{&tg.PhotoSize{Type: "x", W: 1, H: 1, Size: 1}},
			})

			return m, nil
		case tg.MessagesSendMultiMediaRequestTypeID:
			return multiBusinessMessageUpdates(11, 12), nil
		default:
			return nil, fmt.Errorf("unexpected inner request %#x", inner)
		}
	}
}

func TestSendPhotoWithBusinessConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	photo := &InputFileUpload{Name: "p.jpg", Bytes: []byte("img")}
	if _, err := b.SendPhoto(context.Background(), userRef(10, 20), photo, "hi", WithBusinessConnection("bc1")); err != nil {
		t.Fatalf("SendPhoto: %v", err)
	}

	// The file upload must stay in the bot session, unwrapped.
	if !inv.called(tg.UploadSaveFilePartRequestTypeID) {
		t.Fatal("expected a direct file upload")
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMediaRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	send, ok := wrapper.Query.(*tg.MessagesSendMediaRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := send.Media.(*tg.InputMediaUploadedPhoto); !ok {
		t.Fatalf("media = %#v, want uploaded photo", send.Media)
	}
}

func TestSendDocumentWithBusinessConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	doc := &InputFileUpload{Name: "f.bin", Bytes: []byte("data")}
	if _, err := b.SendDocument(context.Background(), userRef(10, 20), doc, "", WithBusinessConnection("bc2")); err != nil {
		t.Fatalf("SendDocument: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMediaRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc2" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	send, ok := wrapper.Query.(*tg.MessagesSendMediaRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := send.Media.(*tg.InputMediaUploadedDocument); !ok {
		t.Fatalf("media = %#v, want uploaded document", send.Media)
	}
}

func TestSendPhotoByFileIDWithBusinessConnection(t *testing.T) {
	fid := documentFileID(t, 0x91)

	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	// A file_id send performs no upload; only the send is wrapped.
	if _, err := b.SendPhoto(context.Background(), userRef(10, 20), FileID(fid), "", WithBusinessConnection("bc3")); err != nil {
		t.Fatalf("SendPhoto: %v", err)
	}

	if inv.called(tg.UploadSaveFilePartRequestTypeID) {
		t.Fatal("a file_id send must not upload")
	}

	if !inv.called(tg.InvokeWithBusinessConnectionRequestTypeID) {
		t.Fatal("expected a wrapped send")
	}
}

func TestSendMediaGroupWithBusinessConnection(t *testing.T) {
	inv := newMockInvoker()

	var connSeen string

	inv.handle(tg.InvokeWithBusinessConnectionRequestTypeID, businessAlbumHandler(t, &connSeen))

	b := newMockBot(inv)

	media := []InputMedia{
		&InputMediaPhoto{Media: &InputFileUpload{Name: "a.jpg", Bytes: []byte("a")}},
		&InputMediaPhoto{Media: &InputFileUpload{Name: "b.jpg", Bytes: []byte("b")}},
	}

	msgs, err := b.SendMediaGroup(context.Background(), userRef(10, 20), media, WithBusinessConnection("bc7"))
	if err != nil {
		t.Fatalf("SendMediaGroup: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("got %d messages", len(msgs))
	}

	if connSeen != "bc7" {
		t.Fatalf("connection id = %q", connSeen)
	}

	// Uploads happen in the bot session, unwrapped.
	if !inv.called(tg.UploadSaveFilePartRequestTypeID) {
		t.Fatal("expected direct file uploads")
	}

	// The final wrapped call is the multi-media send.
	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMultiMediaRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if _, ok := wrapper.Query.(*tg.MessagesSendMultiMediaRequest); !ok {
		t.Fatalf("final query = %#v, want sendMultiMedia", wrapper.Query)
	}
}
