package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestConvertGiftToStars(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.BoolBox{Bool: &tg.BoolTrue{}})

	err := newMockBot(inv).ConvertGiftToStars(context.Background(), "bc1", OwnedGiftFromMessage(123))
	if err != nil {
		t.Fatalf("ConvertGiftToStars: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsConvertStarGiftRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	req, ok := wrapper.Query.(*tg.PaymentsConvertStarGiftRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	gift, ok := req.Stargift.(*tg.InputSavedStarGiftUser)
	if !ok || gift.MsgID != 123 {
		t.Fatalf("stargift = %#v", req.Stargift)
	}
}

func TestUpgradeGift(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.UpdatesBox{Updates: &tg.Updates{}})

	err := newMockBot(inv).UpgradeGift(context.Background(), "bc1", OwnedGiftFromMessage(7), true)
	if err != nil {
		t.Fatalf("UpgradeGift: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsUpgradeStarGiftRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.PaymentsUpgradeStarGiftRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if !req.KeepOriginalDetails {
		t.Fatalf("keep original details not set")
	}

	if g, ok := req.Stargift.(*tg.InputSavedStarGiftUser); !ok || g.MsgID != 7 {
		t.Fatalf("stargift = %#v", req.Stargift)
	}
}

func TestTransferGift(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.UpdatesBox{Updates: &tg.Updates{}})

	err := newMockBot(inv).TransferGift(context.Background(), "bc1", OwnedGiftFromSlug("rare-gift"), userRef(99, 5))
	if err != nil {
		t.Fatalf("TransferGift: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsTransferStarGiftRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.PaymentsTransferStarGiftRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if g, ok := req.Stargift.(*tg.InputSavedStarGiftSlug); !ok || g.Slug != "rare-gift" {
		t.Fatalf("stargift = %#v", req.Stargift)
	}

	if _, ok := req.ToID.(*tg.InputPeerUser); !ok {
		t.Fatalf("to id = %#v, want user", req.ToID)
	}
}

func TestConvertGiftInvalidID(t *testing.T) {
	inv := newMockInvoker()

	err := newMockBot(inv).ConvertGiftToStars(context.Background(), "bc1", "bogus")
	if err == nil {
		t.Fatal("expected error for invalid owned_gift_id")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
