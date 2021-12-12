package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

const (
	// MaxTDLibChatID is maximum chat TDLib ID.
	MaxTDLibChatID = 999999999999
	// MaxTDLibChannelID is maximum channel TDLib ID.
	MaxTDLibChannelID = 1000000000000 - int64(1<<31)
	// ZeroTDLibChannelID is minimum channel TDLib ID.
	ZeroTDLibChannelID = -1000000000000
	// MaxTDLibUserID is maximum user TDLib ID.
	MaxTDLibUserID = (1 << 40) - 1
)

func fullChatInputPeer(full tg.FullChat) tg.InputPeerClass {
	switch full := full.(type) {
	case *tg.Chat:
		return &tg.InputPeerChat{ChatID: full.ID}
	case *tg.Channel:
		return &tg.InputPeerChannel{
			ChannelID:  full.ID,
			AccessHash: full.AccessHash,
		}
	default:
		return &tg.InputPeerEmpty{}
	}
}

func toTDLibID(p tg.InputPeerClass) int64 {
	switch p := p.(type) {
	case *tg.InputPeerUser:
		return p.GetUserID()
	case *tg.InputPeerChat:
		return -p.GetChatID()
	case *tg.InputPeerChannel:
		return ZeroTDLibChannelID - (p.GetChannelID() * -1)
	default:
		return 0
	}
}

func fromTDLibID(id int64) int64 {
	switch {
	case IsUserTDLibID(id):
	case IsChatTDLibID(id):
		id = -id
	case IsChannelTDLibID(id):
		id -= ZeroTDLibChannelID
		id = -id
	}
	return id
}

// IsUserTDLibID whether that given ID is user ID.
func IsUserTDLibID(id int64) bool {
	return id > 0 && id <= MaxTDLibUserID
}

// IsChatTDLibID whether that given ID is chat ID.
func IsChatTDLibID(id int64) bool {
	return id < 0 && -MaxTDLibChatID <= id
}

// IsChannelTDLibID whether that given ID is channel ID.
func IsChannelTDLibID(id int64) bool {
	return id < 0 &&
		id != ZeroTDLibChannelID &&
		!IsChatTDLibID(id) &&
		ZeroTDLibChannelID-MaxTDLibChannelID <= id
}

func (b *BotAPI) getChatByPeer(ctx context.Context, p tg.PeerClass) (oas.Chat, error) {
	var chatID int64
	switch p := p.(type) {
	case *tg.PeerUser:
		user, ok, err := b.peers.FindUser(ctx, p.UserID)
		switch {
		case err != nil:
			return oas.Chat{}, errors.Wrapf(err, "find user: %d", p.UserID)
		case !ok:
			return oas.Chat{}, errors.Errorf("can't find user %d", p.UserID)
		}
		return oas.Chat{
			ID:        toTDLibID(user.AsInputPeer()),
			Type:      oas.ChatTypePrivate,
			Username:  optString(user.GetUsername),
			FirstName: optString(user.GetFirstName),
			LastName:  optString(user.GetLastName),
		}, nil
	case *tg.PeerChat:
		chatID = p.ChatID
	case *tg.PeerChannel:
		chatID = p.ChannelID
	default:
		return oas.Chat{}, errors.Errorf("unexpected type %T", p)
	}

	chat, ok, err := b.peers.FindChat(ctx, chatID)
	switch {
	case err != nil:
		return oas.Chat{}, errors.Wrapf(err, "find chat: %d", chatID)
	case !ok:
		return oas.Chat{}, errors.Errorf("can't find chat %d", chatID)
	}

	r := oas.Chat{
		ID:    toTDLibID(fullChatInputPeer(chat)),
		Type:  oas.ChatTypeGroup,
		Title: oas.NewOptString(chat.GetTitle()),
		// TODO(tdakkota): set more fields, when gotd schema will be updated
		HasProtectedContent: oas.OptBool{},
	}
	switch ch := chat.(type) {
	case *tg.Chat:
	case *tg.Channel:
		if ch.Broadcast {
			r.Type = oas.ChatTypeChannel
		} else {
			r.Type = oas.ChatTypeSupergroup
		}
		r.Username = optString(ch.GetUsername)
	}

	return r, nil
}

func (b *BotAPI) resolveID(ctx context.Context, id oas.ID) (tg.InputPeerClass, error) {
	if id.IsInt64() {
		return b.resolveIntID(ctx, id)
	}

	username := id.String
	if len(username) < 1 || username[0] != '@' {
		return nil, &PeerNotFoundError{ID: id}
	}
	// Cut @.
	username = username[1:]

	p, err := b.resolver.ResolveDomain(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve %q", username)
	}
	switch p.(type) {
	case *tg.InputPeerChat, *tg.InputPeerChannel:
		return p, nil
	default:
		return nil, &PeerNotFoundError{ID: id}
	}
}

func (b *BotAPI) resolveUserID(ctx context.Context, id int64) (*tg.User, error) {
	user, ok, err := b.peers.FindUser(ctx, id)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "find user: %d", id)
	case !ok:
		return nil, &PeerNotFoundError{ID: oas.NewInt64ID(id)}
	}
	return user, nil
}

func (b *BotAPI) resolveIntID(ctx context.Context, chatID oas.ID) (tg.InputPeerClass, error) {
	id := chatID.Int64
	cleanID := fromTDLibID(id)

	if IsUserTDLibID(id) {
		user, err := b.resolveUserID(ctx, cleanID)
		if err != nil {
			return nil, err
		}
		return user.AsInputPeer(), nil
	}

	chat, ok, err := b.peers.FindChat(ctx, cleanID)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "find chat: %d", id)
	case !ok:
		return nil, &PeerNotFoundError{ID: chatID}
	}
	switch chat := chat.(type) {
	case *tg.Chat:
		return chat.AsInputPeer(), nil
	case *tg.Channel:
		return chat.AsInputPeer(), nil
	default:
		return nil, &PeerNotFoundError{ID: chatID}
	}
}
