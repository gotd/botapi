package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestEditMessageMediaVariants(t *testing.T) {
	fid := documentFileID(t, 0x90)
	cases := map[string]InputMedia{
		"video-url":      &InputMediaVideo{Media: FileURL("https://e/v.mp4"), Caption: "c"},
		"animation-id":   &InputMediaAnimation{Media: FileID(fid)},
		"audio-url":      &InputMediaAudio{Media: FileURL("https://e/a.mp3")},
		"document-id":    &InputMediaDocument{Media: FileID(fid)},
		"photo-id":       &InputMediaPhoto{Media: FileID(photoFileID(t, 0x91))},
		"document-bytes": &InputMediaDocument{Media: FileFromBytes("f.bin", []byte("data"))},
	}

	for name, media := range cases {
		t.Run(name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

			b := newMockBot(inv)
			if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, media); err != nil {
				t.Fatalf("EditMessageMedia(%s): %v", name, err)
			}
		})
	}
}
