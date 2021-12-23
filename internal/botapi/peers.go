package botapi

import (
	"context"
	"strconv"

	"github.com/go-faster/errors"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func fillBotAPIChatPrivate(user peers.User) oas.Chat {
	return oas.Chat{
		ID:        int64(user.TDLibPeerID()),
		Type:      oas.ChatTypePrivate,
		Username:  optString(user.Username),
		FirstName: optString(user.FirstName),
		LastName:  optString(user.LastName),
	}
}

func fillBotAPIChatGroup(chat Chat) oas.Chat {
	r := oas.Chat{
		ID:                  int64(chat.TDLibPeerID()),
		Type:                oas.ChatTypeGroup,
		Title:               oas.NewOptString(chat.VisibleName()),
		Username:            optString(chat.Username),
		HasProtectedContent: oas.NewOptBool(chat.NoForwards()),
	}
	switch ch := chat.(type) {
	case peers.Channel:
		if _, ok := ch.ToBroadcast(); ok {
			r.Type = oas.ChatTypeChannel
		} else {
			r.Type = oas.ChatTypeSupergroup
		}
	}
	return r
}

func (b *BotAPI) getChatByPeer(ctx context.Context, p tg.PeerClass) (oas.Chat, error) {
	peer, err := b.peers.ResolvePeer(ctx, p)
	if err != nil {
		return oas.Chat{}, errors.Wrapf(err, "find peer: %+v", p)
	}

	if user, ok := peer.(peers.User); ok {
		return fillBotAPIChatPrivate(user), nil
	}

	ch, ok := peer.(Chat)
	if !ok {
		return oas.Chat{}, errors.Errorf("unexpected type %T", peer)
	}
	return fillBotAPIChatGroup(ch), nil
}

func (b *BotAPI) resolveUserID(ctx context.Context, id int64) (peers.User, error) {
	user, err := b.peers.GetUser(ctx, &tg.InputUser{UserID: id})
	if err != nil {
		return peers.User{}, errors.Wrapf(err, "find user: %d", id)
	}
	return user, nil
}

// Chat is generic interface for peers.Chat, peers.Channel and friends.
type Chat interface {
	peers.Peer
	Creator() bool
	Left() bool
	NoForwards() bool
	CallActive() bool
	CallNotEmpty() bool
	ParticipantsCount() int
	AdminRights() (tg.ChatAdminRights, bool)
	DefaultBannedRights() (tg.ChatBannedRights, bool)

	Leave(ctx context.Context) error
	SetTitle(ctx context.Context, title string) error
	SetDescription(ctx context.Context, about string) error

	InviteLinks() peers.InviteLinks
	ToSupergroup() (peers.Supergroup, bool)
	ToBroadcast() (peers.Broadcast, bool)
}

var _ = []Chat{
	peers.Chat{},
	peers.Channel{},
}

func (b *BotAPI) resolveIDToChat(ctx context.Context, id oas.ID) (Chat, error) {
	p, err := b.resolveID(ctx, id)
	if err != nil {
		return nil, err
	}
	ch, ok := p.(Chat)
	if !ok {
		return nil, chatNotFound()
	}
	return ch, nil
}

func (b *BotAPI) resolveID(ctx context.Context, id oas.ID) (peers.Peer, error) {
	if id.IsInt64() {
		return b.resolveIntID(ctx, id.Int64)
	}

	username := id.String
	switch {
	case len(username) < 1:
		return nil, &BadRequestError{Message: "Bad Request: chat_id is empty"}
	case username[0] != '@':
		parsedID, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			return nil, chatNotFound()
		}
		return b.resolveIntID(ctx, parsedID)
	}
	// Cut @.
	username = username[1:]

	p, err := b.peers.ResolveDomain(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve %q", username)
	}
	return p, nil
}

func (b *BotAPI) resolveIntID(ctx context.Context, id int64) (peers.Peer, error) {
	p, err := b.peers.ResolveTDLibID(ctx, constant.TDLibPeerID(id))
	if err != nil {
		return nil, errors.Wrapf(err, "find peer %d", id)
	}
	return p, nil
}
