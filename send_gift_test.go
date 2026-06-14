package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendGift(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetPaymentFormRequestTypeID, &tg.PaymentsPaymentFormStarGift{
		FormID:  999,
		Invoice: tg.Invoice{Currency: "XTR"},
	})
	inv.reply(tg.PaymentsSendStarsFormRequestTypeID, &tg.PaymentsPaymentResult{
		Updates: &tg.Updates{},
	})

	bot := newMockBot(inv)

	err := bot.SendGift(context.Background(), userRef(42, 7), "100",
		WithGiftText("thanks!"), WithGiftPayForUpgrade())
	if err != nil {
		t.Fatalf("SendGift: %v", err)
	}

	var form tg.PaymentsGetPaymentFormRequest

	inv.decode(t, tg.PaymentsGetPaymentFormRequestTypeID, &form)

	invoice, ok := form.Invoice.(*tg.InputInvoiceStarGift)
	if !ok {
		t.Fatalf("invoice = %#v, want star gift", form.Invoice)
	}

	if invoice.GiftID != 100 {
		t.Fatalf("gift id = %d, want 100", invoice.GiftID)
	}

	if !invoice.IncludeUpgrade {
		t.Fatalf("include upgrade not set")
	}

	if invoice.Message.Text != "thanks!" {
		t.Fatalf("message = %q", invoice.Message.Text)
	}

	if _, ok := invoice.Peer.(*tg.InputPeerUser); !ok {
		t.Fatalf("peer = %#v, want user", invoice.Peer)
	}

	var send tg.PaymentsSendStarsFormRequest

	inv.decode(t, tg.PaymentsSendStarsFormRequestTypeID, &send)

	if send.FormID != 999 {
		t.Fatalf("form id = %d, want 999", send.FormID)
	}

	if _, ok := send.Invoice.(*tg.InputInvoiceStarGift); !ok {
		t.Fatalf("send invoice = %#v, want star gift", send.Invoice)
	}
}

func TestGiftPremiumSubscription(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetPaymentFormRequestTypeID, &tg.PaymentsPaymentFormStars{
		FormID:  555,
		Invoice: tg.Invoice{Currency: "XTR"},
	})
	inv.reply(tg.PaymentsSendStarsFormRequestTypeID, &tg.PaymentsPaymentResult{
		Updates: &tg.Updates{},
	})

	bot := newMockBot(inv)

	err := bot.GiftPremiumSubscription(context.Background(), 42, 6, WithGiftText("enjoy"))
	if err != nil {
		t.Fatalf("GiftPremiumSubscription: %v", err)
	}

	var form tg.PaymentsGetPaymentFormRequest

	inv.decode(t, tg.PaymentsGetPaymentFormRequestTypeID, &form)

	invoice, ok := form.Invoice.(*tg.InputInvoicePremiumGiftStars)
	if !ok {
		t.Fatalf("invoice = %#v, want premium gift stars", form.Invoice)
	}

	if invoice.Months != 6 {
		t.Fatalf("months = %d, want 6", invoice.Months)
	}

	if invoice.Message.Text != "enjoy" {
		t.Fatalf("message = %q", invoice.Message.Text)
	}

	if _, ok := invoice.UserID.(*tg.InputUser); !ok {
		t.Fatalf("user = %#v, want input user", invoice.UserID)
	}

	var send tg.PaymentsSendStarsFormRequest

	inv.decode(t, tg.PaymentsSendStarsFormRequestTypeID, &send)

	if send.FormID != 555 {
		t.Fatalf("form id = %d, want 555", send.FormID)
	}
}

func TestSendGiftInvalidID(t *testing.T) {
	inv := newMockInvoker()

	err := newMockBot(inv).SendGift(context.Background(), userRef(42, 7), "not-a-number")
	if err == nil {
		t.Fatal("expected error for invalid gift id")
	}

	if inv.count() != 0 {
		t.Fatalf("made %d RPC calls, want 0", inv.count())
	}
}
