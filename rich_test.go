package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"
)

func TestSendRichMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	msg := rich.New(rich.Paragraph(rich.Plain("hi"))).Input()
	if _, err := b.SendRichMessage(context.Background(), userRef(10, 20), msg); err != nil {
		t.Fatalf("SendRichMessage: %v", err)
	}
	if !inv.called(tg.MessagesSendMessageRequestTypeID) {
		t.Fatal("expected messages.sendMessage")
	}
}

func TestSendRichHTMLAndMarkdown(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMessageRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	if _, err := b.SendRichHTML(context.Background(), userRef(10, 20), "<p>hi</p>"); err != nil {
		t.Fatalf("SendRichHTML: %v", err)
	}
	if _, err := b.SendRichMarkdown(context.Background(), userRef(10, 20), "**hi**"); err != nil {
		t.Fatalf("SendRichMarkdown: %v", err)
	}
}
