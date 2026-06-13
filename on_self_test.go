package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// TestDispatchSkipsOwnOutgoingMessage guards against the self-reply loop: the
// bot's own messages arrive on the MTProto update stream with Out=true and must
// not be dispatched to handlers.
func TestDispatchSkipsOwnOutgoingMessage(t *testing.T) {
	b := newTestBot(t)
	fired := false
	b.OnMessage(func(*Context) error { fired = true; return nil })

	own := &tg.Message{
		Out:     true,
		ID:      1,
		Message: "echo",
		PeerID:  &tg.PeerUser{UserID: 42},
	}
	b.dispatchMessage(context.Background(), own, false)

	if fired {
		t.Fatal("handler fired for the bot's own outgoing message")
	}
}
