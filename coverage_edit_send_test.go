package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestEditMessageTextOptions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, Message: "new", PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)
	ctx := context.Background()

	// Reply markup + web-preview-disabled exercise the optional builder branches,
	// and the edit response drives editedMessageFromResp through sentMessage.
	if _, err := b.EditMessageText(ctx, userRef(10, 20), 5, "new",
		WithReplyMarkup(goodInlineMarkup), DisableWebPagePreview()); err != nil {
		t.Fatalf("EditMessageText: %v", err)
	}

	if _, err := b.EditMessageCaption(ctx, userRef(10, 20), 5, "cap"); err != nil {
		t.Fatalf("EditMessageCaption: %v", err)
	}

	// A bad reply markup is rejected before any RPC.
	if _, err := b.EditMessageText(ctx, userRef(10, 20), 5, "x", WithReplyMarkup(badInlineMarkup)); err == nil {
		t.Fatal("expected markup error")
	}
}

func TestEditMessageReplyMarkupBranches(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)
	ctx := context.Background()

	if _, err := b.EditMessageReplyMarkup(ctx, userRef(10, 20), 5, goodInlineMarkup); err != nil {
		t.Fatalf("set markup: %v", err)
	}

	if _, err := b.EditMessageReplyMarkup(ctx, userRef(10, 20), 5, nil); err != nil {
		t.Fatalf("clear markup: %v", err)
	}

	if _, err := b.EditMessageReplyMarkup(ctx, userRef(10, 20), 5, badInlineMarkup); err == nil {
		t.Fatal("expected markup error")
	}
}

// TestSendMessageAllOptions covers every branch of applySendConfig via a fully
// configured send.
func TestSendMessageAllOptions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, messageUpdates(&tg.Message{ID: 1, Message: "hi", PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)

	_, err := b.SendMessage(context.Background(), userRef(10, 20), "hi",
		Silent(), ProtectContent(), ReplyTo(3), DisableWebPagePreview(), WithReplyMarkup(goodInlineMarkup))
	if err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	// A bad reply markup is rejected.
	if _, err := b.SendMessage(context.Background(), userRef(10, 20), "hi", WithReplyMarkup(badInlineMarkup)); err == nil {
		t.Fatal("expected markup error")
	}
}

// TestEditedMessageFromResp covers the response-shape branches of the edit
// unpacker.
func TestEditedMessageFromRespShapes(t *testing.T) {
	msg := &tg.Message{ID: 1}
	cases := []struct {
		name string
		resp tg.UpdatesClass
		ok   bool
	}{
		{"updates", &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateEditMessage{Message: msg}}}, true},
		{"combined", &tg.UpdatesCombined{Updates: []tg.UpdateClass{&tg.UpdateEditChannelMessage{Message: msg}}}, true},
		{"short", &tg.UpdateShort{Update: &tg.UpdateEditMessage{Message: msg}}, true},
		{"none", &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateMessageID{}}}, false},
		{"unhandled", &tg.UpdatesTooLong{}, false},
	}

	for _, c := range cases {
		if _, ok := editedMessageFromResp(c.resp); ok != c.ok {
			t.Errorf("%s: ok=%v, want %v", c.name, ok, c.ok)
		}
	}
}
