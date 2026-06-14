package botapi

import (
	"bytes"
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

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
