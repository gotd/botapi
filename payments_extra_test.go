package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestCreateInvoiceLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsExportInvoiceRequestTypeID, &tg.PaymentsExportedInvoice{URL: "https://t.me/invoice/abc"})

	params := InvoiceParams{
		Title:       "Item",
		Description: "An item",
		Payload:     "payload",
		Currency:    "XTR",
		Prices:      []LabeledPrice{{Label: "Item", Amount: 100}},
	}

	link, err := newMockBot(inv).CreateInvoiceLink(context.Background(), params)
	if err != nil {
		t.Fatalf("CreateInvoiceLink: %v", err)
	}

	if link != "https://t.me/invoice/abc" {
		t.Fatalf("link = %q", link)
	}

	var req tg.PaymentsExportInvoiceRequest

	inv.decode(t, tg.PaymentsExportInvoiceRequestTypeID, &req)

	media, ok := req.InvoiceMedia.(*tg.InputMediaInvoice)
	if !ok || media.Title != "Item" {
		t.Fatalf("media = %#v", req.InvoiceMedia)
	}
}

func TestRefundStarPayment(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsRefundStarsChargeRequestTypeID, okUpdates())

	if err := newMockBot(inv).RefundStarPayment(context.Background(), 99, "charge_123"); err != nil {
		t.Fatalf("RefundStarPayment: %v", err)
	}

	var req tg.PaymentsRefundStarsChargeRequest

	inv.decode(t, tg.PaymentsRefundStarsChargeRequestTypeID, &req)

	if req.ChargeID != "charge_123" {
		t.Fatalf("charge id = %q", req.ChargeID)
	}

	if _, ok := req.UserID.(*tg.InputUser); !ok {
		t.Fatalf("user = %#v", req.UserID)
	}
}
