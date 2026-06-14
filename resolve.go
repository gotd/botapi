package botapi

import (
	"context"
	"strings"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// resolvePeer resolves a Bot API ChatID to a gotd peer, pulling the access hash
// from storage. Numeric ids go through the TDLib id convention; @usernames are
// resolved by domain.
//
// The switch over the sealed ChatID union is exhaustive (gochecksumtype).
func (b *Bot) resolvePeer(ctx context.Context, chat ChatID) (peers.Peer, error) {
	switch c := chat.(type) {
	case ChatIDInt:
		p, err := b.peers.ResolveTDLibID(ctx, constant.TDLibPeerID(int64(c)))
		if err != nil {
			return nil, asAPIError(err)
		}

		return p, nil
	case ChatIDUsername:
		name := strings.TrimPrefix(string(c), "@")
		if name == "" {
			return nil, &Error{Code: 400, Description: "Bad Request: chat_id is empty"}
		}

		p, err := b.peers.ResolveDomain(ctx, name)
		if err != nil {
			return nil, asAPIError(err)
		}

		return p, nil
	case chatIDRef:
		// A PeerRef is addressed directly for sending (resolveInputPeer); it
		// cannot back the peers.Peer that chat-management needs.
		return nil, &Error{Code: 400, Description: "Bad Request: peer reference is only usable for sending"}
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid chat_id"}
	}
}

// resolveInputPeer resolves a ChatID to the tg.InputPeerClass the sender needs.
func (b *Bot) resolveInputPeer(ctx context.Context, chat ChatID) (tg.InputPeerClass, error) {
	// A PeerRef carries its own access hash, so it is addressed directly without
	// consulting stored peer data — this is what makes it survive a restart.
	if ref, ok := chat.(chatIDRef); ok {
		return ref.ref.inputPeer()
	}

	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	return p.InputPeer(), nil
}

// resolveInputUser resolves a Bot API user id to the tg.InputUserClass the
// MTProto user-targeting methods need, pulling the access hash from storage.
func (b *Bot) resolveInputUser(ctx context.Context, userID int64) (tg.InputUserClass, error) {
	u, err := b.peers.ResolveUserID(ctx, userID)
	if err != nil {
		return nil, asAPIError(err)
	}

	return u.InputUser(), nil
}
