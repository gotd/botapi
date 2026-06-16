package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSavePreparedInlineMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSavePreparedInlineMessageRequestTypeID, &tg.MessagesBotPreparedInlineMessage{
		ID:         "prep1",
		ExpireDate: 1700000000,
	})

	result := &InlineQueryResultArticle{
		ID:                  "r1",
		Title:               "Title",
		InputMessageContent: &InputTextMessageContent{MessageText: "hi"},
	}

	got, err := newMockBot(inv).SavePreparedInlineMessage(context.Background(), 99, result,
		WithAllowUserChats(), WithAllowChannelChats())
	if err != nil {
		t.Fatalf("SavePreparedInlineMessage: %v", err)
	}

	if got.ID != "prep1" || got.ExpirationDate != 1700000000 {
		t.Fatalf("got = %#v", got)
	}

	var req tg.MessagesSavePreparedInlineMessageRequest

	inv.decode(t, tg.MessagesSavePreparedInlineMessageRequestTypeID, &req)

	if _, ok := req.Result.(*tg.InputBotInlineResult); !ok {
		t.Fatalf("result = %#v, want inline result", req.Result)
	}

	if len(req.PeerTypes) != 2 {
		t.Fatalf("peer types = %#v, want user + broadcast", req.PeerTypes)
	}

	if _, ok := req.PeerTypes[0].(*tg.InlineQueryPeerTypePM); !ok {
		t.Fatalf("peer type 0 = %#v, want PM", req.PeerTypes[0])
	}

	if _, ok := req.PeerTypes[1].(*tg.InlineQueryPeerTypeBroadcast); !ok {
		t.Fatalf("peer type 1 = %#v, want broadcast", req.PeerTypes[1])
	}
}

func TestSavePreparedInlineMessageNilResult(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).SavePreparedInlineMessage(context.Background(), 99, nil); err == nil {
		t.Fatal("expected error for nil result")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
