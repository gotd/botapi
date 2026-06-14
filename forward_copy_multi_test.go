package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func multiMessageUpdates(ids ...int) *tg.Updates {
	upds := make([]tg.UpdateClass, 0, len(ids))

	for _, id := range ids {
		upds = append(upds, &tg.UpdateNewMessage{
			Message: &tg.Message{ID: id, PeerID: &tg.PeerUser{UserID: 10}},
		})
	}

	return &tg.Updates{
		Updates: upds,
		Users:   []tg.UserClass{&tg.User{ID: 10, AccessHash: 1}},
	}
}

func TestForwardMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesForwardMessagesRequestTypeID, multiMessageUpdates(7, 8))

	b := newMockBot(inv)

	msgs, err := b.ForwardMessages(context.Background(), userRef(10, 20), userRef(30, 40), []int{1, 2})
	if err != nil {
		t.Fatalf("ForwardMessages: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("got %d messages", len(msgs))
	}

	var req tg.MessagesForwardMessagesRequest

	inv.decode(t, tg.MessagesForwardMessagesRequestTypeID, &req)

	if len(req.ID) != 2 || req.ID[0] != 1 || req.ID[1] != 2 {
		t.Fatalf("ids = %v", req.ID)
	}

	if req.DropAuthor {
		t.Fatal("forward should not drop author")
	}
}

func TestForwardMessagesEmpty(t *testing.T) {
	inv := newMockInvoker()

	msgs, err := newMockBot(inv).ForwardMessages(context.Background(), userRef(10, 20), userRef(30, 40), nil)
	if err != nil || msgs != nil {
		t.Fatalf("got %v, %v", msgs, err)
	}

	if inv.count() != 0 {
		t.Fatal("empty forward should make no RPC")
	}
}

func TestCopyMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesForwardMessagesRequestTypeID, multiMessageUpdates(7, 8, 9))

	b := newMockBot(inv)

	msgs, err := b.CopyMessages(context.Background(), userRef(10, 20), userRef(30, 40), []int{1, 2, 3})
	if err != nil {
		t.Fatalf("CopyMessages: %v", err)
	}

	if len(msgs) != 3 {
		t.Fatalf("got %d messages", len(msgs))
	}

	var req tg.MessagesForwardMessagesRequest

	inv.decode(t, tg.MessagesForwardMessagesRequestTypeID, &req)

	if !req.DropAuthor {
		t.Fatal("copy should drop author")
	}
}
