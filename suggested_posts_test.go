package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestApproveSuggestedPost(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesToggleSuggestedPostApprovalRequestTypeID, okUpdates())

	if err := newMockBot(inv).ApproveSuggestedPost(context.Background(), tdlibChannel(50), 12,
		WithSuggestedPostSendDate(1700000000)); err != nil {
		t.Fatalf("ApproveSuggestedPost: %v", err)
	}

	var req tg.MessagesToggleSuggestedPostApprovalRequest

	inv.decode(t, tg.MessagesToggleSuggestedPostApprovalRequestTypeID, &req)

	if req.Reject || req.MsgID != 12 {
		t.Fatalf("req = %#v", req)
	}

	if date, ok := req.GetScheduleDate(); !ok || date != 1700000000 {
		t.Fatalf("schedule date = %d ok=%v", date, ok)
	}
}

func TestDeclineSuggestedPost(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesToggleSuggestedPostApprovalRequestTypeID, okUpdates())

	if err := newMockBot(inv).DeclineSuggestedPost(context.Background(), tdlibChannel(50), 12, "not now"); err != nil {
		t.Fatalf("DeclineSuggestedPost: %v", err)
	}

	var req tg.MessagesToggleSuggestedPostApprovalRequest

	inv.decode(t, tg.MessagesToggleSuggestedPostApprovalRequestTypeID, &req)

	if !req.Reject || req.MsgID != 12 {
		t.Fatalf("req = %#v", req)
	}

	if comment, ok := req.GetRejectComment(); !ok || comment != "not now" {
		t.Fatalf("comment = %q ok=%v", comment, ok)
	}
}
