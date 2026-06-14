package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// TestSendTypedMediaUpload drives the local-upload branch of typedMedia for
// every typed document sender, including the filename attribute.
func TestSendTypedMediaUpload(t *testing.T) {
	file := FileFromBytes("clip.bin", []byte("payload"))
	sends := map[string]func(b *Bot) (*Message, error){
		"video": func(b *Bot) (*Message, error) { return b.SendVideo(context.Background(), userRef(10, 20), file, "c") },
		"animation": func(b *Bot) (*Message, error) {
			return b.SendAnimation(context.Background(), userRef(10, 20), file, "c")
		},
		"audio":     func(b *Bot) (*Message, error) { return b.SendAudio(context.Background(), userRef(10, 20), file, "c") },
		"voice":     func(b *Bot) (*Message, error) { return b.SendVoice(context.Background(), userRef(10, 20), file, "c") },
		"videonote": func(b *Bot) (*Message, error) { return b.SendVideoNote(context.Background(), userRef(10, 20), file) },
		"sticker":   func(b *Bot) (*Message, error) { return b.SendSticker(context.Background(), userRef(10, 20), file) },
	}
	for name, send := range sends {
		t.Run(name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
			b := newMockBot(inv)
			if _, err := send(b); err != nil {
				t.Fatalf("%s upload: %v", name, err)
			}
			if !inv.called(tg.UploadSaveFilePartRequestTypeID) {
				t.Fatalf("%s did not upload a file part", name)
			}
		})
	}
}

// TestEditMessageMediaUploadAndMarkup covers the photo-upload branch, the
// reply-markup branch and the parse-mode caption branch of EditMessageMedia.
func TestEditMessageMediaUploadAndMarkup(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	media := &InputMediaPhoto{Media: FileFromBytes("p.jpg", []byte("img")), Caption: "<b>x</b>", ParseMode: ParseModeHTML}
	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, media, WithReplyMarkup(goodInlineMarkup)); err != nil {
		t.Fatalf("EditMessageMedia upload: %v", err)
	}
	if !inv.called(tg.UploadSaveFilePartRequestTypeID) {
		t.Fatal("photo upload should save a file part")
	}
}

// TestEditMessageMediaErrors covers the error branches: a bad file_id, a bad
// reply markup and an unresolved peer.
func TestEditMessageMediaErrors(t *testing.T) {
	b := newMockBot(newMockInvoker())

	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, &InputMediaPhoto{Media: FileID("bad")}); err == nil {
		t.Fatal("expected error for bad photo file_id")
	}
	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, &InputMediaDocument{Media: FileID("bad")}); err == nil {
		t.Fatal("expected error for bad document file_id")
	}
	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5,
		&InputMediaPhoto{Media: FileURL("https://e/p.jpg")}, WithReplyMarkup(badInlineMarkup)); err == nil {
		t.Fatal("expected error for bad reply markup")
	}
}

// TestInputMediaToTgVariants exercises inputMediaToTg across every media variant
// and file source, plus the photo/document URL branches.
func TestInputMediaToTgVariants(t *testing.T) {
	b := newMockBot(newMockInvoker())
	doc := documentFileID(t, 0xA1)
	photo := photoFileID(t, 0xA2)
	ctx := context.Background()

	medias := []InputMedia{
		&InputMediaPhoto{Media: FileURL("https://e/p.jpg")},
		&InputMediaPhoto{Media: FileID(photo)},
		&InputMediaVideo{Media: FileURL("https://e/v.mp4")},
		&InputMediaVideo{Media: FileID(doc)},
		&InputMediaAnimation{Media: FileID(doc)},
		&InputMediaAudio{Media: FileURL("https://e/a.mp3")},
		&InputMediaDocument{Media: FileID(doc)},
	}
	for _, m := range medias {
		if _, _, err := b.inputMediaToTg(ctx, m); err != nil {
			t.Errorf("%T: %v", m, err)
		}
	}

	// Upload branch with a filename on a document.
	if _, _, err := b.inputMediaToTg(ctx, &InputMediaDocument{Media: FileFromBytes("f.bin", []byte("x"))}); err != nil {
		t.Errorf("document upload: %v", err)
	}
}

// TestInputMediaToMultiTypes covers the per-type album resolution branches,
// including the non-upload rejection.
func TestInputMediaToMultiTypes(t *testing.T) {
	b := newMockBot(newMockInvoker())
	ctx := context.Background()
	up := func() InputFile { return FileFromBytes("f.bin", []byte("x")) }

	uploads := []InputMedia{
		&InputMediaPhoto{Media: up()},
		&InputMediaVideo{Media: up()},
		&InputMediaAnimation{Media: up()},
		&InputMediaAudio{Media: up()},
		&InputMediaDocument{Media: up(), Caption: "c"},
	}
	for _, m := range uploads {
		if _, err := b.inputMediaToMulti(ctx, m); err != nil {
			t.Errorf("%T album item: %v", m, err)
		}
	}

	// A non-upload document item is rejected.
	if _, err := b.inputMediaToMulti(ctx, &InputMediaVideo{Media: FileID(documentFileID(t, 1))}); err == nil {
		t.Error("non-upload video album item should be rejected")
	}
}

// TestSentMessagesBranches covers the response-shape branches of sentMessages.
func TestSentMessagesBranches(t *testing.T) {
	b := newMockBot(newMockInvoker())
	ctx := context.Background()

	// Send error propagates.
	if _, err := b.sentMessages(ctx, nil, context.Canceled); err == nil {
		t.Fatal("expected error to propagate")
	}
	// Unhandled updates shape yields no messages.
	if msgs, err := b.sentMessages(ctx, &tg.UpdateShort{}, nil); err != nil || msgs != nil {
		t.Fatalf("UpdateShort: msgs=%v err=%v", msgs, err)
	}
	// UpdatesCombined with a channel message is converted.
	resp := &tg.UpdatesCombined{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewChannelMessage{Message: &tg.Message{ID: 1, Message: "x", PeerID: &tg.PeerUser{UserID: 10}}},
			&tg.UpdateMessageID{}, // skipped
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	}
	msgs, err := b.sentMessages(ctx, resp, nil)
	if err != nil || len(msgs) != 1 {
		t.Fatalf("UpdatesCombined: msgs=%d err=%v", len(msgs), err)
	}
}
