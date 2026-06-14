package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestEditedMessageFromResp(t *testing.T) {
	want := &tg.Message{ID: 42, Message: "edited"}

	cases := []struct {
		name string
		resp tg.UpdatesClass
		ok   bool
	}{
		{
			name: "EditMessage",
			resp: &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateEditMessage{Message: want}}},
			ok:   true,
		},
		{
			name: "EditChannelMessage",
			resp: &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateEditChannelMessage{Message: want}}},
			ok:   true,
		},
		{
			name: "UpdateShort",
			resp: &tg.UpdateShort{Update: &tg.UpdateEditMessage{Message: want}},
			ok:   true,
		},
		{
			name: "NoEdit",
			resp: &tg.Updates{Updates: []tg.UpdateClass{&tg.UpdateNewMessage{Message: want}}},
			ok:   false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := editedMessageFromResp(c.resp)
			if ok != c.ok {
				t.Fatalf("ok: got %v, want %v", ok, c.ok)
			}

			if ok && got.GetID() != 42 {
				t.Fatalf("message id: %d", got.GetID())
			}
		})
	}
}
