package botapi

import "github.com/gotd/td/tg"

// KeyboardButtonRequestUsers defines the criteria used to request suitable users
// when pressed. Maps to a request-peer keyboard button.
type KeyboardButtonRequestUsers struct {
	// RequestID is an identifier echoed back in the users_shared service message.
	RequestID int
	// UserIsBot, when set, requires the selected users to be bots (true) or
	// regular users (false).
	UserIsBot *bool
	// UserIsPremium, when set, requires the selected users to have Telegram
	// Premium (true) or not (false).
	UserIsPremium *bool
	// MaxQuantity is the maximum number of users to select (1-10, default 1).
	MaxQuantity int
	// RequestName, RequestUsername and RequestPhoto request that the selected
	// users' details be included client-side. They are not enforced over MTProto.
	RequestName     bool
	RequestUsername bool
	RequestPhoto    bool
}

// KeyboardButtonRequestChat defines the criteria used to request a suitable chat
// when pressed. Maps to a request-peer keyboard button.
type KeyboardButtonRequestChat struct {
	// RequestID is an identifier echoed back in the chat_shared service message.
	RequestID int
	// ChatIsChannel requires a channel (true) or a group/supergroup (false).
	ChatIsChannel bool
	// ChatIsForum, when set, requires a forum (true) or non-forum (false) chat.
	ChatIsForum *bool
	// ChatHasUsername, when set, requires a chat with (true) or without (false) a
	// public username.
	ChatHasUsername *bool
	// ChatIsCreated requires a chat owned by the user.
	ChatIsCreated bool
	// UserAdministratorRights requires the user to hold these rights in the chat.
	UserAdministratorRights *ChatAdminRights
	// BotAdministratorRights requires the bot to hold these rights in the chat.
	BotAdministratorRights *ChatAdminRights
	// BotIsMember requires the bot to be a member of the chat.
	BotIsMember bool
	// RequestTitle, RequestUsername and RequestPhoto request that the chat's
	// details be included client-side. They are not enforced over MTProto.
	RequestTitle    bool
	RequestUsername bool
	RequestPhoto    bool
}

// requestUsersToTg builds the MTProto request-peer button for a user request.
func requestUsersToTg(text string, r *KeyboardButtonRequestUsers) *tg.KeyboardButtonRequestPeer {
	peerType := &tg.RequestPeerTypeUser{}
	if r.UserIsBot != nil {
		peerType.SetBot(*r.UserIsBot)
	}

	if r.UserIsPremium != nil {
		peerType.SetPremium(*r.UserIsPremium)
	}

	maxQuantity := r.MaxQuantity
	if maxQuantity == 0 {
		maxQuantity = 1
	}

	return &tg.KeyboardButtonRequestPeer{
		Text:        text,
		ButtonID:    r.RequestID,
		PeerType:    peerType,
		MaxQuantity: maxQuantity,
	}
}

// requestChatToTg builds the MTProto request-peer button for a chat request.
func requestChatToTg(text string, r *KeyboardButtonRequestChat) *tg.KeyboardButtonRequestPeer {
	var peerType tg.RequestPeerTypeClass

	if r.ChatIsChannel {
		broadcast := &tg.RequestPeerTypeBroadcast{Creator: r.ChatIsCreated}
		if r.UserAdministratorRights != nil {
			broadcast.SetUserAdminRights(r.UserAdministratorRights.toTg())
		}

		if r.BotAdministratorRights != nil {
			broadcast.SetBotAdminRights(r.BotAdministratorRights.toTg())
		}

		if r.ChatHasUsername != nil {
			broadcast.SetHasUsername(*r.ChatHasUsername)
		}

		peerType = broadcast
	} else {
		chat := &tg.RequestPeerTypeChat{Creator: r.ChatIsCreated, BotParticipant: r.BotIsMember}
		if r.UserAdministratorRights != nil {
			chat.SetUserAdminRights(r.UserAdministratorRights.toTg())
		}

		if r.BotAdministratorRights != nil {
			chat.SetBotAdminRights(r.BotAdministratorRights.toTg())
		}

		if r.ChatHasUsername != nil {
			chat.SetHasUsername(*r.ChatHasUsername)
		}

		if r.ChatIsForum != nil {
			chat.SetForum(*r.ChatIsForum)
		}

		peerType = chat
	}

	return &tg.KeyboardButtonRequestPeer{
		Text:     text,
		ButtonID: r.RequestID,
		PeerType: peerType,
	}
}
