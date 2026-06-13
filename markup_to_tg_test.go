package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestReplyMarkupToTg_Inline(t *testing.T) {
	q := "go"
	m := &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{
		{Text: "site", URL: "https://example.com"},
		{Text: "ok", CallbackData: "ok:1"},
		{Text: "switch", SwitchInlineQuery: &q},
		{Text: "pay", Pay: true},
	}}}

	got, err := replyMarkupToTg(m)
	if err != nil {
		t.Fatal(err)
	}
	inline, isInline := got.(*tg.ReplyInlineMarkup)
	if !isInline || len(inline.Rows) != 1 || len(inline.Rows[0].Buttons) != 4 {
		t.Fatalf("unexpected: %#v", got)
	}
	if _, isURL := inline.Rows[0].Buttons[0].(*tg.KeyboardButtonURL); !isURL {
		t.Fatalf("button 0: want URL, got %T", inline.Rows[0].Buttons[0])
	}
	cb, isCb := inline.Rows[0].Buttons[1].(*tg.KeyboardButtonCallback)
	if !isCb || string(cb.Data) != "ok:1" {
		t.Fatalf("button 1: want callback ok:1, got %#v", inline.Rows[0].Buttons[1])
	}
	si, isSwitch := inline.Rows[0].Buttons[2].(*tg.KeyboardButtonSwitchInline)
	if !isSwitch || si.Query != "go" || si.SamePeer {
		t.Fatalf("button 2: want switch-inline, got %#v", inline.Rows[0].Buttons[2])
	}
}

func TestReplyMarkupToTg_InlineTextButtonRejected(t *testing.T) {
	m := &InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "plain"}}}}
	if _, err := replyMarkupToTg(m); err == nil {
		t.Fatal("plain text inline button must be rejected")
	}
}

func TestReplyMarkupToTg_Reply(t *testing.T) {
	m := &ReplyKeyboardMarkup{
		Keyboard:        [][]KeyboardButton{{Button("hi"), ButtonContact("phone"), ButtonLocation("where")}},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	got, err := replyMarkupToTg(m)
	if err != nil {
		t.Fatal(err)
	}
	kb, ok := got.(*tg.ReplyKeyboardMarkup)
	if !ok || !kb.Resize || !kb.SingleUse {
		t.Fatalf("unexpected: %#v", got)
	}
	if _, ok := kb.Rows[0].Buttons[1].(*tg.KeyboardButtonRequestPhone); !ok {
		t.Fatalf("want request-phone, got %T", kb.Rows[0].Buttons[1])
	}
	if _, ok := kb.Rows[0].Buttons[2].(*tg.KeyboardButtonRequestGeoLocation); !ok {
		t.Fatalf("want request-geo, got %T", kb.Rows[0].Buttons[2])
	}
}

func TestReplyMarkupToTg_RemoveAndForceReply(t *testing.T) {
	if _, err := replyMarkupToTg(RemoveKeyboard()); err != nil {
		t.Fatal(err)
	}
	got, err := replyMarkupToTg(&ForceReply{ForceReply: true, InputFieldPlaceholder: "type"})
	if err != nil {
		t.Fatal(err)
	}
	fr, ok := got.(*tg.ReplyKeyboardForceReply)
	if !ok || fr.Placeholder != "type" {
		t.Fatalf("unexpected force reply: %#v", got)
	}
}
