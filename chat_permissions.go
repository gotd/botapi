package botapi

import "github.com/gotd/td/tg"

// ChatPermissions describes the actions a non-administrator user is allowed to
// take in a chat. A nil/false field denies the action.
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendAudios         bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`
	CanSendVideos         bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
}

// ChatAdminRights describes the administrator privileges granted to a user in a
// supergroup or channel. It mirrors the Bot API promoteChatMember parameters.
type ChatAdminRights struct {
	IsAnonymous         bool `json:"is_anonymous,omitempty"`
	CanManageChat       bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool `json:"can_change_info,omitempty"`
	CanInviteUsers      bool `json:"can_invite_users,omitempty"`
	CanPostMessages     bool `json:"can_post_messages,omitempty"`
	CanEditMessages     bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanManageTopics     bool `json:"can_manage_topics,omitempty"`
	// CustomTitle is the administrator's rank (custom title), if any.
	CustomTitle string `json:"custom_title,omitempty"`
}

// toTg converts the Bot API admin rights to the MTProto representation.
func (r ChatAdminRights) toTg() tg.ChatAdminRights {
	return tg.ChatAdminRights{
		ChangeInfo:     r.CanChangeInfo,
		PostMessages:   r.CanPostMessages,
		EditMessages:   r.CanEditMessages,
		DeleteMessages: r.CanDeleteMessages,
		BanUsers:       r.CanRestrictMembers,
		InviteUsers:    r.CanInviteUsers,
		PinMessages:    r.CanPinMessages,
		AddAdmins:      r.CanPromoteMembers,
		Anonymous:      r.IsAnonymous,
		ManageCall:     r.CanManageVideoChats,
		Other:          r.CanManageChat,
		ManageTopics:   r.CanManageTopics,
	}
}

// toBannedRights turns the allow-list permissions into MTProto banned rights,
// where a set bit denies the action. untilDate is a Unix time (0 = forever).
func (p ChatPermissions) toBannedRights(untilDate int) tg.ChatBannedRights {
	canSendAnyMedia := p.CanSendAudios || p.CanSendDocuments || p.CanSendPhotos ||
		p.CanSendVideos || p.CanSendVideoNotes || p.CanSendVoiceNotes
	return tg.ChatBannedRights{
		SendMessages:    !p.CanSendMessages,
		SendMedia:       !canSendAnyMedia,
		SendAudios:      !p.CanSendAudios,
		SendDocs:        !p.CanSendDocuments,
		SendPhotos:      !p.CanSendPhotos,
		SendVideos:      !p.CanSendVideos,
		SendRoundvideos: !p.CanSendVideoNotes,
		SendVoices:      !p.CanSendVoiceNotes,
		SendPolls:       !p.CanSendPolls,
		SendStickers:    !p.CanSendOtherMessages,
		SendGifs:        !p.CanSendOtherMessages,
		SendGames:       !p.CanSendOtherMessages,
		SendInline:      !p.CanSendOtherMessages,
		EmbedLinks:      !p.CanAddWebPagePreviews,
		ChangeInfo:      !p.CanChangeInfo,
		InviteUsers:     !p.CanInviteUsers,
		PinMessages:     !p.CanPinMessages,
		ManageTopics:    !p.CanManageTopics,
		UntilDate:       untilDate,
	}
}
