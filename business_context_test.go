package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestBusinessContextConnectionID(t *testing.T) {
	bc := newMockBot(newMockInvoker()).Business("bc9")

	if bc.ConnectionID() != "bc9" {
		t.Fatalf("connection id = %q", bc.ConnectionID())
	}
}

func TestBusinessContextSetNameScopesConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.User{ID: 10})

	if err := newMockBot(inv).Business("bc9").SetName(context.Background(), "First", "Last"); err != nil {
		t.Fatalf("SetName: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.AccountUpdateProfileRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc9" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	profile, ok := wrapper.Query.(*tg.AccountUpdateProfileRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if first, _ := profile.GetFirstName(); first != "First" {
		t.Fatalf("first name = %q", first)
	}
}

func TestBusinessContextStarBalance(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PaymentsStarsStatus{
		Balance: &tg.StarsAmount{Amount: 7},
	})

	got, err := newMockBot(inv).Business("bc9").StarBalance(context.Background())
	if err != nil {
		t.Fatalf("StarBalance: %v", err)
	}

	if got.Amount != 7 {
		t.Fatalf("balance = %#v", got)
	}
}

func TestBusinessContextDeleteMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.MessagesAffectedMessages{})

	if err := newMockBot(inv).Business("bc9").DeleteMessages(context.Background(), []int{1, 2}); err != nil {
		t.Fatalf("DeleteMessages: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.MessagesDeleteMessagesRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	del, ok := wrapper.Query.(*tg.MessagesDeleteMessagesRequest)
	if !ok || len(del.ID) != 2 {
		t.Fatalf("query = %#v", wrapper.Query)
	}
}

func TestBusinessContextConnection(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.AccountGetBotBusinessConnectionRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateBotBusinessConnect{Connection: tg.BotBusinessConnection{ConnectionID: "bc9", UserID: 10}},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, FirstName: "Biz"}},
	})

	conn, err := newMockBot(inv).Business("bc9").Connection(context.Background())
	if err != nil {
		t.Fatalf("Connection: %v", err)
	}

	if conn.ID != "bc9" {
		t.Fatalf("conn = %#v", conn)
	}

	var req tg.AccountGetBotBusinessConnectionRequest

	inv.decode(t, tg.AccountGetBotBusinessConnectionRequestTypeID, &req)

	if req.ConnectionID != "bc9" {
		t.Fatalf("requested connection id = %q", req.ConnectionID)
	}
}

func TestBusinessContextRemoveProfilePhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	if err := newMockBot(inv).Business("bc9").RemoveProfilePhoto(context.Background(), false); err != nil {
		t.Fatalf("RemoveProfilePhoto: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PhotosUpdateProfilePhotoRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc9" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}
}
