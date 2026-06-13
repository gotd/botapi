package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestChatMemberFromParticipant(t *testing.T) {
	users := usersByID([]tg.UserClass{
		&tg.User{ID: 1, FirstName: "Owner"},
		&tg.User{ID: 2, FirstName: "Admin"},
		&tg.User{ID: 3, FirstName: "Plain"},
		&tg.User{ID: 4, FirstName: "Restricted"},
	})

	creator := chatMemberFromParticipant(&tg.ChannelParticipantCreator{
		UserID:      1,
		Rank:        "boss",
		AdminRights: tg.ChatAdminRights{Anonymous: true},
	}, users)
	owner, ok := creator.(*ChatMemberOwner)
	if !ok || owner.Status != StatusCreator || !owner.IsAnonymous || owner.CustomTitle != "boss" {
		t.Fatalf("creator: %#v", creator)
	}
	if owner.User.FirstName != "Owner" {
		t.Fatalf("creator user not resolved: %#v", owner.User)
	}

	admin := chatMemberFromParticipant(&tg.ChannelParticipantAdmin{
		UserID:      2,
		CanEdit:     true,
		AdminRights: tg.ChatAdminRights{BanUsers: true, DeleteMessages: true},
	}, users)
	a, ok := admin.(*ChatMemberAdministrator)
	if !ok || !a.CanRestrictMembers || !a.CanDeleteMessages || !a.CanBeEdited {
		t.Fatalf("admin: %#v", admin)
	}

	plain := chatMemberFromParticipant(&tg.ChannelParticipant{UserID: 3}, users)
	if m, isMember := plain.(*ChatMemberMember); !isMember || m.User.FirstName != "Plain" {
		t.Fatalf("member: %#v", plain)
	}

	restricted := chatMemberFromParticipant(&tg.ChannelParticipantBanned{
		Peer:         &tg.PeerUser{UserID: 4},
		BannedRights: tg.ChatBannedRights{SendMedia: true, UntilDate: 123},
	}, users)
	r, ok := restricted.(*ChatMemberRestricted)
	if !ok {
		t.Fatalf("restricted: %#v", restricted)
	}
	if r.CanSendMediaMessages {
		t.Fatal("restricted: media should be denied")
	}
	if !r.CanSendMessages {
		t.Fatal("restricted: messages should be allowed")
	}
	if r.UntilDate != 123 {
		t.Fatalf("restricted: until date %d", r.UntilDate)
	}

	banned := chatMemberFromParticipant(&tg.ChannelParticipantBanned{
		Peer:         &tg.PeerUser{UserID: 4},
		BannedRights: tg.ChatBannedRights{ViewMessages: true},
	}, users)
	if _, ok := banned.(*ChatMemberBanned); !ok {
		t.Fatalf("banned: %#v", banned)
	}
}

func TestChatPermissionsToBannedRights(t *testing.T) {
	// Allow text only; everything else denied.
	perms := ChatPermissions{CanSendMessages: true}
	br := perms.toBannedRights(0)
	if br.SendMessages {
		t.Fatal("SendMessages should be allowed (not banned)")
	}
	if !br.SendPolls || !br.SendMedia || !br.EmbedLinks {
		t.Fatal("non-text actions should be banned")
	}

	all := ChatPermissions{
		CanSendMessages: true, CanSendPhotos: true, CanSendVideos: true,
		CanSendAudios: true, CanSendDocuments: true, CanSendVideoNotes: true,
		CanSendVoiceNotes: true, CanSendPolls: true, CanSendOtherMessages: true,
		CanAddWebPagePreviews: true,
	}
	br = all.toBannedRights(0)
	if br.SendMessages || br.SendMedia || br.SendPolls || br.EmbedLinks || br.SendStickers {
		t.Fatalf("fully-permitted user should have no media bans: %#v", br)
	}
}

func TestChatAdminRightsToTg(t *testing.T) {
	r := ChatAdminRights{CanManageChat: true, CanPromoteMembers: true, IsAnonymous: true}
	tgr := r.toTg()
	if !tgr.Other || !tgr.AddAdmins || !tgr.Anonymous {
		t.Fatalf("admin rights mapping: %#v", tgr)
	}
	if tgr.PostMessages {
		t.Fatal("PostMessages should be false")
	}
}
