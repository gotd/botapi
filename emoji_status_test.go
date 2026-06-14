package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetUserEmojiStatus(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsUpdateUserEmojiStatusRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SetUserEmojiStatus(context.Background(), 99, "555", WithEmojiStatusExpiration(1000)); err != nil {
		t.Fatalf("SetUserEmojiStatus: %v", err)
	}

	var req tg.BotsUpdateUserEmojiStatusRequest

	inv.decode(t, tg.BotsUpdateUserEmojiStatusRequestTypeID, &req)

	status, ok := req.EmojiStatus.(*tg.EmojiStatus)
	if !ok || status.DocumentID != 555 {
		t.Fatalf("status = %#v", req.EmojiStatus)
	}

	if until, ok := status.GetUntil(); !ok || until != 1000 {
		t.Fatalf("until = %d, ok=%v", until, ok)
	}
}

func TestSetUserEmojiStatusClear(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsUpdateUserEmojiStatusRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SetUserEmojiStatus(context.Background(), 99, ""); err != nil {
		t.Fatalf("SetUserEmojiStatus: %v", err)
	}

	var req tg.BotsUpdateUserEmojiStatusRequest

	inv.decode(t, tg.BotsUpdateUserEmojiStatusRequestTypeID, &req)

	if _, ok := req.EmojiStatus.(*tg.EmojiStatusEmpty); !ok {
		t.Fatalf("status = %#v, want empty", req.EmojiStatus)
	}
}
