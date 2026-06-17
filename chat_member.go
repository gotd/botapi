package botapi

// ChatMember is a sealed union describing one member of a chat. The concrete
// type corresponds to the member's status.
//
// Concrete variants: *ChatMemberOwner, *ChatMemberAdministrator,
// *ChatMemberMember, *ChatMemberRestricted, *ChatMemberLeft, *ChatMemberBanned.
type ChatMember interface {
	isChatMember()
}

// ChatMemberOwner is the chat creator.
type ChatMemberOwner struct {
	Status      ChatMemberStatus `json:"status"`
	User        User             `json:"user"`
	IsAnonymous bool             `json:"is_anonymous,omitempty"`
	CustomTitle string           `json:"custom_title,omitempty"`
}

// ChatMemberAdministrator is a chat administrator with granted privileges.
type ChatMemberAdministrator struct {
	Status              ChatMemberStatus `json:"status"`
	User                User             `json:"user"`
	CanBeEdited         bool             `json:"can_be_edited,omitempty"`
	IsAnonymous         bool             `json:"is_anonymous,omitempty"`
	CanManageChat       bool             `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool             `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool             `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool             `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool             `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool             `json:"can_change_info,omitempty"`
	CanInviteUsers      bool             `json:"can_invite_users,omitempty"`
	CanPostMessages     bool             `json:"can_post_messages,omitempty"`
	CanEditMessages     bool             `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool             `json:"can_pin_messages,omitempty"`
	CustomTitle         string           `json:"custom_title,omitempty"`
}

// ChatMemberMember is an ordinary member with no special restrictions.
type ChatMemberMember struct {
	Status    ChatMemberStatus `json:"status"`
	User      User             `json:"user"`
	UntilDate int              `json:"until_date,omitempty"`
	Tag       string           `json:"tag,omitempty"`
}

// ChatMemberRestricted is a member subject to restrictions.
type ChatMemberRestricted struct {
	Status                ChatMemberStatus `json:"status"`
	User                  User             `json:"user"`
	IsMember              bool             `json:"is_member,omitempty"`
	CanSendMessages       bool             `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool             `json:"can_send_media_messages,omitempty"`
	CanSendPolls          bool             `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool             `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool             `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool             `json:"can_change_info,omitempty"`
	CanInviteUsers        bool             `json:"can_invite_users,omitempty"`
	CanPinMessages        bool             `json:"can_pin_messages,omitempty"`
	UntilDate             int              `json:"until_date,omitempty"`
	Tag                   string           `json:"tag,omitempty"`
}

// ChatMemberLeft is a user who is not and was not a member of the chat.
type ChatMemberLeft struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`
}

// ChatMemberBanned is a user banned from the chat.
type ChatMemberBanned struct {
	Status    ChatMemberStatus `json:"status"`
	User      User             `json:"user"`
	UntilDate int              `json:"until_date,omitempty"`
}

func (*ChatMemberOwner) isChatMember()         {}
func (*ChatMemberAdministrator) isChatMember() {}
func (*ChatMemberMember) isChatMember()        {}
func (*ChatMemberRestricted) isChatMember()    {}
func (*ChatMemberLeft) isChatMember()          {}
func (*ChatMemberBanned) isChatMember()        {}

// ChatInviteLink represents an invite link for a chat.
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request,omitempty"`
	IsPrimary               bool   `json:"is_primary,omitempty"`
	IsRevoked               bool   `json:"is_revoked,omitempty"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int    `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
	// SubscriptionPeriod is the number of seconds the subscription is active
	// before the next payment, for subscription invite links.
	SubscriptionPeriod int `json:"subscription_period,omitempty"`
	// SubscriptionPrice is the amount of Telegram Stars a user must pay initially
	// and after each subsequent subscription period to join via the link.
	SubscriptionPrice int `json:"subscription_price,omitempty"`
}

// ChatMemberUpdated represents a change in the status of a chat member.
type ChatMemberUpdated struct {
	Chat          Chat            `json:"chat"`
	From          User            `json:"from"`
	Date          int             `json:"date"`
	OldChatMember ChatMember      `json:"old_chat_member"`
	NewChatMember ChatMember      `json:"new_chat_member"`
	InviteLink    *ChatInviteLink `json:"invite_link,omitempty"`
}
