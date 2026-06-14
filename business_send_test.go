package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func businessSendReply() *tg.Updates {
	return messageUpdates(&tg.Message{ID: 1, Message: "hi", PeerID: &tg.PeerUser{UserID: 10}})
}

func TestSendMessageWithBusinessConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	m, err := b.SendMessage(context.Background(), userRef(10, 20), "hi", WithBusinessConnection("bc1"))
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	if m.MessageID != 1 {
		t.Fatalf("message id = %d", m.MessageID)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	sm, ok := wrapper.Query.(*tg.MessagesSendMessageRequest)
	if !ok || sm.Message != "hi" {
		t.Fatalf("query = %#v", wrapper.Query)
	}
}

func TestSendMessageWithoutBusinessIsDirect(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	if _, err := b.SendMessage(context.Background(), userRef(10, 20), "hi"); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	if inv.called(tg.InvokeWithBusinessConnectionRequestTypeID) {
		t.Fatal("a non-business send must not be wrapped")
	}

	if !inv.called(tg.MessagesSendMessageRequestTypeID) {
		t.Fatal("expected a direct messages.sendMessage")
	}
}

func TestBusinessContextSendMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)

	if _, err := b.Business("bc3").SendMessage(context.Background(), userRef(10, 20), "hi"); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc3" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}
}

func TestContextReplyBusiness(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		BusinessMessage: &Message{
			MessageID:            9,
			BusinessConnectionID: "bc4",
			Chat:                 Chat{ID: 10, Type: ChatTypePrivate},
		},
	}}

	if _, err := c.Reply("ok"); err != nil {
		t.Fatalf("Reply: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc4" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	sm, ok := wrapper.Query.(*tg.MessagesSendMessageRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := sm.GetReplyTo(); !ok {
		t.Fatal("business reply should set reply_to")
	}
}

func TestBusinessReplyUsesEntitiesPeer(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)
	b.OnBusinessMessage(func(c *Context) error {
		_, err := c.Reply("hi")

		return err
	})

	// The update delivers the chat's access hash (555) in its entities; the reply
	// must address the peer with that hash, not the bot's stored one.
	e := tg.Entities{Users: map[int64]*tg.User{10: {ID: 10, AccessHash: 555}}}
	msg := &tg.Message{ID: 7, PeerID: &tg.PeerUser{UserID: 10}}

	b.dispatchBusinessMessage(context.Background(), e, "bc1", msg, false)

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	sm, ok := wrapper.Query.(*tg.MessagesSendMessageRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	peer, ok := sm.Peer.(*tg.InputPeerUser)
	if !ok || peer.UserID != 10 || peer.AccessHash != 555 {
		t.Fatalf("peer = %#v, want user 10 with access hash 555", sm.Peer)
	}
}

func TestContextSendBusiness(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, businessSendReply())

	b := newMockBot(inv)
	c := &Context{Context: context.Background(), Bot: b, Update: &Update{
		BusinessMessage: &Message{
			MessageID:            9,
			BusinessConnectionID: "bc5",
			Chat:                 Chat{ID: 10, Type: ChatTypePrivate},
		},
	}}

	if _, err := c.Send("hi"); err != nil {
		t.Fatalf("Send: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesSendMessageRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc5" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}
}
