package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/tg"
)

// Peer kinds for a PeerRef.
const (
	peerKindUser    = "user"
	peerKindChat    = "chat"
	peerKindChannel = "channel"
)

// PeerRef is a serializable reference to a chat that carries its access hash.
//
// Sending to a chat requires its MTProto access hash. The bot harvests and
// persists access hashes for peers it has seen, but addressing a chat after a
// restart otherwise depends on that stored peer data. A PeerRef captures the
// id and access hash in a self-contained, JSON-serializable form: persist it
// yourself (in a DB, a file, …) and, after a restart, send to it with Peer —
// directly from the reference, no re-resolution.
//
//	ref, _ := bot.PeerRef(ctx, botapi.ID(chatID)) // resolve once, capture the hash
//	data, _ := json.Marshal(ref)                  // persist it
//	// … restart …
//	var ref botapi.PeerRef
//	_ = json.Unmarshal(data, &ref)
//	bot.SendMessage(ctx, botapi.Peer(ref), "still works")
type PeerRef struct {
	Kind       string `json:"kind"` // "user", "chat" or "channel"
	ID         int64  `json:"id"`
	AccessHash int64  `json:"access_hash,omitempty"`
}

// PeerRef resolves a chat to a serializable reference including its access hash.
func (b *Bot) PeerRef(ctx context.Context, chat ChatID) (PeerRef, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return PeerRef{}, err
	}

	return peerRefFromInputPeer(peer)
}

// peerRefFromInputPeer extracts a PeerRef from a resolved input peer.
func peerRefFromInputPeer(p tg.InputPeerClass) (PeerRef, error) {
	switch v := p.(type) {
	case *tg.InputPeerUser:
		return PeerRef{Kind: peerKindUser, ID: v.UserID, AccessHash: v.AccessHash}, nil
	case *tg.InputPeerChannel:
		return PeerRef{Kind: peerKindChannel, ID: v.ChannelID, AccessHash: v.AccessHash}, nil
	case *tg.InputPeerChat:
		return PeerRef{Kind: peerKindChat, ID: v.ChatID}, nil
	default:
		return PeerRef{}, &Error{Code: 400, Description: "Bad Request: chat can't be referenced"}
	}
}

// inputPeer reconstructs the MTProto input peer from the reference.
func (r PeerRef) inputPeer() (tg.InputPeerClass, error) {
	switch r.Kind {
	case peerKindUser:
		return &tg.InputPeerUser{UserID: r.ID, AccessHash: r.AccessHash}, nil
	case peerKindChannel:
		return &tg.InputPeerChannel{ChannelID: r.ID, AccessHash: r.AccessHash}, nil
	case peerKindChat:
		return &tg.InputPeerChat{ChatID: r.ID}, nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid peer reference"}
	}
}

// tdlibID returns the TDLib (Bot API) chat id for the reference.
func (r PeerRef) tdlibID() int64 {
	var id constant.TDLibPeerID

	switch r.Kind {
	case peerKindUser:
		id.User(r.ID)
	case peerKindChat:
		id.Chat(r.ID)
	case peerKindChannel:
		id.Channel(r.ID)
	}

	return int64(id)
}

// chatIDRef adapts a PeerRef to the ChatID union.
type chatIDRef struct{ ref PeerRef }

func (chatIDRef) isChatID() {}

// MarshalJSON encodes the reference as its bare Bot API chat id.
func (c chatIDRef) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(c.ref.tdlibID(), 10)), nil
}

// Peer targets a chat by a serializable reference captured with Bot.PeerRef.
// The send is addressed directly from the reference's access hash — no stored
// peer data and no re-resolution — so it works across restarts.
func Peer(ref PeerRef) ChatID { return chatIDRef{ref} }
