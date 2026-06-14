package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestBanChatMember(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditBannedRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.BanChatMember(context.Background(), tdlibChannel(50), 99); err != nil {
		t.Fatalf("BanChatMember: %v", err)
	}

	var req tg.ChannelsEditBannedRequest

	inv.decode(t, tg.ChannelsEditBannedRequestTypeID, &req)

	ch, ok := req.Channel.(*tg.InputChannel)
	if !ok || ch.ChannelID != 50 {
		t.Fatalf("channel = %#v", req.Channel)
	}

	if !req.BannedRights.ViewMessages {
		t.Fatal("ban should set ViewMessages")
	}
}

func TestUnbanChatMember(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditBannedRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.UnbanChatMember(context.Background(), tdlibChannel(50), 99); err != nil {
		t.Fatalf("UnbanChatMember: %v", err)
	}

	var req tg.ChannelsEditBannedRequest

	inv.decode(t, tg.ChannelsEditBannedRequestTypeID, &req)

	if req.BannedRights.ViewMessages {
		t.Fatal("unban should clear ViewMessages")
	}
}

func TestPromoteChatMember(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditAdminRequestTypeID, okUpdates())

	b := newMockBot(inv)

	rights := ChatAdminRights{CanDeleteMessages: true, CanPinMessages: true}
	if err := b.PromoteChatMember(context.Background(), tdlibChannel(50), 99, rights); err != nil {
		t.Fatalf("PromoteChatMember: %v", err)
	}

	var req tg.ChannelsEditAdminRequest

	inv.decode(t, tg.ChannelsEditAdminRequestTypeID, &req)

	if !req.AdminRights.DeleteMessages || !req.AdminRights.PinMessages {
		t.Fatalf("admin rights = %#v", req.AdminRights)
	}
}

func TestRestrictChatMember(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditBannedRequestTypeID, okUpdates())

	b := newMockBot(inv)

	perms := ChatPermissions{CanSendMessages: false, CanSendPolls: false}
	if err := b.RestrictChatMember(context.Background(), tdlibChannel(50), 99, perms, 0); err != nil {
		t.Fatalf("RestrictChatMember: %v", err)
	}

	if !inv.called(tg.ChannelsEditBannedRequestTypeID) {
		t.Fatal("expected channels.editBanned")
	}
}

func TestGetChatMember(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsGetParticipantRequestTypeID, &tg.ChannelsChannelParticipant{
		Participant: &tg.ChannelParticipantAdmin{UserID: 99},
		Users:       []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})

	b := newMockBot(inv)

	m, err := b.GetChatMember(context.Background(), tdlibChannel(50), 99)
	if err != nil {
		t.Fatalf("GetChatMember: %v", err)
	}

	admin, ok := m.(*ChatMemberAdministrator)
	if !ok {
		t.Fatalf("member = %T", m)
	}

	if admin.Status != StatusAdministrator || admin.User.ID != 99 {
		t.Fatalf("admin = %#v", admin)
	}
}

func TestGetChatMemberCount(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsGetParticipantsRequestTypeID, &tg.ChannelsChannelParticipants{Count: 1234})

	b := newMockBot(inv)

	n, err := b.GetChatMemberCount(context.Background(), tdlibChannel(50))
	if err != nil {
		t.Fatalf("GetChatMemberCount: %v", err)
	}

	if n != 1234 {
		t.Fatalf("count = %d", n)
	}
}

func TestGetChatAdministrators(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsGetParticipantsRequestTypeID, &tg.ChannelsChannelParticipants{
		Count: 2,
		Participants: []tg.ChannelParticipantClass{
			&tg.ChannelParticipantCreator{UserID: 1},
			&tg.ChannelParticipantAdmin{UserID: 2},
		},
		Users: []tg.UserClass{
			&tg.User{ID: 1, AccessHash: 1},
			&tg.User{ID: 2, AccessHash: 2},
		},
	})

	b := newMockBot(inv)

	admins, err := b.GetChatAdministrators(context.Background(), tdlibChannel(50))
	if err != nil {
		t.Fatalf("GetChatAdministrators: %v", err)
	}

	if len(admins) != 2 {
		t.Fatalf("admins = %d", len(admins))
	}

	if _, ok := admins[0].(*ChatMemberOwner); !ok {
		t.Fatalf("first admin = %T, want owner", admins[0])
	}
}

func TestSetChatAdministratorCustomTitle(t *testing.T) {
	inv := newMockInvoker()
	// It reads the current participant rights first, then re-applies them + rank.
	inv.reply(tg.ChannelsGetParticipantRequestTypeID, &tg.ChannelsChannelParticipant{
		Participant: &tg.ChannelParticipantAdmin{
			UserID:      99,
			AdminRights: tg.ChatAdminRights{PinMessages: true},
		},
		Users: []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})
	inv.reply(tg.ChannelsEditAdminRequestTypeID, okUpdates())

	b := newMockBot(inv)

	err := b.SetChatAdministratorCustomTitle(context.Background(), tdlibChannel(50), 99, "Boss")
	if err != nil {
		t.Fatalf("SetChatAdministratorCustomTitle: %v", err)
	}

	var req tg.ChannelsEditAdminRequest

	inv.decode(t, tg.ChannelsEditAdminRequestTypeID, &req)

	if req.Rank != "Boss" {
		t.Fatalf("rank = %q", req.Rank)
	}
}
