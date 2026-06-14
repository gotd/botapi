package botapi

import (
	"context"
	"sync"
	"testing"

	"github.com/go-faster/errors"

	"github.com/gotd/log"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
)

// mockInvoker is an in-memory tg.Invoker for hermetic method tests. It dispatches
// by request TypeID to registered handlers, capturing each request so tests can
// assert on the exact MTProto call a Bot method makes, and returns a canned
// response that the real client decodes back into the method's result type.
//
// It ships with default handlers for the peer-resolution RPCs (users.getUsers,
// channels.getChannels) so methods that resolve a chat or convert an incoming
// message work without a network or seeded storage; tests override any RPC with
// reply/handle.
type mockInvoker struct {
	mu       sync.Mutex
	handlers map[uint32]func(buf *bin.Buffer) (bin.Encoder, error)
	calls    []uint32
	last     map[uint32][]byte
}

func newMockInvoker() *mockInvoker {
	m := &mockInvoker{
		handlers: map[uint32]func(*bin.Buffer) (bin.Encoder, error){},
		last:     map[uint32][]byte{},
	}
	// Default resolution: echo requested users/channels back as full entities.
	m.handle(tg.UsersGetUsersRequestTypeID, func(buf *bin.Buffer) (bin.Encoder, error) {
		req := &tg.UsersGetUsersRequest{}
		if err := req.Decode(buf); err != nil {
			return nil, err
		}
		var users []tg.UserClass
		for _, iu := range req.ID {
			switch u := iu.(type) {
			case *tg.InputUserSelf:
				users = append(users, &tg.User{ID: 1, Self: true, Bot: true, AccessHash: 1, Username: "test_bot"})
			case *tg.InputUser:
				users = append(users, &tg.User{ID: u.UserID, AccessHash: u.AccessHash})
			case *tg.InputUserFromMessage:
				users = append(users, &tg.User{ID: u.UserID})
			}
		}
		return &tg.UserClassVector{Elems: users}, nil
	})
	m.handle(tg.ChannelsGetChannelsRequestTypeID, func(buf *bin.Buffer) (bin.Encoder, error) {
		req := &tg.ChannelsGetChannelsRequest{}
		if err := req.Decode(buf); err != nil {
			return nil, err
		}
		var chats []tg.ChatClass
		for _, ic := range req.ID {
			if c, ok := ic.(*tg.InputChannel); ok {
				chats = append(chats, &tg.Channel{
					ID:         c.ChannelID,
					AccessHash: c.AccessHash,
					Title:      "channel",
					Photo:      &tg.ChatPhotoEmpty{},
				})
			}
		}
		return &tg.MessagesChats{Chats: chats}, nil
	})
	return m
}

// handle registers a handler for an RPC request TypeID. The handler receives a
// buffer positioned at the start of the request and returns the response to
// decode back, or an error (use a *tgerr.Error to exercise error mapping).
func (m *mockInvoker) handle(id uint32, h func(buf *bin.Buffer) (bin.Encoder, error)) {
	m.mu.Lock()
	m.handlers[id] = h
	m.mu.Unlock()
}

// reply registers a constant response for an RPC request TypeID.
func (m *mockInvoker) reply(id uint32, resp bin.Encoder) {
	m.handle(id, func(*bin.Buffer) (bin.Encoder, error) { return resp, nil })
}

// fail registers an RPC that returns the given Telegram error.
func (m *mockInvoker) fail(id uint32, err *tgerr.Error) {
	m.handle(id, func(*bin.Buffer) (bin.Encoder, error) { return nil, err })
}

// called reports whether an RPC with the given TypeID was invoked.
func (m *mockInvoker) called(id uint32) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.last[id]
	return ok
}

// count returns how many RPCs were invoked in total.
func (m *mockInvoker) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.calls)
}

// decode decodes the most recent request with the given TypeID into dst.
func (m *mockInvoker) decode(t *testing.T, id uint32, dst bin.Decoder) {
	t.Helper()
	m.mu.Lock()
	raw, ok := m.last[id]
	m.mu.Unlock()
	if !ok {
		t.Fatalf("no recorded request for %#x", id)
	}
	if err := dst.Decode(&bin.Buffer{Buf: raw}); err != nil {
		t.Fatalf("decode request %#x: %v", id, err)
	}
}

func (m *mockInvoker) Invoke(_ context.Context, input bin.Encoder, output bin.Decoder) error {
	var buf bin.Buffer
	if err := input.Encode(&buf); err != nil {
		return err
	}
	id, err := buf.PeekID()
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.calls = append(m.calls, id)
	m.last[id] = append([]byte(nil), buf.Buf...)
	h := m.handlers[id]
	m.mu.Unlock()

	if h == nil {
		return errors.Errorf("mockInvoker: unhandled request %#x", id)
	}
	resp, err := h(&buf)
	if err != nil {
		return err
	}
	var rb bin.Buffer
	if err := resp.Encode(&rb); err != nil {
		return err
	}
	return output.Decode(&rb)
}

// newMockBot builds a Bot wired to a mock invoker, with in-memory peers and a
// canned self user, so every RPC method can be driven offline.
func newMockBot(inv *mockInvoker) *Bot {
	raw := tg.NewClient(inv)
	return &Bot{
		log:    log.Nop,
		raw:    raw,
		sender: message.NewSender(raw),
		peers:  peers.Options{}.Build(raw),
		disp:   tg.NewUpdateDispatcher(),
		self:   &tg.User{ID: 1, Bot: true, AccessHash: 1, Username: "test_bot"},
	}
}

// userRef is a PeerRef addressing a user with a known access hash, so sends to
// it skip peer resolution.
func userRef(id, hash int64) ChatID {
	return Peer(PeerRef{Kind: peerKindUser, ID: id, AccessHash: hash})
}

// channelRef is a PeerRef addressing a channel with a known access hash.
func channelRef(id, hash int64) ChatID {
	return Peer(PeerRef{Kind: peerKindChannel, ID: id, AccessHash: hash})
}

// tdlibChannel/tdlibUser/tdlibChat build a numeric ChatID (TDLib id convention)
// that resolves through the peers manager — the mock invoker answers the
// resulting users.getUsers / channels.getChannels. Use these for the
// chat-management methods, which reject a PeerRef.
func tdlibChannel(id int64) ChatID {
	var p constant.TDLibPeerID
	p.Channel(id)
	return ID(int64(p))
}

func tdlibUser(id int64) ChatID {
	var p constant.TDLibPeerID
	p.User(id)
	return ID(int64(p))
}

func tdlibChat(id int64) ChatID {
	var p constant.TDLibPeerID
	p.Chat(id)
	return ID(int64(p))
}

// okUpdates is an empty successful Updates response, used by methods that ignore
// their result.
func okUpdates() *tg.Updates { return &tg.Updates{} }

// messageUpdates wraps a single new message (plus its peer user) in an Updates,
// the shape send/edit methods unpack into a Message.
func messageUpdates(msg *tg.Message) *tg.Updates {
	return &tg.Updates{
		Updates: []tg.UpdateClass{&tg.UpdateNewMessage{Message: msg}},
		Users:   []tg.UserClass{&tg.User{ID: 1, AccessHash: 1, Bot: true, Username: "test_bot"}},
	}
}
