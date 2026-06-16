package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

func TestReadBusinessMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.MessagesAffectedMessages{})

	if err := newMockBot(inv).ReadBusinessMessage(context.Background(), "bc1", userRef(10, 20), 77); err != nil {
		t.Fatalf("ReadBusinessMessage: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesReadHistoryRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	rh, ok := wrapper.Query.(*tg.MessagesReadHistoryRequest)
	if !ok || rh.MaxID != 77 {
		t.Fatalf("query = %#v", wrapper.Query)
	}
}

func TestTransferBusinessAccountStars(t *testing.T) {
	inv := newMockInvoker()

	// The transfer makes two business-wrapped calls under the same outer type:
	// getPaymentForm then sendStarsForm. Respond by call order.
	calls := 0

	inv.handle(tg.InvokeWithBusinessConnectionRequestTypeID, func(*bin.Buffer) (bin.Encoder, error) {
		calls++
		if calls == 1 {
			return &tg.PaymentsPaymentFormStars{FormID: 555, Invoice: tg.Invoice{Currency: "XTR"}}, nil
		}

		return &tg.PaymentsPaymentResult{Updates: &tg.Updates{}}, nil
	})

	if err := newMockBot(inv).TransferBusinessAccountStars(context.Background(), "bc2", 100); err != nil {
		t.Fatalf("TransferBusinessAccountStars: %v", err)
	}

	// The last wrapped call is sendStarsForm.
	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsSendStarsFormRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc2" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	send, ok := wrapper.Query.(*tg.PaymentsSendStarsFormRequest)
	if !ok || send.FormID != 555 {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	invoice, ok := send.Invoice.(*tg.InputInvoiceBusinessBotTransferStars)
	if !ok || invoice.Stars != 100 {
		t.Fatalf("invoice = %#v", send.Invoice)
	}

	if _, ok := invoice.Bot.(*tg.InputUserSelf); !ok {
		t.Fatalf("bot = %#v, want self", invoice.Bot)
	}
}
