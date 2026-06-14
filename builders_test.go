package botapi

import (
	"bytes"
	"testing"
)

func TestBuilders(t *testing.T) {
	if kb := Keyboard(Row(Button("a"), ButtonContact("c")), Row(Button("b"))); len(kb.Keyboard) != 2 {
		t.Fatalf("keyboard = %#v", kb)
	}

	if ik := InlineKeyboard([]InlineKeyboardButton{InlineButtonData("x", "d"), InlineButtonURL("u", "https://e")}); len(ik.InlineKeyboard) != 1 {
		t.Fatalf("inline keyboard = %#v", ik)
	}

	for _, f := range []InputFile{
		FileID("id"), FileURL("https://e"), FileFromPath("/tmp/x"),
		FileFromBytes("n", []byte("d")), FileFromReader("n", bytes.NewReader(nil)),
	} {
		if f == nil {
			t.Fatal("nil input file")
		}
	}

	if e := Emoji("👍"); e == nil {
		t.Fatal("nil emoji reaction")
	}

	if e := CustomEmoji("123"); e == nil {
		t.Fatal("nil custom emoji reaction")
	}
}
