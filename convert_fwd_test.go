package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestForwardOriginHiddenUser(t *testing.T) {
	b := newTestBot(t)
	h := &tg.MessageFwdHeader{Date: 123}
	h.SetFromName("Anon")

	origin, err := b.forwardOrigin(context.Background(), h)
	if err != nil {
		t.Fatal(err)
	}

	hidden, ok := origin.(*MessageOriginHiddenUser)
	if !ok || hidden.SenderUserName != "Anon" || hidden.Date != 123 {
		t.Fatalf("expected hidden-user origin, got %#v", origin)
	}
}

func TestForwardOriginNoSender(t *testing.T) {
	b := newTestBot(t)
	origin, err := b.forwardOrigin(context.Background(), &tg.MessageFwdHeader{Date: 1})

	if err != nil || origin != nil {
		t.Fatalf("no sender should yield (nil, nil), got (%#v, %v)", origin, err)
	}
}
