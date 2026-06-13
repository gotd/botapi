package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestUserFromTgUser(t *testing.T) {
	u := &tg.User{ID: 42, Bot: true, FirstName: "Ada", LastName: "L", Username: "ada", LangCode: "en"}
	u.SetPremium(true)
	got := userFromTgUser(u)
	want := User{ID: 42, IsBot: true, FirstName: "Ada", LastName: "L", Username: "ada", LanguageCode: "en", IsPremium: true}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestInlineKeyboardFromTg(t *testing.T) {
	mkp := &tg.ReplyInlineMarkup{Rows: []tg.KeyboardButtonRow{{Buttons: []tg.KeyboardButtonClass{
		&tg.KeyboardButtonURL{Text: "site", URL: "https://example.com"},
		&tg.KeyboardButtonCallback{Text: "ok", Data: []byte("ok:1")},
		&tg.KeyboardButtonSwitchInline{Text: "here", Query: "q", SamePeer: true},
	}}}}

	got := inlineKeyboardFromTg(mkp)
	if len(got.InlineKeyboard) != 1 || len(got.InlineKeyboard[0]) != 3 {
		t.Fatalf("unexpected shape: %+v", got.InlineKeyboard)
	}
	row := got.InlineKeyboard[0]
	if row[0].URL != "https://example.com" || row[1].CallbackData != "ok:1" {
		t.Fatalf("url/callback lost: %+v", row)
	}
	if row[2].SwitchInlineQueryCurrentChat == nil || *row[2].SwitchInlineQueryCurrentChat != "q" {
		t.Fatalf("same-peer switch-inline lost: %+v", row[2])
	}
}
