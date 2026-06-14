package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendInvoice(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, sendMediaOK())
	b := newMockBot(inv)

	params := InvoiceParams{
		Title:       "Item",
		Description: "An item",
		Payload:     "payload",
		Currency:    "USD",
		Prices:      []LabeledPrice{{Label: "Item", Amount: 1000}},
	}
	if _, err := b.SendInvoice(context.Background(), userRef(10, 20), params); err != nil {
		t.Fatalf("SendInvoice: %v", err)
	}
	var req tg.MessagesSendMediaRequest
	inv.decode(t, tg.MessagesSendMediaRequestTypeID, &req)
	if _, ok := req.Media.(*tg.InputMediaInvoice); !ok {
		t.Fatalf("media = %#v", req.Media)
	}
}
