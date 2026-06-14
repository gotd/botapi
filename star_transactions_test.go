package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetStarTransactions(t *testing.T) {
	incoming := tg.StarsTransaction{
		ID:     "in-1",
		Amount: &tg.StarsAmount{Amount: 100, Nanos: 5},
		Date:   1000,
		Peer:   &tg.StarsTransactionPeer{Peer: &tg.PeerUser{UserID: 42}},
	}
	incoming.SetBotPayload([]byte("payload"))

	outgoing := tg.StarsTransaction{
		ID:     "out-1",
		Amount: &tg.StarsAmount{Amount: -250, Nanos: 0},
		Date:   2000,
		Peer:   &tg.StarsTransactionPeerFragment{},
	}
	outgoing.SetTransactionDate(2100)
	outgoing.SetTransactionURL("https://fragment.example/tx")

	api := tg.StarsTransaction{
		ID:     "api-1",
		Amount: &tg.StarsAmount{Amount: -7, Nanos: 0},
		Date:   3000,
		Peer:   &tg.StarsTransactionPeerAPI{},
	}
	api.SetFloodskipNumber(7)

	status := &tg.PaymentsStarsStatus{
		Balance: &tg.StarsAmount{Amount: 1},
		Users: []tg.UserClass{
			&tg.User{ID: 42, FirstName: "Buyer", Username: "buyer"},
		},
	}
	status.SetHistory([]tg.StarsTransaction{incoming, outgoing, api})

	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetStarsTransactionsRequestTypeID, status)

	got, err := newMockBot(inv).GetStarTransactions(context.Background(), 5, 50)
	if err != nil {
		t.Fatalf("GetStarTransactions: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}

	// Incoming user transaction: positive amount, partner is the source.
	in := got[0]
	if in.ID != "in-1" || in.Amount != 100 || in.NanostarAmount != 5 || in.Date != 1000 {
		t.Fatalf("incoming = %#v", in)
	}

	if in.Receiver != nil {
		t.Fatalf("incoming receiver = %#v, want nil", in.Receiver)
	}

	user, ok := in.Source.(TransactionPartnerUser)
	if !ok {
		t.Fatalf("incoming source = %#v, want user", in.Source)
	}

	if user.User.ID != 42 || user.User.FirstName != "Buyer" {
		t.Fatalf("partner user = %#v", user.User)
	}

	if user.TransactionType != transactionTypeInvoicePayment {
		t.Fatalf("transaction_type = %q", user.TransactionType)
	}

	if user.InvoicePayload != "payload" {
		t.Fatalf("invoice payload = %q", user.InvoicePayload)
	}

	// Outgoing fragment withdrawal: amount made positive, partner is receiver.
	out := got[1]
	if out.Amount != 250 || out.Source != nil {
		t.Fatalf("outgoing = %#v", out)
	}

	frag, ok := out.Receiver.(TransactionPartnerFragment)
	if !ok {
		t.Fatalf("outgoing receiver = %#v, want fragment", out.Receiver)
	}

	state, ok := frag.WithdrawalState.(RevenueWithdrawalStateSucceeded)
	if !ok {
		t.Fatalf("withdrawal state = %#v, want succeeded", frag.WithdrawalState)
	}

	if state.Date != 2100 || state.URL != "https://fragment.example/tx" {
		t.Fatalf("succeeded state = %#v", state)
	}

	// Telegram API usage.
	apiTx := got[2]

	tapi, ok := apiTx.Receiver.(TransactionPartnerTelegramApi)
	if !ok {
		t.Fatalf("api receiver = %#v, want telegram_api", apiTx.Receiver)
	}

	if tapi.RequestCount != 7 || apiTx.Amount != 7 {
		t.Fatalf("api tx = %#v / %#v", apiTx, tapi)
	}

	var req tg.PaymentsGetStarsTransactionsRequest

	inv.decode(t, tg.PaymentsGetStarsTransactionsRequestTypeID, &req)

	if _, ok := req.Peer.(*tg.InputPeerSelf); !ok {
		t.Fatalf("peer = %#v, want self", req.Peer)
	}

	if req.Limit != 50 {
		t.Fatalf("limit = %d, want 50", req.Limit)
	}

	if req.Offset != "5" {
		t.Fatalf("offset = %q, want 5", req.Offset)
	}
}

func TestGetStarTransactionsDefaultLimit(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetStarsTransactionsRequestTypeID, &tg.PaymentsStarsStatus{
		Balance: &tg.StarsAmount{Amount: 0},
	})

	if _, err := newMockBot(inv).GetStarTransactions(context.Background(), 0, 0); err != nil {
		t.Fatalf("GetStarTransactions: %v", err)
	}

	var req tg.PaymentsGetStarsTransactionsRequest

	inv.decode(t, tg.PaymentsGetStarsTransactionsRequestTypeID, &req)

	if req.Limit != 100 {
		t.Fatalf("limit = %d, want default 100", req.Limit)
	}

	if req.Offset != "" {
		t.Fatalf("offset = %q, want empty", req.Offset)
	}
}
