package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
)

// testTgErr is a canned Telegram error used to drive asAPIError branches.
func testTgErr() *tgerr.Error {
	return &tgerr.Error{Code: 400, Type: "TEST_ERROR", Message: "TEST_ERROR"}
}

// botFailing returns a Bot whose given RPC fails, leaving the default
// resolution handlers intact.
func botFailing(id uint32) *Bot {
	inv := newMockInvoker()
	inv.fail(id, testTgErr())

	return newMockBot(inv)
}

// TestChatAdminMethodErrors covers the RPC-error branch of the direct-RPC chat
// administration methods.
func TestChatAdminMethodErrors(t *testing.T) {
	ctx := context.Background()
	ch := tdlibChannel(123)

	checks := []struct {
		name string
		call func(b *Bot) error
		rpc  uint32
	}{
		{"pin", func(b *Bot) error { return b.PinChatMessage(ctx, userRef(10, 20), 5) }, tg.MessagesUpdatePinnedMessageRequestTypeID},
		{"unpin", func(b *Bot) error { return b.UnpinChatMessage(ctx, userRef(10, 20), 5) }, tg.MessagesUpdatePinnedMessageRequestTypeID},
		{"unpin-all", func(b *Bot) error { return b.UnpinAllChatMessages(ctx, userRef(10, 20)) }, tg.MessagesUnpinAllMessagesRequestTypeID},
		{"permissions", func(b *Bot) error { return b.SetChatPermissions(ctx, userRef(10, 20), ChatPermissions{}) }, tg.MessagesEditChatDefaultBannedRightsRequestTypeID},
		{"set-sticker-set", func(b *Bot) error { return b.SetChatStickerSet(ctx, ch, "pack") }, tg.ChannelsSetStickersRequestTypeID},
		{"del-sticker-set", func(b *Bot) error { return b.DeleteChatStickerSet(ctx, ch) }, tg.ChannelsSetStickersRequestTypeID},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if err := c.call(botFailing(c.rpc)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

// TestChatMemberQueryErrors covers the RPC-error branch of the participant
// query methods.
func TestChatMemberQueryErrors(t *testing.T) {
	ctx := context.Background()
	ch := tdlibChannel(123)

	if _, err := botFailing(tg.ChannelsGetParticipantRequestTypeID).GetChatMember(ctx, ch, 10); err == nil {
		t.Fatal("GetChatMember: expected error")
	}

	if _, err := botFailing(tg.ChannelsGetParticipantsRequestTypeID).GetChatMemberCount(ctx, ch); err == nil {
		t.Fatal("GetChatMemberCount: expected error")
	}

	if _, err := botFailing(tg.ChannelsGetParticipantsRequestTypeID).GetChatAdministrators(ctx, ch); err == nil {
		t.Fatal("GetChatAdministrators: expected error")
	}
}

// TestChatAdminRejectsPrivate covers errNotInPrivateChat: the shared-admin
// methods are not available on a private (user) chat.
func TestChatAdminRejectsPrivate(t *testing.T) {
	ctx := context.Background()
	b := newMockBot(newMockInvoker())

	if err := b.SetChatTitle(ctx, userRef(10, 20), "t"); err == nil {
		t.Fatal("SetChatTitle on private chat should fail")
	}

	if err := b.SetChatDescription(ctx, userRef(10, 20), "d"); err == nil {
		t.Fatal("SetChatDescription on private chat should fail")
	}

	if err := b.LeaveChat(ctx, userRef(10, 20)); err == nil {
		t.Fatal("LeaveChat on private chat should fail")
	}
}

// TestResolveChannelRejectsUser covers the resolveChannel error path: a numeric
// user id is not a channel.
func TestResolveChannelRejectsUser(t *testing.T) {
	ctx := context.Background()
	b := newMockBot(newMockInvoker())

	if _, err := b.GetChatMemberCount(ctx, ID(10)); err == nil {
		t.Fatal("GetChatMemberCount on a user id should fail")
	}

	if err := b.SetChatStickerSet(ctx, ID(10), "p"); err == nil {
		t.Fatal("SetChatStickerSet on a user id should fail")
	}
}
