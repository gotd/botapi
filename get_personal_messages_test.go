package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetUserPersonalChatMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesGetPersonalChannelHistoryRequestTypeID, &tg.MessagesMessages{
		Messages: []tg.MessageClass{
			&tg.Message{ID: 7, Message: "hello", PeerID: &tg.PeerUser{UserID: 99}},
		},
		Users: []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})

	got, err := newMockBot(inv).GetUserPersonalChatMessages(context.Background(), 99, 10)
	if err != nil {
		t.Fatalf("GetUserPersonalChatMessages: %v", err)
	}

	if len(got) != 1 || got[0].MessageID != 7 || got[0].Text != "hello" {
		t.Fatalf("messages = %#v", got)
	}

	var req tg.MessagesGetPersonalChannelHistoryRequest

	inv.decode(t, tg.MessagesGetPersonalChannelHistoryRequestTypeID, &req)

	if req.Limit != 10 {
		t.Fatalf("limit = %d", req.Limit)
	}
}

func TestGetUserPersonalChatMessagesBadLimit(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).GetUserPersonalChatMessages(context.Background(), 99, 0); err == nil {
		t.Fatal("expected error for non-positive limit")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
