package botapi

import (
	"context"

	"github.com/gotd/log"
	"github.com/gotd/td/tg"
)

// installBusinessHandlers wires the business-connection updates to the router. It
// is called from installHandlers.
func (b *Bot) installBusinessHandlers() {
	b.disp.OnBotBusinessConnect(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotBusinessConnect) error {
		b.route(ctx, &Update{BusinessConnection: b.businessConnectionFromTg(u.Connection, e.Users)})

		return nil
	})
	b.disp.OnBotNewBusinessMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotNewBusinessMessage) error {
		b.dispatchBusinessMessage(ctx, e, u.ConnectionID, u.Message, false)

		return nil
	})
	b.disp.OnBotEditBusinessMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateBotEditBusinessMessage) error {
		b.dispatchBusinessMessage(ctx, e, u.ConnectionID, u.Message, true)

		return nil
	})
	b.disp.OnBotDeleteBusinessMessage(func(ctx context.Context, _ tg.Entities, u *tg.UpdateBotDeleteBusinessMessage) error {
		chat, err := b.chatByPeer(ctx, u.Peer)
		if err != nil {
			b.logger().Error(ctx, "Convert business deleted messages", log.Error(err))

			return nil
		}

		b.route(ctx, &Update{DeletedBusinessMessages: &BusinessMessagesDeleted{
			BusinessConnectionID: u.ConnectionID,
			Chat:                 chat,
			MessageIDs:           u.Messages,
		}})

		return nil
	})
}

// dispatchBusinessMessage converts a message delivered over a business connection
// and routes it as a (edited) business message. Unlike dispatchMessage it keeps
// outgoing messages: the Bot API delivers the whole business conversation,
// including messages the account itself sent.
func (b *Bot) dispatchBusinessMessage(ctx context.Context, e tg.Entities, connectionID string, msg tg.MessageClass, edited bool) {
	m, err := b.messageFromTg(ctx, msg)
	if err != nil {
		b.logger().Error(ctx, "Convert business message", log.Error(err))

		return
	}

	if m == nil {
		return
	}

	m.BusinessConnectionID = connectionID
	if m.raw != nil {
		// Drop redelivered messages so the bot does not reply to the same one
		// twice (Telegram may resend updates, especially to bots on reconnect).
		if b.businessSeen != nil && !b.businessSeen.fresh(businessMessageKey(connectionID, m.raw)) {
			return
		}

		m.businessPeer = inputPeerFromEntities(m.raw.PeerID, e)
	}

	u := &Update{}
	if edited {
		u.EditedBusinessMessage = m
	} else {
		u.BusinessMessage = m
	}

	b.route(ctx, u)
}

// inputPeerFromEntities builds the input peer for a message's chat from the
// access hashes delivered with the update, so a business send addresses the chat
// with the account's own access hash. Returns nil when the peer is absent.
func inputPeerFromEntities(peer tg.PeerClass, e tg.Entities) tg.InputPeerClass {
	switch p := peer.(type) {
	case *tg.PeerUser:
		if u, ok := e.Users[p.UserID]; ok {
			return &tg.InputPeerUser{UserID: p.UserID, AccessHash: u.AccessHash}
		}
	case *tg.PeerChannel:
		if ch, ok := e.Channels[p.ChannelID]; ok {
			return &tg.InputPeerChannel{ChannelID: p.ChannelID, AccessHash: ch.AccessHash}
		}
	case *tg.PeerChat:
		return &tg.InputPeerChat{ChatID: p.ChatID}
	}

	return nil
}

// businessChatID returns a send target for a business message's chat that carries
// the account-scoped access hash, bypassing the bot's peer store.
func businessChatID(m *Message) (ChatID, bool) {
	if m == nil || m.businessPeer == nil {
		return nil, false
	}

	ref, err := peerRefFromInputPeer(m.businessPeer)
	if err != nil {
		return nil, false
	}

	return Peer(ref), true
}

// businessConnectionID returns the connection id the update belongs to, or empty
// when the update is not a business update.
func (u *Update) businessConnectionID() string {
	switch {
	case u.BusinessMessage != nil:
		return u.BusinessMessage.BusinessConnectionID
	case u.EditedBusinessMessage != nil:
		return u.EditedBusinessMessage.BusinessConnectionID
	case u.BusinessConnection != nil:
		return u.BusinessConnection.ID
	case u.DeletedBusinessMessages != nil:
		return u.DeletedBusinessMessages.BusinessConnectionID
	default:
		return ""
	}
}

// Kind predicates for business updates.
func hasBusinessMessage(c *Context) bool         { return c.Update.BusinessMessage != nil }
func hasEditedBusinessMessage(c *Context) bool   { return c.Update.EditedBusinessMessage != nil }
func hasBusinessConnection(c *Context) bool      { return c.Update.BusinessConnection != nil }
func hasDeletedBusinessMessages(c *Context) bool { return c.Update.DeletedBusinessMessages != nil }

// OnBusinessMessage registers a handler for new messages from a connected
// business account.
func (b *Bot) OnBusinessMessage(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasBusinessMessage, predicates)...)
}

// OnEditedBusinessMessage registers a handler for edited messages from a
// connected business account.
func (b *Bot) OnEditedBusinessMessage(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasEditedBusinessMessage, predicates)...)
}

// OnBusinessConnection registers a handler for business connection updates (the
// bot was connected to, disconnected from, or had its rights changed on a
// business account).
func (b *Bot) OnBusinessConnection(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasBusinessConnection, predicates)...)
}

// OnDeletedBusinessMessages registers a handler for messages deleted from a
// connected business account.
func (b *Bot) OnDeletedBusinessMessages(h Handler, predicates ...Predicate) {
	b.on(h, prepend(hasDeletedBusinessMessages, predicates)...)
}

// BusinessMessage returns the business message the update carries (new or
// edited), or nil when the update carries none.
func (c *Context) BusinessMessage() *Message {
	switch {
	case c.Update.BusinessMessage != nil:
		return c.Update.BusinessMessage
	case c.Update.EditedBusinessMessage != nil:
		return c.Update.EditedBusinessMessage
	default:
		return nil
	}
}

// BusinessChat returns a send target for the business message's chat that
// carries the account-scoped access hash, and whether the update is a business
// message. Use it to address sends on behalf of the connected account; the bot's
// stored access hash is invalid in the business context.
func (c *Context) BusinessChat() (ChatID, bool) {
	return businessChatID(c.BusinessMessage())
}

// Business returns a BusinessContext scoped to the update's business connection,
// and whether the update belongs to one. Use it to act on behalf of the
// connected account from within a business update handler.
func (c *Context) Business() (*BusinessContext, bool) {
	id := c.Update.businessConnectionID()
	if id == "" {
		return nil, false
	}

	return c.Bot.Business(id), true
}
