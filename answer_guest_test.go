package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestAnswerGuestQuery(t *testing.T) {
	mid := &tg.InputBotInlineMessageID64{DCID: 2, OwnerID: 1, ID: 9, AccessHash: 42}

	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotGuestChatResultRequestTypeID, mid)

	result := &InlineQueryResultArticle{
		ID:                  "r1",
		Title:               "Title",
		InputMessageContent: &InputTextMessageContent{MessageText: "hi"},
	}

	got, err := newMockBot(inv).AnswerGuestQuery(context.Background(), "12345", result)
	if err != nil {
		t.Fatalf("AnswerGuestQuery: %v", err)
	}

	decoded, err := decodeInlineMessageID(got.InlineMessageID)
	if err != nil {
		t.Fatalf("decode inline message id: %v", err)
	}

	if decoded.String() != mid.String() {
		t.Fatalf("inline message id = %#v, want %#v", decoded, mid)
	}

	var req tg.MessagesSetBotGuestChatResultRequest

	inv.decode(t, tg.MessagesSetBotGuestChatResultRequestTypeID, &req)

	if req.QueryID != 12345 {
		t.Fatalf("query id = %d", req.QueryID)
	}

	if _, ok := req.Result.(*tg.InputBotInlineResult); !ok {
		t.Fatalf("result = %#v, want inline result", req.Result)
	}
}

func TestAnswerGuestQueryInvalid(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).AnswerGuestQuery(context.Background(), "12345", nil); err == nil {
		t.Fatal("expected error for nil result")
	}

	if _, err := newMockBot(inv).AnswerGuestQuery(context.Background(), "notanumber",
		&InlineQueryResultArticle{ID: "r1", Title: "T", InputMessageContent: &InputTextMessageContent{MessageText: "x"}}); err == nil {
		t.Fatal("expected error for invalid guest_query_id")
	}
}
