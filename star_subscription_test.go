package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestEditUserStarSubscription(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsChangeStarsSubscriptionRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).EditUserStarSubscription(context.Background(), 99, "charge123", true); err != nil {
		t.Fatalf("EditUserStarSubscription: %v", err)
	}

	var req tg.PaymentsChangeStarsSubscriptionRequest

	inv.decode(t, tg.PaymentsChangeStarsSubscriptionRequestTypeID, &req)

	if req.SubscriptionID != "charge123" {
		t.Fatalf("subscription id = %q", req.SubscriptionID)
	}

	if canceled, ok := req.GetCanceled(); !ok || !canceled {
		t.Fatalf("canceled = %v ok=%v", canceled, ok)
	}

	if _, ok := req.Peer.(*tg.InputPeerUser); !ok {
		t.Fatalf("peer = %#v, want user", req.Peer)
	}
}
