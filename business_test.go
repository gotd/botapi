package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetBusinessConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.AccountGetBotBusinessConnectionRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateBotBusinessConnect{Connection: tg.BotBusinessConnection{
				ConnectionID: "bc1",
				UserID:       10,
				Date:         123,
				Rights:       tg.BusinessBotRights{Reply: true, EditBio: true},
			}},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 1, FirstName: "Biz"}},
	})

	conn, err := newMockBot(inv).GetBusinessConnection(context.Background(), "bc1")
	if err != nil {
		t.Fatalf("GetBusinessConnection: %v", err)
	}

	if conn.ID != "bc1" || conn.UserChatID != 10 || conn.Date != 123 || !conn.IsEnabled {
		t.Fatalf("conn = %#v", conn)
	}

	if conn.User.FirstName != "Biz" {
		t.Fatalf("user = %#v", conn.User)
	}

	if conn.Rights == nil || !conn.Rights.CanReply || !conn.Rights.CanEditBio {
		t.Fatalf("rights = %#v", conn.Rights)
	}
}

func TestSetBusinessAccountName(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.User{ID: 10})

	if err := newMockBot(inv).SetBusinessAccountName(context.Background(), "bc1", "First", "Last"); err != nil {
		t.Fatalf("SetBusinessAccountName: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.AccountUpdateProfileRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	profile, ok := wrapper.Query.(*tg.AccountUpdateProfileRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if first, _ := profile.GetFirstName(); first != "First" {
		t.Fatalf("first name = %q", first)
	}

	if last, _ := profile.GetLastName(); last != "Last" {
		t.Fatalf("last name = %q", last)
	}
}

func TestSetBusinessAccountBio(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.User{ID: 10})

	if err := newMockBot(inv).SetBusinessAccountBio(context.Background(), "bc1", "my bio"); err != nil {
		t.Fatalf("SetBusinessAccountBio: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.AccountUpdateProfileRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	profile, ok := wrapper.Query.(*tg.AccountUpdateProfileRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if about, _ := profile.GetAbout(); about != "my bio" {
		t.Fatalf("about = %q", about)
	}
}

func TestSetBusinessAccountUsername(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.User{ID: 10})

	if err := newMockBot(inv).SetBusinessAccountUsername(context.Background(), "bc1", "biz_user"); err != nil {
		t.Fatalf("SetBusinessAccountUsername: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.AccountUpdateUsernameRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	username, ok := wrapper.Query.(*tg.AccountUpdateUsernameRequest)
	if !ok || username.Username != "biz_user" {
		t.Fatalf("query = %#v", wrapper.Query)
	}
}

func TestGetBusinessAccountStarBalance(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PaymentsStarsStatus{
		Balance: &tg.StarsAmount{Amount: 42, Nanos: 5},
	})

	got, err := newMockBot(inv).GetBusinessAccountStarBalance(context.Background(), "bc1")
	if err != nil {
		t.Fatalf("GetBusinessAccountStarBalance: %v", err)
	}

	if got.Amount != 42 || got.NanostarAmount != 5 {
		t.Fatalf("balance = %#v", got)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsGetStarsStatusRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	status, ok := wrapper.Query.(*tg.PaymentsGetStarsStatusRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := status.Peer.(*tg.InputPeerSelf); !ok {
		t.Fatalf("peer = %#v, want self", status.Peer)
	}
}

func TestDeleteBusinessMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.MessagesAffectedMessages{})

	if err := newMockBot(inv).DeleteBusinessMessages(context.Background(), "bc1", []int{5, 6}); err != nil {
		t.Fatalf("DeleteBusinessMessages: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesDeleteMessagesRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	del, ok := wrapper.Query.(*tg.MessagesDeleteMessagesRequest)
	if !ok || len(del.ID) != 2 || !del.Revoke {
		t.Fatalf("query = %#v", wrapper.Query)
	}
}

func TestDeleteBusinessMessagesEmpty(t *testing.T) {
	inv := newMockInvoker()

	if err := newMockBot(inv).DeleteBusinessMessages(context.Background(), "bc1", nil); err != nil {
		t.Fatalf("DeleteBusinessMessages: %v", err)
	}

	if inv.count() != 0 {
		t.Fatal("empty delete should make no RPC")
	}
}
