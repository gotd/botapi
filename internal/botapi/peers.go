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

func (b *BotAPI) getChatByPeer(ctx context.Context, p tg.PeerClass) (oas.Chat, error) {
	peer, err := b.peers.ResolvePeer(ctx, p)
	if err != nil {
		return oas.Chat{}, errors.Wrapf(err, "find peer: %+v", p)
	}

	tdlibID := peer.TDLibPeerID()
	if user, ok := peer.(peers.User); ok {
		return oas.Chat{
			ID:        int64(tdlibID),
			Type:      oas.ChatTypePrivate,
			Username:  optString(user.Username),
			FirstName: optString(user.FirstName),
			LastName:  optString(user.LastName),
		}, nil
	}

	r := oas.Chat{
		ID:                  int64(tdlibID),
		Type:                oas.ChatTypeGroup,
		Title:               oas.NewOptString(peer.VisibleName()),
		HasProtectedContent: oas.OptBool{},
	}
	switch ch := peer.(type) {
	case peers.Chat:
		r.HasProtectedContent.SetTo(ch.NoForwards())
	case peers.Channel:
		if _, ok := ch.ToBroadcast(); ok {
			r.Type = oas.ChatTypeChannel
		} else {
			r.Type = oas.ChatTypeSupergroup
		}
		r.Username = optString(ch.Username)
		r.HasProtectedContent.SetTo(ch.NoForwards())
	}

	return r, nil
}

func (b *BotAPI) resolveID(ctx context.Context, id oas.ID) (tg.InputPeerClass, error) {
	if id.IsInt64() {
		return b.resolveIntID(ctx, id.Int64)
	}

	username := id.String
	if len(username) < 1 {
		return nil, &PeerNotFoundError{ID: id}
	}
	switch {
	case len(username) < 1:
		return nil, &PeerNotFoundError{ID: id}
	case username[0] != '@':
		parsedID, err := strconv.ParseInt(username, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse id")
		}
		return b.resolveIntID(ctx, parsedID)
	}
	// Cut @.
	username = username[1:]

	p, err := b.peers.ResolveDomain(ctx, username)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve %q", username)
	}
	return p.InputPeer(), nil
}

func (b *BotAPI) resolveUserID(ctx context.Context, id int64) (*tg.User, error) {
	user, err := b.peers.GetUser(ctx, &tg.InputUser{UserID: id})
	if err != nil {
		return nil, errors.Wrapf(err, "find user: %d", id)
	}
	return user.Raw(), nil
}

func (b *BotAPI) resolveIntID(ctx context.Context, id int64) (tg.InputPeerClass, error) {
	p, err := b.peers.ResolveTDLibID(ctx, constant.TDLibPeerID(id))
	if err != nil {
		return nil, errors.Wrapf(err, "find peer %d", id)
	}
	return p.InputPeer(), nil
}