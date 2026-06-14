package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSendChatActionAll(t *testing.T) {
	actions := []ChatAction{
		ChatActionTyping, ChatActionUploadPhoto, ChatActionRecordVideo, ChatActionUploadVideo,
		ChatActionRecordVoice, ChatActionUploadVoice, ChatActionUploadDocument, ChatActionChooseSticker,
		ChatActionFindLocation, ChatActionRecordVideoNote, ChatActionUploadVideoNote,
	}
	for _, action := range actions {
		inv := newMockInvoker()
		inv.reply(tg.MessagesSetTypingRequestTypeID, &tg.BoolTrue{})

		b := newMockBot(inv)
		if err := b.SendChatAction(context.Background(), userRef(10, 20), action); err != nil {
			t.Fatalf("SendChatAction(%q): %v", action, err)
		}
	}

	// Unknown action is rejected before the wire.
	b := newMockBot(newMockInvoker())
	if err := b.SendChatAction(context.Background(), userRef(10, 20), ChatAction("nonsense")); err == nil {
		t.Fatal("unknown action should error")
	}
}
