package botapi

import (
	"context"
	"testing"
)

func TestContextChatAndSender(t *testing.T) {
	c := &Context{
		Context: context.Background(),
		Update: &Update{Message: &Message{
			MessageID: 7,
			Chat:      Chat{ID: -100500, Type: ChatTypeSupergroup},
			From:      &User{ID: 11, Username: "ada"},
		}},
	}

	chat, ok := c.Chat()
	if !ok {
		t.Fatal("Chat should be present")
	}
	if id, isInt := chat.(ChatIDInt); !isInt || int64(id) != -100500 {
		t.Fatalf("Chat id = %v", chat)
	}
	if s := c.Sender(); s == nil || s.ID != 11 {
		t.Fatalf("Sender = %+v", c.Sender())
	}
}

func TestContextSenderFromCallback(t *testing.T) {
	c := &Context{
		Context: context.Background(),
		Update:  &Update{CallbackQuery: &CallbackQuery{From: User{ID: 99}}},
	}
	if s := c.Sender(); s == nil || s.ID != 99 {
		t.Fatalf("callback sender = %+v", c.Sender())
	}
}

func TestContextReplyWithoutMessage(t *testing.T) {
	c := &Context{Context: context.Background(), Bot: newTestBot(t), Update: &Update{InlineQuery: &InlineQuery{}}}
	if _, err := c.Reply("hi"); err == nil {
		t.Fatal("Reply with no message should error")
	}
	if _, ok := c.Chat(); ok {
		t.Fatal("inline-query-only update should have no chat")
	}
}
