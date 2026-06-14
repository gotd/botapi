package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func exportedInvite() *tg.ChatInviteExported {
	return &tg.ChatInviteExported{Link: "https://t.me/+abc", AdminID: 1}
}

func TestExportChatInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesExportChatInviteRequestTypeID, exportedInvite())

	b := newMockBot(inv)

	link, err := b.ExportChatInviteLink(context.Background(), tdlibChannel(50))
	if err != nil {
		t.Fatalf("ExportChatInviteLink: %v", err)
	}

	if link != "https://t.me/+abc" {
		t.Fatalf("link = %q", link)
	}
}

func TestCreateChatInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesExportChatInviteRequestTypeID, exportedInvite())

	b := newMockBot(inv)

	link, err := b.CreateChatInviteLink(context.Background(), tdlibChannel(50),
		WithInviteLinkName("promo"), WithInviteLinkMemberLimit(100))
	if err != nil {
		t.Fatalf("CreateChatInviteLink: %v", err)
	}

	if link.InviteLink != "https://t.me/+abc" {
		t.Fatalf("link = %#v", link)
	}

	var req tg.MessagesExportChatInviteRequest

	inv.decode(t, tg.MessagesExportChatInviteRequestTypeID, &req)

	if req.Title != "promo" || req.UsageLimit != 100 {
		t.Fatalf("req = %#v", req)
	}
}

func TestRevokeChatInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditExportedChatInviteRequestTypeID, &tg.MessagesExportedChatInvite{
		Invite: &tg.ChatInviteExported{Link: "https://t.me/+abc", AdminID: 1, Revoked: true},
	})

	b := newMockBot(inv)

	link, err := b.RevokeChatInviteLink(context.Background(), tdlibChannel(50), "https://t.me/+abc")
	if err != nil {
		t.Fatalf("RevokeChatInviteLink: %v", err)
	}

	if !link.IsRevoked {
		t.Fatal("link should be revoked")
	}
}

func TestEditChatInviteLink(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditExportedChatInviteRequestTypeID, &tg.MessagesExportedChatInvite{
		Invite: &tg.ChatInviteExported{Link: "https://t.me/+abc", AdminID: 1, Title: "renamed"},
	})

	b := newMockBot(inv)

	link, err := b.EditChatInviteLink(context.Background(), tdlibChannel(50), "https://t.me/+abc",
		WithInviteLinkName("renamed"))
	if err != nil {
		t.Fatalf("EditChatInviteLink: %v", err)
	}

	if link.Name != "renamed" {
		t.Fatalf("name = %q", link.Name)
	}
}
