package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSavePreparedKeyboardButtonUsers(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsRequestWebViewButtonRequestTypeID, &tg.BotsRequestedButton{WebappReqID: "btn1"})

	isBot := false
	button := KeyboardButton{
		Text:         "Pick a user",
		RequestUsers: &KeyboardButtonRequestUsers{RequestID: 7, UserIsBot: &isBot, MaxQuantity: 3},
	}

	got, err := newMockBot(inv).SavePreparedKeyboardButton(context.Background(), 99, button)
	if err != nil {
		t.Fatalf("SavePreparedKeyboardButton: %v", err)
	}

	if got.ID != "btn1" {
		t.Fatalf("id = %q", got.ID)
	}

	var req tg.BotsRequestWebViewButtonRequest

	inv.decode(t, tg.BotsRequestWebViewButtonRequestTypeID, &req)

	rp, ok := req.Button.(*tg.KeyboardButtonRequestPeer)
	if !ok || rp.ButtonID != 7 || rp.MaxQuantity != 3 {
		t.Fatalf("button = %#v", req.Button)
	}

	ut, ok := rp.PeerType.(*tg.RequestPeerTypeUser)
	if !ok {
		t.Fatalf("peer type = %#v, want user", rp.PeerType)
	}

	if bot, ok := ut.GetBot(); !ok || bot {
		t.Fatalf("bot = %v ok=%v, want false", bot, ok)
	}
}

func TestSavePreparedKeyboardButtonChat(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsRequestWebViewButtonRequestTypeID, &tg.BotsRequestedButton{WebappReqID: "btn2"})

	button := KeyboardButton{
		Text:        "Pick a channel",
		RequestChat: &KeyboardButtonRequestChat{RequestID: 9, ChatIsChannel: true, ChatIsCreated: true},
	}

	if _, err := newMockBot(inv).SavePreparedKeyboardButton(context.Background(), 99, button); err != nil {
		t.Fatalf("SavePreparedKeyboardButton: %v", err)
	}

	var req tg.BotsRequestWebViewButtonRequest

	inv.decode(t, tg.BotsRequestWebViewButtonRequestTypeID, &req)

	rp := req.Button.(*tg.KeyboardButtonRequestPeer)

	bc, ok := rp.PeerType.(*tg.RequestPeerTypeBroadcast)
	if !ok || !bc.Creator {
		t.Fatalf("peer type = %#v, want broadcast creator", rp.PeerType)
	}
}

func TestSavePreparedKeyboardButtonInvalid(t *testing.T) {
	inv := newMockInvoker()

	if _, err := newMockBot(inv).SavePreparedKeyboardButton(context.Background(), 99, KeyboardButton{Text: "x"}); err == nil {
		t.Fatal("expected error for non-request button")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
