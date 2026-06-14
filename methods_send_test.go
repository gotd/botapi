package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// editUpdates wraps a single edited message in an Updates, the shape edit
// methods unpack.
func editUpdates(msg *tg.Message) *tg.Updates {
	return &tg.Updates{
		Updates: []tg.UpdateClass{&tg.UpdateEditMessage{Message: msg}},
		Users:   []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
	}
}

func TestSendDice(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SendDice(context.Background(), userRef(10, 20), DiceDie); err != nil {
		t.Fatalf("SendDice: %v", err)
	}
	if !inv.called(tg.MessagesSendMediaRequestTypeID) {
		t.Fatal("expected messages.sendMedia")
	}
}

func TestSendLocation(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SendLocation(context.Background(), userRef(10, 20), 55.75, 37.61); err != nil {
		t.Fatalf("SendLocation: %v", err)
	}
}

func TestSendVenue(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SendVenue(context.Background(), userRef(10, 20), 55.75, 37.61, "Red Square", "Moscow"); err != nil {
		t.Fatalf("SendVenue: %v", err)
	}
}

func TestSendContact(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SendContact(context.Background(), userRef(10, 20), "+1555", "Ada", "Lovelace"); err != nil {
		t.Fatalf("SendContact: %v", err)
	}
}

func TestSendPoll(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSendMediaRequestTypeID, messageUpdates(&tg.Message{ID: 1, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SendPoll(context.Background(), userRef(10, 20), "Q?", []string{"a", "b"}); err != nil {
		t.Fatalf("SendPoll: %v", err)
	}
}

func TestEditMessageText(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{
		ID:      5,
		Message: "edited",
		PeerID:  &tg.PeerUser{UserID: 10},
	}))
	b := newMockBot(inv)

	m, err := b.EditMessageText(context.Background(), userRef(10, 20), 5, "edited")
	if err != nil {
		t.Fatalf("EditMessageText: %v", err)
	}
	if m.Text != "edited" {
		t.Fatalf("text = %q", m.Text)
	}
	var req tg.MessagesEditMessageRequest
	inv.decode(t, tg.MessagesEditMessageRequestTypeID, &req)
	if req.ID != 5 || req.Message != "edited" {
		t.Fatalf("req = %#v", req)
	}
}

func TestEditMessageCaption(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.EditMessageCaption(context.Background(), userRef(10, 20), 5, "cap"); err != nil {
		t.Fatalf("EditMessageCaption: %v", err)
	}
}

func TestEditMessageReplyMarkup(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	kb := InlineKeyboard([]InlineKeyboardButton{InlineButtonData("ok", "data")})
	if _, err := b.EditMessageReplyMarkup(context.Background(), userRef(10, 20), 5, kb); err != nil {
		t.Fatalf("EditMessageReplyMarkup: %v", err)
	}
	var req tg.MessagesEditMessageRequest
	inv.decode(t, tg.MessagesEditMessageRequestTypeID, &req)
	if _, ok := req.ReplyMarkup.(*tg.ReplyInlineMarkup); !ok {
		t.Fatalf("markup = %#v", req.ReplyMarkup)
	}
}

func TestEditMessageLiveLocation(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.EditMessageLiveLocation(context.Background(), userRef(10, 20), 5, 55.75, 37.61); err != nil {
		t.Fatalf("EditMessageLiveLocation: %v", err)
	}
}

func TestStopMessageLiveLocation(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.StopMessageLiveLocation(context.Background(), userRef(10, 20), 5, nil); err != nil {
		t.Fatalf("StopMessageLiveLocation: %v", err)
	}
}

func TestSetGameScore(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetGameScoreRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))
	b := newMockBot(inv)

	if _, err := b.SetGameScore(context.Background(), userRef(10, 20), 5, 99, 1000); err != nil {
		t.Fatalf("SetGameScore: %v", err)
	}
	var req tg.MessagesSetGameScoreRequest
	inv.decode(t, tg.MessagesSetGameScoreRequestTypeID, &req)
	if req.Score != 1000 {
		t.Fatalf("score = %d", req.Score)
	}
}

func TestGetGameHighScores(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesGetGameHighScoresRequestTypeID, &tg.MessagesHighScores{
		Scores: []tg.HighScore{{Pos: 1, UserID: 99, Score: 1000}},
		Users:  []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})
	b := newMockBot(inv)

	scores, err := b.GetGameHighScores(context.Background(), userRef(10, 20), 5, 99)
	if err != nil {
		t.Fatalf("GetGameHighScores: %v", err)
	}
	if len(scores) != 1 || scores[0].Score != 1000 {
		t.Fatalf("scores = %#v", scores)
	}
}
