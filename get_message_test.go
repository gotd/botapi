package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetMessageExposesRaw(t *testing.T) {
	inv := newMockInvoker()

	raw := &tg.Message{ID: 7, Message: "hi", PeerID: &tg.PeerChannel{ChannelID: 50}}
	inv.reply(tg.ChannelsGetMessagesRequestTypeID, &tg.MessagesChannelMessages{
		Messages: []tg.MessageClass{raw},
	})

	b := newMockBot(inv)

	m, err := b.GetMessage(context.Background(), tdlibChannel(50), 7)
	if err != nil {
		t.Fatalf("GetMessage: %v", err)
	}

	if m.MessageID != 7 || m.Text != "hi" {
		t.Fatalf("message = %#v", m)
	}

	if m.Raw() == nil || m.Raw().ID != 7 || m.Raw().Message != "hi" {
		t.Fatalf("Raw() = %#v, want the fetched tg.Message", m.Raw())
	}

	var req tg.ChannelsGetMessagesRequest

	inv.decode(t, tg.ChannelsGetMessagesRequestTypeID, &req)

	if len(req.ID) != 1 {
		t.Fatalf("requested ids = %#v", req.ID)
	}
}

func TestGetMessageNotFound(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsGetMessagesRequestTypeID, &tg.MessagesChannelMessages{})

	b := newMockBot(inv)

	if _, err := b.GetMessage(context.Background(), tdlibChannel(50), 7); err == nil {
		t.Fatal("expected error for missing message")
	}
}
