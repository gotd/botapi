package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

// broadcastInvoker returns a mock invoker whose channel resolution yields a
// broadcast channel, so messages addressed to it convert to channel posts.
func broadcastInvoker() *mockInvoker {
	inv := newMockInvoker()
	inv.handle(tg.ChannelsGetChannelsRequestTypeID, func(*bin.Buffer) (bin.Encoder, error) {
		return &tg.MessagesChats{Chats: []tg.ChatClass{
			&tg.Channel{ID: 50, AccessHash: 1, Title: "ch", Broadcast: true, Photo: &tg.ChatPhotoEmpty{}},
		}}, nil
	})
	return inv
}

// TestDispatchMessageRouting drives every routing branch of dispatchMessage.
func TestDispatchMessageRouting(t *testing.T) {
	b := newMockBot(broadcastInvoker())
	ctx := context.Background()

	var fired []string
	b.OnMessage(func(c *Context) error { fired = append(fired, "message"); return nil })
	b.OnEditedMessage(func(c *Context) error { fired = append(fired, "edited"); return nil })
	b.OnChannelPost(func(c *Context) error { fired = append(fired, "channel"); return nil })

	regular := &tg.Message{ID: 1, Message: "x", PeerID: &tg.PeerUser{UserID: 10}}
	regular.SetFromID(&tg.PeerUser{UserID: 10})
	b.dispatchMessage(ctx, regular, false)
	b.dispatchMessage(ctx, regular, true)

	channelMsg := &tg.Message{ID: 2, Message: "x", PeerID: &tg.PeerChannel{ChannelID: 50}}
	b.dispatchMessage(ctx, channelMsg, false)
	b.dispatchMessage(ctx, channelMsg, true) // edited channel post (no registrar; exercises the switch)

	// A service message converts to nil and routes nothing.
	b.dispatchMessage(ctx, &tg.MessageService{ID: 3, PeerID: &tg.PeerUser{UserID: 10}}, false)

	want := []string{"message", "edited", "channel"}
	if len(fired) != len(want) {
		t.Fatalf("fired = %v, want %v", fired, want)
	}
	for i := range want {
		if fired[i] != want[i] {
			t.Fatalf("fired = %v, want %v", fired, want)
		}
	}
}

// TestInstallHandlers wires the dispatcher and feeds it a new-message update,
// covering the installHandlers callbacks.
func TestInstallHandlers(t *testing.T) {
	b := newMockBot(newMockInvoker())
	b.installHandlers()
	fired := false
	b.OnMessage(func(c *Context) error { fired = true; return nil })

	msg := &tg.Message{ID: 1, Message: "hi", PeerID: &tg.PeerUser{UserID: 10}}
	msg.SetFromID(&tg.PeerUser{UserID: 10})
	if err := b.disp.Handle(context.Background(), &tg.Updates{
		Updates: []tg.UpdateClass{&tg.UpdateNewMessage{Message: msg}},
		Users:   []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	}); err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if !fired {
		t.Fatal("new-message handler did not fire through the dispatcher")
	}
}

// TestPollFromTgQuiz covers the quiz, correct-option and solution branches of
// pollFromTg.
func TestPollFromTgQuiz(t *testing.T) {
	poll := &tg.Poll{
		ID:       7,
		Closed:   true,
		Quiz:     true,
		Question: tg.TextWithEntities{Text: "2+2?"},
		Answers: []tg.PollAnswerClass{
			&tg.PollAnswer{Text: tg.TextWithEntities{Text: "3"}, Option: []byte{0}},
			&tg.PollAnswer{Text: tg.TextWithEntities{Text: "4"}, Option: []byte{1}},
		},
	}
	results := &tg.PollResults{
		TotalVoters: 5,
		Results: []tg.PollAnswerVoters{
			{Option: []byte{0}, Voters: 2},
			{Option: []byte{1}, Voters: 3, Correct: true},
		},
	}
	results.SetSolution("two plus two")
	out := pollFromTg(poll, results)
	if out.Type != PollQuiz || !out.IsClosed || out.CorrectOptionID != 1 {
		t.Fatalf("poll = %#v", out)
	}
	if out.Explanation != "two plus two" || out.TotalVoterCount != 5 {
		t.Fatalf("poll explanation/total = %#v", out)
	}
}

// TestChatPhotoBranches covers the non-upload and private-chat rejection paths
// of SetChatPhoto/DeleteChatPhoto.
func TestChatPhotoBranches(t *testing.T) {
	b := newMockBot(newMockInvoker())
	ctx := context.Background()

	if err := b.SetChatPhoto(ctx, userRef(10, 20), FileID("x")); err == nil {
		t.Fatal("SetChatPhoto with non-upload file should fail")
	}
	if err := b.SetChatPhoto(ctx, userRef(10, 20), FileFromBytes("p.jpg", []byte("img"))); err == nil {
		t.Fatal("SetChatPhoto on a private chat should fail")
	}
	if err := b.DeleteChatPhoto(ctx, userRef(10, 20)); err == nil {
		t.Fatal("DeleteChatPhoto on a private chat should fail")
	}
}
