package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestCallbackQueryFromTg(t *testing.T) {
	e := tg.Entities{Users: map[int64]*tg.User{5: {ID: 5, FirstName: "Ada", Username: "ada"}}}
	u := &tg.UpdateBotCallbackQuery{QueryID: 100, UserID: 5, ChatInstance: 77, Data: []byte("vote:1")}

	cq := callbackQueryFromTg(e, u)
	if cq.ID != "100" || cq.ChatInstance != "77" || cq.Data != "vote:1" {
		t.Fatalf("scalar fields: %+v", cq)
	}

	if cq.From.ID != 5 || cq.From.Username != "ada" {
		t.Fatalf("from from entities: %+v", cq.From)
	}
}

func TestCallbackQueryFromTgUnknownUser(t *testing.T) {
	cq := callbackQueryFromTg(tg.Entities{}, &tg.UpdateBotCallbackQuery{QueryID: 1, UserID: 9})
	if cq.From.ID != 9 {
		t.Fatalf("fallback from id: %+v", cq.From)
	}
}

func TestInlineQueryFromTg(t *testing.T) {
	e := tg.Entities{Users: map[int64]*tg.User{3: {ID: 3, FirstName: "Bob"}}}
	u := &tg.UpdateBotInlineQuery{QueryID: 42, UserID: 3, Query: "cats", Offset: "10"}

	iq := inlineQueryFromTg(e, u)
	if iq.ID != "42" || iq.Query != "cats" || iq.Offset != "10" || iq.From.ID != 3 {
		t.Fatalf("inline query: %+v", iq)
	}
}
