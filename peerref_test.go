package botapi

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gotd/td/tg"
)

func TestPeerRefInputPeerRoundTrip(t *testing.T) {
	cases := []struct {
		ref  PeerRef
		want tg.InputPeerClass
	}{
		{PeerRef{Kind: peerKindUser, ID: 7, AccessHash: 11}, &tg.InputPeerUser{UserID: 7, AccessHash: 11}},
		{PeerRef{Kind: peerKindChannel, ID: 9, AccessHash: 13}, &tg.InputPeerChannel{ChannelID: 9, AccessHash: 13}},
		{PeerRef{Kind: peerKindChat, ID: 5}, &tg.InputPeerChat{ChatID: 5}},
	}
	for _, c := range cases {
		got, err := c.ref.inputPeer()
		if err != nil {
			t.Fatal(err)
		}
		if got.String() != c.want.String() {
			t.Fatalf("inputPeer: got %#v, want %#v", got, c.want)
		}
		// Extracting a ref back from the input peer yields the original.
		back, err := peerRefFromInputPeer(got)
		if err != nil || back != c.ref {
			t.Fatalf("round-trip: got %#v (%v), want %#v", back, err, c.ref)
		}
	}
}

func TestPeerSendsDirectlyFromRef(t *testing.T) {
	// resolveInputPeer must build the input peer straight from the ref, with no
	// peer manager or network — this is what survives a restart.
	b := &Bot{}
	ref := PeerRef{Kind: peerKindChannel, ID: 100, AccessHash: 200}
	got, err := b.resolveInputPeer(context.Background(), Peer(ref))
	if err != nil {
		t.Fatal(err)
	}
	ch, ok := got.(*tg.InputPeerChannel)
	if !ok || ch.ChannelID != 100 || ch.AccessHash != 200 {
		t.Fatalf("input peer: %#v", got)
	}
}

func TestPeerRefJSONRoundTrip(t *testing.T) {
	ref := PeerRef{Kind: peerKindUser, ID: 42, AccessHash: 99}
	data, err := json.Marshal(ref)
	if err != nil {
		t.Fatal(err)
	}
	var back PeerRef
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatal(err)
	}
	if back != ref {
		t.Fatalf("json round-trip: got %#v, want %#v", back, ref)
	}
}

func TestPeerChatIDMarshalsToTDLibID(t *testing.T) {
	// Peer(ref) is a ChatID; its JSON form is the bare Bot API chat id.
	data, err := json.Marshal(Peer(PeerRef{Kind: peerKindUser, ID: 42}))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "42" {
		t.Fatalf("user chat id: %s", data)
	}
	// Channels use the -100… supergroup id space.
	data, _ = json.Marshal(Peer(PeerRef{Kind: peerKindChannel, ID: 123}))
	if string(data) != "-1000000000123" {
		t.Fatalf("channel chat id: %s", data)
	}
}
