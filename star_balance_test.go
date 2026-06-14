package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetMyStarBalance(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetStarsStatusRequestTypeID, &tg.PaymentsStarsStatus{
		Balance: &tg.StarsAmount{Amount: 1500, Nanos: 250},
	})

	got, err := newMockBot(inv).GetMyStarBalance(context.Background())
	if err != nil {
		t.Fatalf("GetMyStarBalance: %v", err)
	}

	if got.Amount != 1500 || got.NanostarAmount != 250 {
		t.Fatalf("balance = %#v", got)
	}

	var req tg.PaymentsGetStarsStatusRequest

	inv.decode(t, tg.PaymentsGetStarsStatusRequestTypeID, &req)

	if _, ok := req.Peer.(*tg.InputPeerSelf); !ok {
		t.Fatalf("peer = %#v, want self", req.Peer)
	}
}
