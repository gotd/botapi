package botapi

import (
	"context"
	"testing"
)

// TestInviteLinkErrorsOnPrivate covers the resolveInviteLinks error branch of
// every invite-link method: a private (user) chat has no invite links.
func TestInviteLinkErrorsOnPrivate(t *testing.T) {
	b := newMockBot(newMockInvoker())
	ctx := context.Background()
	chat := userRef(10, 20)

	if _, err := b.ExportChatInviteLink(ctx, chat); err == nil {
		t.Fatal("ExportChatInviteLink on private chat should fail")
	}
	if _, err := b.CreateChatInviteLink(ctx, chat); err == nil {
		t.Fatal("CreateChatInviteLink on private chat should fail")
	}
	if _, err := b.EditChatInviteLink(ctx, chat, "https://t.me/+x"); err == nil {
		t.Fatal("EditChatInviteLink on private chat should fail")
	}
	if _, err := b.RevokeChatInviteLink(ctx, chat, "https://t.me/+x"); err == nil {
		t.Fatal("RevokeChatInviteLink on private chat should fail")
	}
}
