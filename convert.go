package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// userFromTgUser converts a raw tg.User into a Bot API User. Pure: optional
// fields are taken as their zero values when absent.
func userFromTgUser(u *tg.User) User {
	return User{
		ID:           u.ID,
		IsBot:        u.Bot,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Username:     u.Username,
		LanguageCode: u.LangCode,
		IsPremium:    u.Premium,
	}
}

// chatFromPeer builds a Bot API Chat from a resolved gotd peer.
//
// The default arm keeps the switch closed against future peer kinds
// (default-signifies-exhaustive).
func chatFromPeer(p peers.Peer) Chat {
	switch p := p.(type) {
	case peers.User:
		c := Chat{ID: int64(p.TDLibPeerID()), Type: ChatTypePrivate}
		if v, ok := p.Username(); ok {
			c.Username = v
		}

		if v, ok := p.FirstName(); ok {
			c.FirstName = v
		}

		if v, ok := p.LastName(); ok {
			c.LastName = v
		}

		return c
	case peers.Channel:
		c := Chat{ID: int64(p.TDLibPeerID()), Title: p.VisibleName(), Type: ChatTypeSupergroup}
		if _, ok := p.ToBroadcast(); ok {
			c.Type = ChatTypeChannel
		}

		if v, ok := p.Username(); ok {
			c.Username = v
		}

		return c
	case peers.Chat:
		return Chat{ID: int64(p.TDLibPeerID()), Type: ChatTypeGroup, Title: p.VisibleName()}
	default:
		return Chat{ID: int64(p.TDLibPeerID())}
	}
}

// inlineKeyboardFromTg converts an MTProto inline markup back into a Bot API
// InlineKeyboardMarkup. Pure.
func inlineKeyboardFromTg(mkp *tg.ReplyInlineMarkup) *InlineKeyboardMarkup {
	rows := make([][]InlineKeyboardButton, len(mkp.Rows))
	for i, row := range mkp.Rows {
		buttons := make([]InlineKeyboardButton, len(row.Buttons))
		for j, b := range row.Buttons {
			btn := InlineKeyboardButton{Text: b.GetText()}
			switch b := b.(type) {
			case *tg.KeyboardButtonURL:
				btn.URL = b.URL
			case *tg.KeyboardButtonCallback:
				btn.CallbackData = string(b.Data)
			case *tg.KeyboardButtonBuy:
				btn.Pay = true
			case *tg.KeyboardButtonSwitchInline:
				q := b.Query
				if b.SamePeer {
					btn.SwitchInlineQueryCurrentChat = &q
				} else {
					btn.SwitchInlineQuery = &q
				}
			case *tg.KeyboardButtonURLAuth:
				// login_url buttons are represented as ordinary url buttons.
				btn.URL = b.URL
			}

			buttons[j] = btn
		}

		rows[i] = buttons
	}

	return &InlineKeyboardMarkup{InlineKeyboard: rows}
}

// chatByPeer resolves a tg.PeerClass to a Bot API Chat via the peer manager.
func (b *Bot) chatByPeer(ctx context.Context, p tg.PeerClass) (Chat, error) {
	peer, err := b.peers.ResolvePeer(ctx, p)
	if err != nil {
		return Chat{}, asAPIError(err)
	}

	return chatFromPeer(peer), nil
}

// fillFrom sets Message.From (for users) or Message.SenderChat (for chats and
// channels) from a tg sender peer.
func (b *Bot) fillFrom(ctx context.Context, from tg.PeerClass, m *Message) error {
	switch from := from.(type) {
	case *tg.PeerUser:
		u, err := b.peers.GetUser(ctx, &tg.InputUser{UserID: from.UserID})
		if err != nil {
			return asAPIError(err)
		}

		user := userFromTgUser(u.Raw())

		m.From = &user
	case *tg.PeerChat, *tg.PeerChannel:
		chat, err := b.chatByPeer(ctx, from)
		if err != nil {
			return err
		}

		m.SenderChat = &chat
	}

	return nil
}

