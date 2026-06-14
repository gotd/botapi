package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetChatPhotoUpload(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditPhotoRequestTypeID, okUpdates())

	b := newMockBot(inv)

	err := b.SetChatPhoto(context.Background(), tdlibChannel(50), FileFromBytes("p.jpg", []byte("img")))
	if err != nil {
		t.Fatalf("SetChatPhoto: %v", err)
	}

	if !inv.called(tg.ChannelsEditPhotoRequestTypeID) {
		t.Fatal("expected channels.editPhoto")
	}
}

func TestSetChatPhotoRejectsNonUpload(t *testing.T) {
	b := newMockBot(newMockInvoker())
	if err := b.SetChatPhoto(context.Background(), tdlibChannel(50), FileID("x")); err == nil {
		t.Fatal("SetChatPhoto should reject a non-uploaded file")
	}
}
