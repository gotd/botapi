package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestVerifyUser(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetCustomVerificationRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).VerifyUser(context.Background(), 99, WithVerificationDescription("trusted")); err != nil {
		t.Fatalf("VerifyUser: %v", err)
	}

	var req tg.BotsSetCustomVerificationRequest

	inv.decode(t, tg.BotsSetCustomVerificationRequestTypeID, &req)

	if !req.GetEnabled() {
		t.Fatal("enabled should be true")
	}

	if desc, ok := req.GetCustomDescription(); !ok || desc != "trusted" {
		t.Fatalf("description = %q, ok=%v", desc, ok)
	}

	if _, ok := req.Peer.(*tg.InputPeerUser); !ok {
		t.Fatalf("peer = %#v, want user", req.Peer)
	}
}

func TestVerifyChat(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetCustomVerificationRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).VerifyChat(context.Background(), tdlibChannel(50)); err != nil {
		t.Fatalf("VerifyChat: %v", err)
	}

	var req tg.BotsSetCustomVerificationRequest

	inv.decode(t, tg.BotsSetCustomVerificationRequestTypeID, &req)

	if _, ok := req.Peer.(*tg.InputPeerChannel); !ok {
		t.Fatalf("peer = %#v, want channel", req.Peer)
	}
}

func TestRemoveUserVerification(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetCustomVerificationRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).RemoveUserVerification(context.Background(), 99); err != nil {
		t.Fatalf("RemoveUserVerification: %v", err)
	}

	var req tg.BotsSetCustomVerificationRequest

	inv.decode(t, tg.BotsSetCustomVerificationRequestTypeID, &req)

	if req.GetEnabled() {
		t.Fatal("enabled should be false")
	}
}

func TestRemoveChatVerification(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetCustomVerificationRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).RemoveChatVerification(context.Background(), tdlibChannel(50)); err != nil {
		t.Fatalf("RemoveChatVerification: %v", err)
	}

	if !inv.called(tg.BotsSetCustomVerificationRequestTypeID) {
		t.Fatal("verification not called")
	}
}
