package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestOwnedGiftToTg(t *testing.T) {
	t.Run("message", func(t *testing.T) {
		got, err := ownedGiftToTg(OwnedGiftFromMessage(42))
		if err != nil {
			t.Fatalf("decode: %v", err)
		}

		g, ok := got.(*tg.InputSavedStarGiftUser)
		if !ok || g.MsgID != 42 {
			t.Fatalf("got %#v", got)
		}
	})

	t.Run("slug", func(t *testing.T) {
		got, err := ownedGiftToTg(OwnedGiftFromSlug("abc"))
		if err != nil {
			t.Fatalf("decode: %v", err)
		}

		g, ok := got.(*tg.InputSavedStarGiftSlug)
		if !ok || g.Slug != "abc" {
			t.Fatalf("got %#v", got)
		}
	})

	t.Run("chat", func(t *testing.T) {
		got, err := ownedGiftToTg("chat:777:888:5")
		if err != nil {
			t.Fatalf("decode: %v", err)
		}

		g, ok := got.(*tg.InputSavedStarGiftChat)
		if !ok || g.SavedID != 5 {
			t.Fatalf("got %#v", got)
		}

		peer, ok := g.Peer.(*tg.InputPeerChannel)
		if !ok || peer.ChannelID != 777 || peer.AccessHash != 888 {
			t.Fatalf("peer %#v", g.Peer)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for _, id := range []string{"", "nope", "msg:", "msg:x", "chat:1:2", "slug:"} {
			if _, err := ownedGiftToTg(id); err == nil {
				t.Fatalf("expected error for %q", id)
			}
		}
	})
}
