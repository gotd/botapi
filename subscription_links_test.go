package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestCreateChatSubscriptionInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesExportChatInviteRequestTypeID, &tg.ChatInviteExported{
		Link:  "https://t.me/+sub",
		Title: "Subs",
		SubscriptionPricing: tg.StarsSubscriptionPricing{
			Period: 2592000,
			Amount: 100,
		},
	})

	b := newMockBot(inv)

	link, err := b.CreateChatSubscriptionInviteLink(
		context.Background(), tdlibChannel(50), 2592000, 100, WithInviteLinkName("Subs"),
	)
	if err != nil {
		t.Fatalf("CreateChatSubscriptionInviteLink: %v", err)
	}

	if link.SubscriptionPeriod != 2592000 || link.SubscriptionPrice != 100 || link.Name != "Subs" {
		t.Fatalf("link = %#v", link)
	}

	var req tg.MessagesExportChatInviteRequest

	inv.decode(t, tg.MessagesExportChatInviteRequestTypeID, &req)

	pricing, ok := req.GetSubscriptionPricing()
	if !ok || pricing.Period != 2592000 || pricing.Amount != 100 {
		t.Fatalf("pricing = %#v, ok=%v", pricing, ok)
	}

	if title, ok := req.GetTitle(); !ok || title != "Subs" {
		t.Fatalf("title = %q, ok=%v", title, ok)
	}
}

func TestEditChatSubscriptionInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditExportedChatInviteRequestTypeID, &tg.MessagesExportedChatInvite{
		Invite: &tg.ChatInviteExported{Link: "https://t.me/+sub", Title: "Renamed"},
	})

	b := newMockBot(inv)

	link, err := b.EditChatSubscriptionInviteLink(
		context.Background(), tdlibChannel(50), "https://t.me/+sub", WithInviteLinkName("Renamed"),
	)
	if err != nil {
		t.Fatalf("EditChatSubscriptionInviteLink: %v", err)
	}

	if link.Name != "Renamed" {
		t.Fatalf("link = %#v", link)
	}

	var req tg.MessagesEditExportedChatInviteRequest

	inv.decode(t, tg.MessagesEditExportedChatInviteRequestTypeID, &req)

	if req.Link != "https://t.me/+sub" {
		t.Fatalf("link = %q", req.Link)
	}
}
