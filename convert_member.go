package botapi

import (
	"github.com/gotd/td/constant"
	"github.com/gotd/td/tg"
)

// chatMemberFromParticipant converts an MTProto channel participant into a Bot
// API ChatMember. users maps user id to the resolved tg.User harvested from the
// same response so the member's User can be filled in.
func chatMemberFromParticipant(p tg.ChannelParticipantClass, users map[int64]*tg.User) ChatMember {
	user := func(id int64) User {
		if u, ok := users[id]; ok {
			return userFromTgUser(u)
		}

		return User{ID: id}
	}

	switch v := p.(type) {
	case *tg.ChannelParticipantCreator:
		return &ChatMemberOwner{
			Status:      StatusCreator,
			User:        user(v.UserID),
			IsAnonymous: v.AdminRights.Anonymous,
			CustomTitle: v.Rank,
		}
	case *tg.ChannelParticipantAdmin:
		r := v.AdminRights

		return &ChatMemberAdministrator{
			Status:              StatusAdministrator,
			User:                user(v.UserID),
			CanBeEdited:         v.CanEdit,
			IsAnonymous:         r.Anonymous,
			CanManageChat:       r.Other,
			CanDeleteMessages:   r.DeleteMessages,
			CanManageVideoChats: r.ManageCall,
			CanRestrictMembers:  r.BanUsers,
			CanPromoteMembers:   r.AddAdmins,
			CanChangeInfo:       r.ChangeInfo,
			CanInviteUsers:      r.InviteUsers,
			CanPostMessages:     r.PostMessages,
			CanEditMessages:     r.EditMessages,
			CanPinMessages:      r.PinMessages,
			CustomTitle:         v.Rank,
		}
	case *tg.ChannelParticipantSelf:
		return &ChatMemberMember{Status: StatusMember, User: user(v.UserID), Tag: v.Rank}
	case *tg.ChannelParticipant:
		return &ChatMemberMember{Status: StatusMember, User: user(v.UserID), Tag: v.Rank}
	case *tg.ChannelParticipantBanned:
		uid := peerUserID(v.Peer)
		if v.Left {
			return &ChatMemberLeft{Status: StatusLeft, User: user(uid)}
		}

		br := v.BannedRights
		// A member who can still view messages is "restricted"; one who cannot is
		// fully banned ("kicked").
		if br.ViewMessages {
			return &ChatMemberBanned{
				Status:    StatusBanned,
				User:      user(uid),
				UntilDate: br.UntilDate,
			}
		}

		return &ChatMemberRestricted{
			Status:                StatusRestricted,
			User:                  user(uid),
			IsMember:              !v.Left,
			CanSendMessages:       !br.SendMessages,
			CanSendMediaMessages:  !br.SendMedia,
			CanSendPolls:          !br.SendPolls,
			CanSendOtherMessages:  !(br.SendStickers || br.SendGifs || br.SendGames || br.SendInline),
			CanAddWebPagePreviews: !br.EmbedLinks,
			CanChangeInfo:         !br.ChangeInfo,
			CanInviteUsers:        !br.InviteUsers,
			CanPinMessages:        !br.PinMessages,
			Tag:                   v.Rank,
			UntilDate:             br.UntilDate,
		}
	case *tg.ChannelParticipantLeft:
		return &ChatMemberLeft{Status: StatusLeft, User: user(peerUserID(v.Peer))}
	default:
		return &ChatMemberMember{Status: StatusMember}
	}
}

// peerUserID extracts the user id from a peer, or 0 if it is not a user.
func peerUserID(p tg.PeerClass) int64 {
	if u, ok := p.(*tg.PeerUser); ok {
		return u.UserID
	}

	return 0
}

// usersByID indexes resolved tg.User values from an MTProto response by id.
func usersByID(users []tg.UserClass) map[int64]*tg.User {
	m := make(map[int64]*tg.User, len(users))
	for _, u := range users {
		if user, ok := u.(*tg.User); ok {
			m[user.ID] = user
		}
	}

	return m
}

// chatsByID indexes resolved tg.Chat/tg.Channel values from an MTProto response
// by their raw (unmarked) id.
func chatsByID(chats []tg.ChatClass) map[int64]tg.ChatClass {
	m := make(map[int64]tg.ChatClass, len(chats))
	for _, c := range chats {
		switch c := c.(type) {
		case *tg.Chat:
			m[c.ID] = c
		case *tg.ChatForbidden:
			m[c.ID] = c
		case *tg.Channel:
			m[c.ID] = c
		case *tg.ChannelForbidden:
			m[c.ID] = c
		}
	}

	return m
}

// chatFromRaw builds a Bot API Chat from a raw chat/channel peer, enriching it
// from the resolved chats map when the peer is present. It works offline: when
// the chat is absent only the marked id and kind are filled.
func chatFromRaw(p tg.PeerClass, chats map[int64]tg.ChatClass) Chat {
	var id constant.TDLibPeerID

	switch p := p.(type) {
	case *tg.PeerChat:
		id.Chat(p.ChatID)

		c := Chat{ID: int64(id), Type: ChatTypeGroup}

		switch raw := chats[p.ChatID].(type) {
		case *tg.Chat:
			c.Title = raw.Title
		case *tg.ChatForbidden:
			c.Title = raw.Title
		}

		return c
	case *tg.PeerChannel:
		id.Channel(p.ChannelID)

		c := Chat{ID: int64(id), Type: ChatTypeSupergroup}

		switch raw := chats[p.ChannelID].(type) {
		case *tg.Channel:
			c.Title = raw.Title
			c.Username = raw.Username
			c.IsForum = raw.Forum
			c.Type = channelChatType(raw.Broadcast)
		case *tg.ChannelForbidden:
			c.Title = raw.Title
			c.Type = channelChatType(raw.Broadcast)
		}

		return c
	default:
		return Chat{}
	}
}

// channelChatType reports the Bot API chat type for a channel by its broadcast
// flag: broadcast channels are "channel", supergroups are "supergroup".
func channelChatType(broadcast bool) ChatType {
	if broadcast {
		return ChatTypeChannel
	}

	return ChatTypeSupergroup
}
