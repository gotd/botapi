package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestAnswerChatJoinRequestQuery(t *testing.T) {
	for _, c := range []struct {
		result string
		want   tg.JoinChatBotResultClass
	}{
		{"approve", &tg.JoinChatBotResultApproved{}},
		{"decline", &tg.JoinChatBotResultDeclined{}},
		{"queue", &tg.JoinChatBotResultQueued{}},
	} {
		t.Run(c.result, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.BotsSetJoinChatResultsRequestTypeID, &tg.BoolTrue{})

			if err := newMockBot(inv).AnswerChatJoinRequestQuery(context.Background(), "123", c.result); err != nil {
				t.Fatalf("AnswerChatJoinRequestQuery: %v", err)
			}

			var req tg.BotsSetJoinChatResultsRequest

			inv.decode(t, tg.BotsSetJoinChatResultsRequestTypeID, &req)

			if req.QueryID != 123 {
				t.Fatalf("query id = %d", req.QueryID)
			}

			if req.Result.TypeID() != c.want.TypeID() {
				t.Fatalf("result = %#v, want %#v", req.Result, c.want)
			}
		})
	}
}

func TestAnswerChatJoinRequestQueryInvalid(t *testing.T) {
	inv := newMockInvoker()

	if err := newMockBot(inv).AnswerChatJoinRequestQuery(context.Background(), "123", "maybe"); err == nil {
		t.Fatal("expected error for invalid result")
	}

	if err := newMockBot(inv).AnswerChatJoinRequestQuery(context.Background(), "notanumber", "approve"); err == nil {
		t.Fatal("expected error for invalid query id")
	}
}

func TestSendChatJoinRequestWebApp(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetJoinChatResultsRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SendChatJoinRequestWebApp(context.Background(), "456", "https://example.com/app"); err != nil {
		t.Fatalf("SendChatJoinRequestWebApp: %v", err)
	}

	var req tg.BotsSetJoinChatResultsRequest

	inv.decode(t, tg.BotsSetJoinChatResultsRequestTypeID, &req)

	wv, ok := req.Result.(*tg.JoinChatBotResultWebView)
	if !ok || wv.URL != "https://example.com/app" {
		t.Fatalf("result = %#v", req.Result)
	}
}
