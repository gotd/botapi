package botapi

import (
	"encoding/json"
	"testing"
)

func TestChatIDMarshal(t *testing.T) {
	got, err := json.Marshal(struct {
		ChatID ChatID `json:"chat_id"`
	}{ChatID: ID(-100123)})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"chat_id":-100123}` {
		t.Fatalf("numeric id: got %s", got)
	}

	got, err = json.Marshal(struct {
		ChatID ChatID `json:"chat_id"`
	}{ChatID: Username("@durov")})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"chat_id":"@durov"}` {
		t.Fatalf("username: got %s", got)
	}
}

func TestInlineKeyboardBuilder(t *testing.T) {
	kb := InlineKeyboard(
		InlineRow(InlineButtonURL("site", "https://example.com"), InlineButtonData("ok", "ok:1")),
	)
	if len(kb.InlineKeyboard) != 1 || len(kb.InlineKeyboard[0]) != 2 {
		t.Fatalf("unexpected shape: %+v", kb.InlineKeyboard)
	}
	if kb.InlineKeyboard[0][1].CallbackData != "ok:1" {
		t.Fatalf("callback data not set: %+v", kb.InlineKeyboard[0][1])
	}

	// The builder result must satisfy the sealed ReplyMarkup union.
	var _ ReplyMarkup = kb
}

func TestInputFileConstructors(t *testing.T) {
	if _, ok := FileID("abc").(InputFileID); !ok {
		t.Fatal("FileID should yield InputFileID")
	}
	up, ok := FileFromBytes("a.txt", []byte("hi")).(*InputFileUpload)
	if !ok || up.Name != "a.txt" || string(up.Bytes) != "hi" {
		t.Fatalf("FileFromBytes: %+v ok=%v", up, ok)
	}
}
