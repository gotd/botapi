package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendMessageDraft(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetTypingRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SendMessageDraft(context.Background(), userRef(10, 20), 999, "draft text",
		WithDraftMessageThread(5)); err != nil {
		t.Fatalf("SendMessageDraft: %v", err)
	}

	var req tg.MessagesSetTypingRequest

	inv.decode(t, tg.MessagesSetTypingRequestTypeID, &req)

	if topMsg, ok := req.GetTopMsgID(); !ok || topMsg != 5 {
		t.Fatalf("top msg id = %d ok=%v", topMsg, ok)
	}

	action, ok := req.Action.(*tg.SendMessageTextDraftAction)
	if !ok {
		t.Fatalf("action = %#v, want text draft", req.Action)
	}

	if action.RandomID != 999 || action.Text.Text != "draft text" {
		t.Fatalf("action = %#v", action)
	}
}

func TestSendRichMessageDraft(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetTypingRequestTypeID, &tg.BoolTrue{})

	message := RichMessage{
		Blocks: []tg.PageBlockClass{&tg.PageBlockParagraph{Text: &tg.TextPlain{Text: "partial"}}},
		RTL:    true,
		Part:   true,
	}

	if err := newMockBot(inv).SendRichMessageDraft(context.Background(), userRef(10, 20), 1000, message); err != nil {
		t.Fatalf("SendRichMessageDraft: %v", err)
	}

	var req tg.MessagesSetTypingRequest

	inv.decode(t, tg.MessagesSetTypingRequestTypeID, &req)

	action, ok := req.Action.(*tg.SendMessageRichMessageDraftAction)
	if !ok {
		t.Fatalf("action = %#v, want rich draft", req.Action)
	}

	if action.RandomID != 1000 || len(action.RichMessage.Blocks) != 1 {
		t.Fatalf("action = %#v", action)
	}

	if !action.RichMessage.GetRtl() {
		t.Fatal("rtl flag should be set")
	}

	if _, ok := action.RichMessage.Blocks[0].(*tg.PageBlockParagraph); !ok {
		t.Fatalf("block 0 = %#v, want paragraph", action.RichMessage.Blocks[0])
	}
}
