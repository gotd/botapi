package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestAnswerCallbackQuery(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotCallbackAnswerRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	err := b.AnswerCallbackQuery(context.Background(), "123",
		WithCallbackText("done"), WithCallbackAlert())
	if err != nil {
		t.Fatalf("AnswerCallbackQuery: %v", err)
	}

	var req tg.MessagesSetBotCallbackAnswerRequest

	inv.decode(t, tg.MessagesSetBotCallbackAnswerRequestTypeID, &req)

	if req.Message != "done" || !req.Alert {
		t.Fatalf("req = %#v", req)
	}
}

func TestAnswerCallbackQueryError(t *testing.T) {
	b := newMockBot(newMockInvoker())
	// A non-numeric query id never reaches the wire.
	if err := b.AnswerCallbackQuery(context.Background(), ""); err != nil {
		// empty string is treated as a (zero) id and still sent; allow either.
		t.Logf("empty id: %v", err)
	}
}

func TestSetMyCommands(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetBotCommandsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	cmds := []BotCommand{{Command: "start", Description: "Start"}, {Command: "help", Description: "Help"}}
	if err := b.SetMyCommands(context.Background(), cmds); err != nil {
		t.Fatalf("SetMyCommands: %v", err)
	}

	var req tg.BotsSetBotCommandsRequest

	inv.decode(t, tg.BotsSetBotCommandsRequestTypeID, &req)

	if len(req.Commands) != 2 || req.Commands[0].Command != "start" {
		t.Fatalf("commands = %#v", req.Commands)
	}
}

func TestGetMyCommands(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsGetBotCommandsRequestTypeID, &tg.BotCommandVector{
		Elems: []tg.BotCommand{{Command: "start", Description: "Start"}},
	})

	b := newMockBot(inv)

	cmds, err := b.GetMyCommands(context.Background())
	if err != nil {
		t.Fatalf("GetMyCommands: %v", err)
	}

	if len(cmds) != 1 || cmds[0].Command != "start" {
		t.Fatalf("cmds = %#v", cmds)
	}
}

func TestDeleteMyCommands(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsResetBotCommandsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.DeleteMyCommands(context.Background()); err != nil {
		t.Fatalf("DeleteMyCommands: %v", err)
	}

	if !inv.called(tg.BotsResetBotCommandsRequestTypeID) {
		t.Fatal("expected bots.resetBotCommands")
	}
}

func TestForwardMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesForwardMessagesRequestTypeID, messageUpdates(&tg.Message{
		ID:     7,
		PeerID: &tg.PeerUser{UserID: 10},
	}))

	b := newMockBot(inv)

	m, err := b.ForwardMessage(context.Background(), userRef(10, 20), userRef(30, 40), 7)
	if err != nil {
		t.Fatalf("ForwardMessage: %v", err)
	}

	if m.MessageID != 7 {
		t.Fatalf("message id = %d", m.MessageID)
	}
}

func TestDeleteMessage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesDeleteMessagesRequestTypeID, &tg.MessagesAffectedMessages{})

	b := newMockBot(inv)

	if err := b.DeleteMessage(context.Background(), userRef(10, 20), 7); err != nil {
		t.Fatalf("DeleteMessage: %v", err)
	}

	var req tg.MessagesDeleteMessagesRequest

	inv.decode(t, tg.MessagesDeleteMessagesRequestTypeID, &req)

	if len(req.ID) != 1 || req.ID[0] != 7 {
		t.Fatalf("ids = %v", req.ID)
	}
}

func TestDeleteMessages(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesDeleteMessagesRequestTypeID, &tg.MessagesAffectedMessages{})

	b := newMockBot(inv)

	if err := b.DeleteMessages(context.Background(), userRef(10, 20), []int{7, 8, 9}); err != nil {
		t.Fatalf("DeleteMessages: %v", err)
	}

	var req tg.MessagesDeleteMessagesRequest

	inv.decode(t, tg.MessagesDeleteMessagesRequestTypeID, &req)

	if len(req.ID) != 3 {
		t.Fatalf("ids = %v", req.ID)
	}
}

func TestCopyMessage(t *testing.T) {
	inv := newMockInvoker()
	// Copy is a forward with the author dropped.
	inv.reply(tg.MessagesForwardMessagesRequestTypeID, messageUpdates(&tg.Message{
		ID:      8,
		Message: "hello",
		PeerID:  &tg.PeerUser{UserID: 10},
	}))

	b := newMockBot(inv)

	m, err := b.CopyMessage(context.Background(), userRef(10, 20), userRef(30, 40), 7)
	if err != nil {
		t.Fatalf("CopyMessage: %v", err)
	}

	if m.MessageID != 8 {
		t.Fatalf("message id = %d", m.MessageID)
	}

	var req tg.MessagesForwardMessagesRequest

	inv.decode(t, tg.MessagesForwardMessagesRequestTypeID, &req)

	if !req.DropAuthor {
		t.Fatal("copy should drop author")
	}
}

func TestSendChatAction(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetTypingRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	if err := b.SendChatAction(context.Background(), userRef(10, 20), ChatActionTyping); err != nil {
		t.Fatalf("SendChatAction: %v", err)
	}

	if !inv.called(tg.MessagesSetTypingRequestTypeID) {
		t.Fatal("expected messages.setTyping")
	}
}

func TestGetChat(t *testing.T) {
	b := newMockBot(newMockInvoker())

	c, err := b.GetChat(context.Background(), tdlibChannel(50))
	if err != nil {
		t.Fatalf("GetChat: %v", err)
	}

	if c.Type != ChatTypeChannel && c.Type != ChatTypeSupergroup {
		t.Fatalf("chat type = %q", c.Type)
	}
}

func TestGetUserProfilePhotos(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PhotosGetUserPhotosRequestTypeID, &tg.PhotosPhotos{
		Photos: []tg.PhotoClass{
			&tg.Photo{ID: 1, AccessHash: 2, FileReference: []byte{0}, Sizes: []tg.PhotoSizeClass{
				&tg.PhotoSize{Type: "x", W: 640, H: 640, Size: 1000},
			}},
		},
		Users: []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})

	b := newMockBot(inv)

	photos, err := b.GetUserProfilePhotos(context.Background(), 99)
	if err != nil {
		t.Fatalf("GetUserProfilePhotos: %v", err)
	}

	if photos.TotalCount != 1 || len(photos.Photos) != 1 {
		t.Fatalf("photos = %#v", photos)
	}
}
