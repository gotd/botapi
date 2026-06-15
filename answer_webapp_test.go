package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestAnswerWebAppQuery(t *testing.T) {
	mid := &tg.InputBotInlineMessageID64{DCID: 2, OwnerID: 1, ID: 9, AccessHash: 42}

	sent := &tg.WebViewMessageSent{}
	sent.SetMsgID(mid)

	inv := newMockInvoker()
	inv.reply(tg.MessagesSendWebViewResultMessageRequestTypeID, sent)

	result := &InlineQueryResultArticle{
		ID:                  "r1",
		Title:               "Title",
		InputMessageContent: &InputTextMessageContent{MessageText: "hi"},
	}

	got, err := newMockBot(inv).AnswerWebAppQuery(context.Background(), "wq1", result)
	if err != nil {
		t.Fatalf("AnswerWebAppQuery: %v", err)
	}

	decoded, err := decodeInlineMessageID(got.InlineMessageID)
	if err != nil {
		t.Fatalf("decode inline message id %q: %v", got.InlineMessageID, err)
	}

	if decoded.String() != mid.String() {
		t.Fatalf("inline message id = %#v, want %#v", decoded, mid)
	}

	var req tg.MessagesSendWebViewResultMessageRequest

	inv.decode(t, tg.MessagesSendWebViewResultMessageRequestTypeID, &req)

	if req.BotQueryID != "wq1" {
		t.Fatalf("bot query id = %q", req.BotQueryID)
	}

	if _, ok := req.Result.(*tg.InputBotInlineResult); !ok {
		t.Fatalf("result = %#v, want inline result", req.Result)
	}
}

func TestAnswerWebAppQueryNilResult(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).AnswerWebAppQuery(context.Background(), "wq1", nil); err == nil {
		t.Fatal("expected error for nil result")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
