package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestCreateForumTopic(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesCreateForumTopicRequestTypeID, &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewChannelMessage{Message: &tg.MessageService{
				ID:     777,
				PeerID: &tg.PeerChannel{ChannelID: 50},
				Action: &tg.MessageActionTopicCreate{
					Title:       "Bugs",
					IconColor:   0x6FB9F0,
					IconEmojiID: 12345,
				},
			}},
		},
	})

	b := newMockBot(inv)

	topic, err := b.CreateForumTopic(context.Background(), tdlibChannel(50), "Bugs",
		WithForumTopicIconColor(0x6FB9F0),
		WithForumTopicIconCustomEmojiID("12345"),
	)
	if err != nil {
		t.Fatalf("CreateForumTopic: %v", err)
	}

	if topic.MessageThreadID != 777 || topic.Name != "Bugs" || topic.IconColor != 0x6FB9F0 || topic.IconCustomEmojiID != "12345" {
		t.Fatalf("topic = %#v", topic)
	}

	var req tg.MessagesCreateForumTopicRequest

	inv.decode(t, tg.MessagesCreateForumTopicRequestTypeID, &req)

	if req.Title != "Bugs" || req.IconColor != 0x6FB9F0 || req.IconEmojiID != 12345 {
		t.Fatalf("req = %#v", req)
	}
}

func TestCreateForumTopicInvalidEmoji(t *testing.T) {
	b := newMockBot(newMockInvoker())

	if _, err := b.CreateForumTopic(context.Background(), tdlibChannel(50), "X",
		WithForumTopicIconCustomEmojiID("notanumber")); err == nil {
		t.Fatal("expected error for invalid icon_custom_emoji_id")
	}
}

func TestEditForumTopic(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditForumTopicRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.EditForumTopic(context.Background(), tdlibChannel(50), 777,
		WithForumTopicName("Renamed"),
		WithForumTopicIconCustomEmojiID(""),
	); err != nil {
		t.Fatalf("EditForumTopic: %v", err)
	}

	var req tg.MessagesEditForumTopicRequest

	inv.decode(t, tg.MessagesEditForumTopicRequestTypeID, &req)

	title, ok := req.GetTitle()
	if !ok || title != "Renamed" {
		t.Fatalf("title = %q ok=%v", title, ok)
	}

	if emoji, ok := req.GetIconEmojiID(); !ok || emoji != 0 {
		t.Fatalf("icon emoji = %d ok=%v, want 0 (removed)", emoji, ok)
	}

	if req.TopicID != 777 {
		t.Fatalf("topic id = %d", req.TopicID)
	}
}

func TestForumTopicCloseReopenHide(t *testing.T) {
	for _, c := range []struct {
		name    string
		call    func(*Bot) error
		topicID int
		closed  *bool
		hidden  *bool
	}{
		{"close", func(b *Bot) error {
			return b.CloseForumTopic(context.Background(), tdlibChannel(50), 5)
		}, 5, ptrBool(true), nil},
		{"reopen", func(b *Bot) error {
			return b.ReopenForumTopic(context.Background(), tdlibChannel(50), 5)
		}, 5, ptrBool(false), nil},
		{"close-general", func(b *Bot) error {
			return b.CloseGeneralForumTopic(context.Background(), tdlibChannel(50))
		}, generalForumTopicID, ptrBool(true), nil},
		{"reopen-general", func(b *Bot) error {
			return b.ReopenGeneralForumTopic(context.Background(), tdlibChannel(50))
		}, generalForumTopicID, ptrBool(false), nil},
		{"hide-general", func(b *Bot) error {
			return b.HideGeneralForumTopic(context.Background(), tdlibChannel(50))
		}, generalForumTopicID, nil, ptrBool(true)},
		{"unhide-general", func(b *Bot) error {
			return b.UnhideGeneralForumTopic(context.Background(), tdlibChannel(50))
		}, generalForumTopicID, nil, ptrBool(false)},
	} {
		t.Run(c.name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesEditForumTopicRequestTypeID, okUpdates())

			if err := c.call(newMockBot(inv)); err != nil {
				t.Fatalf("call: %v", err)
			}

			var req tg.MessagesEditForumTopicRequest

			inv.decode(t, tg.MessagesEditForumTopicRequestTypeID, &req)

			if req.TopicID != c.topicID {
				t.Fatalf("topic id = %d, want %d", req.TopicID, c.topicID)
			}

			if c.closed != nil {
				closed, ok := req.GetClosed()
				if !ok || closed != *c.closed {
					t.Fatalf("closed = %v ok=%v, want %v", closed, ok, *c.closed)
				}
			}

			if c.hidden != nil {
				hidden, ok := req.GetHidden()
				if !ok || hidden != *c.hidden {
					t.Fatalf("hidden = %v ok=%v, want %v", hidden, ok, *c.hidden)
				}
			}
		})
	}
}

func TestEditGeneralForumTopic(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditForumTopicRequestTypeID, okUpdates())

	if err := newMockBot(inv).EditGeneralForumTopic(context.Background(), tdlibChannel(50), "General!"); err != nil {
		t.Fatalf("EditGeneralForumTopic: %v", err)
	}

	var req tg.MessagesEditForumTopicRequest

	inv.decode(t, tg.MessagesEditForumTopicRequestTypeID, &req)

	if req.TopicID != generalForumTopicID {
		t.Fatalf("topic id = %d", req.TopicID)
	}

	if title, ok := req.GetTitle(); !ok || title != "General!" {
		t.Fatalf("title = %q ok=%v", title, ok)
	}
}

func TestDeleteForumTopic(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesDeleteTopicHistoryRequestTypeID, &tg.MessagesAffectedHistory{})

	if err := newMockBot(inv).DeleteForumTopic(context.Background(), tdlibChannel(50), 777); err != nil {
		t.Fatalf("DeleteForumTopic: %v", err)
	}

	var req tg.MessagesDeleteTopicHistoryRequest

	inv.decode(t, tg.MessagesDeleteTopicHistoryRequestTypeID, &req)

	if req.TopMsgID != 777 {
		t.Fatalf("top msg id = %d", req.TopMsgID)
	}
}

func TestUnpinAllForumTopicMessages(t *testing.T) {
	for _, c := range []struct {
		name    string
		call    func(*Bot) error
		topicID int
	}{
		{"topic", func(b *Bot) error {
			return b.UnpinAllForumTopicMessages(context.Background(), tdlibChannel(50), 777)
		}, 777},
		{"general", func(b *Bot) error {
			return b.UnpinAllGeneralForumTopicMessages(context.Background(), tdlibChannel(50))
		}, generalForumTopicID},
	} {
		t.Run(c.name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesUnpinAllMessagesRequestTypeID, &tg.MessagesAffectedHistory{})

			if err := c.call(newMockBot(inv)); err != nil {
				t.Fatalf("call: %v", err)
			}

			var req tg.MessagesUnpinAllMessagesRequest

			inv.decode(t, tg.MessagesUnpinAllMessagesRequestTypeID, &req)

			if topMsg, ok := req.GetTopMsgID(); !ok || topMsg != c.topicID {
				t.Fatalf("top msg id = %d ok=%v, want %d", topMsg, ok, c.topicID)
			}
		})
	}
}

func ptrBool(b bool) *bool { return &b }