// forwardOrigin maps an MTProto forward header into a Bot API MessageOrigin.
// A hidden sender yields a hidden-user origin (no resolution); otherwise the
// user or chat/channel is resolved via the peer manager.
func (b *Bot) forwardOrigin(ctx context.Context, h *tg.MessageFwdHeader) (MessageOrigin, error) {
	if name, ok := h.GetFromName(); ok {
		return &MessageOriginHiddenUser{Type: OriginHiddenUser, Date: h.Date, SenderUserName: name}, nil
	}

	fromID, ok := h.GetFromID()
	if !ok {
		return nil, nil
	}

	if pu, ok := fromID.(*tg.PeerUser); ok {
		u, err := b.peers.GetUser(ctx, &tg.InputUser{UserID: pu.UserID})
		if err != nil {
			return nil, asAPIError(err)
		}

		return &MessageOriginUser{Type: OriginUser, Date: h.Date, SenderUser: userFromTgUser(u.Raw())}, nil
	}

	chat, err := b.chatByPeer(ctx, fromID)
	if err != nil {
		return nil, err
	}

	author, _ := h.GetPostAuthor()
	if post, ok := h.GetChannelPost(); ok && chat.Type == ChatTypeChannel {
		return &MessageOriginChannel{
			Type:            OriginChannel,
			Date:            h.Date,
			Chat:            chat,
			MessageID:       post,
			AuthorSignature: author,
		}, nil
	}

	return &MessageOriginChat{
		Type:            OriginChat,
		Date:            h.Date,
		SenderChat:      chat,
		AuthorSignature: author,
	}, nil
}

// convertMessage translates a tg.Message into a Bot API Message, resolving the
// chat and sender via the peer manager. It fills the core fields (text,
// entities, reply target, inline markup); richer media mapping arrives with the
// media send/receive work.
func (b *Bot) convertMessage(ctx context.Context, m *tg.Message) (*Message, error) {
	chat, err := b.chatByPeer(ctx, m.PeerID)
	if err != nil {
		return nil, err
	}

	r := &Message{
		MessageID:           m.ID,
		Date:                m.Date,
		Chat:                chat,
		HasProtectedContent: m.Noforwards,
	}
	if v, ok := m.GetEditDate(); ok {
		r.EditDate = v
	}

	if v, ok := m.GetPostAuthor(); ok {
		r.AuthorSignature = v
	}

	if m.Out {
		if self, selfErr := b.peers.Self(ctx); selfErr == nil {
			user := userFromTgUser(self.Raw())

			r.From = &user
		}
	} else if from, ok := m.GetFromID(); ok {
		if err := b.fillFrom(ctx, from, r); err != nil {
			return nil, err
		}
	}

	if rh, ok := m.ReplyTo.(*tg.MessageReplyHeader); ok && rh.ReplyToMsgID != 0 {
		r.ReplyToMessage = &Message{MessageID: rh.ReplyToMsgID}
	}

	if h, ok := m.GetFwdFrom(); ok {
		origin, err := b.forwardOrigin(ctx, &h)
		if err != nil {
			return nil, err
		}

		r.ForwardOrigin = origin
	}

	if media, ok := m.GetMedia(); ok {
		convertMessageMedia(media, r)

		// On a media message, m.Message is the caption.
		if m.Message != "" {
			r.Caption = m.Message
			if len(m.Entities) > 0 {
				r.CaptionEntities = entitiesFromTg(m.Entities)
			}
		}
	} else {
		if m.Message != "" {
			r.Text = m.Message
		}

		if len(m.Entities) > 0 {
			r.Entities = entitiesFromTg(m.Entities)
		}
	}

	if mkp, ok := m.ReplyMarkup.(*tg.ReplyInlineMarkup); ok {
		r.ReplyMarkup = inlineKeyboardFromTg(mkp)
	}

	return r, nil
}
