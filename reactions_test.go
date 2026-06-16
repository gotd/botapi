package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/tg"
)

func TestSetMessageReaction(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendReactionRequestTypeID, okUpdates())

	b := newMockBot(inv)

	reactions := []ReactionType{Emoji("👍")}
	if err := b.SetMessageReaction(context.Background(), userRef(10, 1), 42, reactions, WithBigReaction()); err != nil {
		t.Fatalf("SetMessageReaction: %v", err)
	}

	var req tg.MessagesSendReactionRequest

	inv.decode(t, tg.MessagesSendReactionRequestTypeID, &req)

	if req.MsgID != 42 || !req.Big {
		t.Fatalf("req = %#v", req)
	}

	if len(req.Reaction) != 1 {
		t.Fatalf("reactions = %#v", req.Reaction)
	}

	emoji, ok := req.Reaction[0].(*tg.ReactionEmoji)
	if !ok || emoji.Emoticon != "👍" {
		t.Fatalf("reaction = %#v", req.Reaction[0])
	}
}

func TestSetMessageReactionEmptyClears(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendReactionRequestTypeID, okUpdates())

	if err := newMockBot(inv).SetMessageReaction(context.Background(), userRef(10, 1), 1, nil); err != nil {
		t.Fatalf("SetMessageReaction: %v", err)
	}

	var req tg.MessagesSendReactionRequest

	inv.decode(t, tg.MessagesSendReactionRequestTypeID, &req)

	if len(req.Reaction) != 0 {
		t.Fatalf("reactions = %#v, want empty", req.Reaction)
	}
}

func TestReactionToTgCustomEmoji(t *testing.T) {
	got, err := reactionToTg(CustomEmoji("12345"))
	if err != nil {
		t.Fatalf("reactionToTg: %v", err)
	}

	custom, ok := got.(*tg.ReactionCustomEmoji)
	if !ok || custom.DocumentID != 12345 {
		t.Fatalf("reaction = %#v", got)
	}

	if _, err := reactionToTg(CustomEmoji("notanumber")); err == nil {
		t.Fatal("expected error for invalid custom_emoji_id")
	}
}

func TestDeleteMessageReaction(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesDeleteParticipantReactionRequestTypeID, okUpdates())

	if err := newMockBot(inv).DeleteMessageReaction(context.Background(), tdlibChannel(50), 42, userRef(60, 1)); err != nil {
		t.Fatalf("DeleteMessageReaction: %v", err)
	}

	var req tg.MessagesDeleteParticipantReactionRequest

	inv.decode(t, tg.MessagesDeleteParticipantReactionRequestTypeID, &req)

	if req.MsgID != 42 {
		t.Fatalf("msg id = %d", req.MsgID)
	}

	if _, ok := req.Participant.(*tg.InputPeerUser); !ok {
		t.Fatalf("participant = %#v, want user", req.Participant)
	}
}

func TestDeleteAllMessageReactions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesDeleteParticipantReactionsRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).DeleteAllMessageReactions(context.Background(), tdlibChannel(50), userRef(60, 1)); err != nil {
		t.Fatalf("DeleteAllMessageReactions: %v", err)
	}

	var req tg.MessagesDeleteParticipantReactionsRequest

	inv.decode(t, tg.MessagesDeleteParticipantReactionsRequestTypeID, &req)

	if _, ok := req.Participant.(*tg.InputPeerUser); !ok {
		t.Fatalf("participant = %#v, want user", req.Participant)
	}
}

func TestChatJoinRequests(t *testing.T) {
	for _, c := range []struct {
		name     string
		call     func(*Bot) error
		approved bool
	}{
		{"approve", func(b *Bot) error {
			return b.ApproveChatJoinRequest(context.Background(), tdlibChannel(50), 99)
		}, true},
		{"decline", func(b *Bot) error {
			return b.DeclineChatJoinRequest(context.Background(), tdlibChannel(50), 99)
		}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.MessagesHideChatJoinRequestRequestTypeID, okUpdates())

			if err := c.call(newMockBot(inv)); err != nil {
				t.Fatalf("call: %v", err)
			}

			var req tg.MessagesHideChatJoinRequestRequest

			inv.decode(t, tg.MessagesHideChatJoinRequestRequestTypeID, &req)

			if req.Approved != c.approved {
				t.Fatalf("approved = %v, want %v", req.Approved, c.approved)
			}
		})
	}
}

func TestBanUnbanChatSenderChat(t *testing.T) {
	var p constant.TDLibPeerID

	p.Channel(60)

	senderID := int64(p)

	for _, c := range []struct {
		name string
		call func(*Bot) error
		view bool
	}{
		{"ban", func(b *Bot) error {
			return b.BanChatSenderChat(context.Background(), tdlibChannel(50), senderID)
		}, true},
		{"unban", func(b *Bot) error {
			return b.UnbanChatSenderChat(context.Background(), tdlibChannel(50), senderID)
		}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.ChannelsEditBannedRequestTypeID, okUpdates())

			if err := c.call(newMockBot(inv)); err != nil {
				t.Fatalf("call: %v", err)
			}

			var req tg.ChannelsEditBannedRequest

			inv.decode(t, tg.ChannelsEditBannedRequestTypeID, &req)

			if _, ok := req.Participant.(*tg.InputPeerChannel); !ok {
				t.Fatalf("participant = %#v, want channel", req.Participant)
			}

			if req.BannedRights.ViewMessages != c.view {
				t.Fatalf("ViewMessages = %v, want %v", req.BannedRights.ViewMessages, c.view)
			}
		})
	}
}
