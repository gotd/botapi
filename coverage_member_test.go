package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

// TestUserToInputPeer covers every branch of the input-user to input-peer
// mapping used by the participant-editing methods.
func TestUserToInputPeer(t *testing.T) {
	cases := []struct {
		in   tg.InputUserClass
		want tg.InputPeerClass
	}{
		{&tg.InputUser{UserID: 1, AccessHash: 2}, &tg.InputPeerUser{UserID: 1, AccessHash: 2}},
		{&tg.InputUserFromMessage{Peer: &tg.InputPeerSelf{}, MsgID: 3, UserID: 4}, &tg.InputPeerUserFromMessage{Peer: &tg.InputPeerSelf{}, MsgID: 3, UserID: 4}},
		{&tg.InputUserSelf{}, &tg.InputPeerSelf{}},
		{&tg.InputUserEmpty{}, &tg.InputPeerEmpty{}},
	}
	for _, c := range cases {
		if got := userToInputPeer(c.in); got.TypeID() != c.want.TypeID() {
			t.Errorf("userToInputPeer(%T) = %T, want %T", c.in, got, c.want)
		}
	}
}

// TestPeerUserIDAndUsersByID covers the non-user fallthrough of peerUserID and
// the non-user skip of usersByID.
func TestPeerUserIDAndUsersByID(t *testing.T) {
	if got := peerUserID(&tg.PeerUser{UserID: 7}); got != 7 {
		t.Fatalf("peerUserID(user) = %d", got)
	}

	if got := peerUserID(&tg.PeerChannel{ChannelID: 9}); got != 0 {
		t.Fatalf("peerUserID(channel) = %d, want 0", got)
	}

	m := usersByID([]tg.UserClass{&tg.User{ID: 1}, &tg.UserEmpty{ID: 2}})
	if len(m) != 1 || m[1] == nil {
		t.Fatalf("usersByID = %#v", m)
	}
}
