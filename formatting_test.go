package botapi

import "testing"

func TestMessage_OriginalMD(t *testing.T) {
	msg := &Message{
		Text: "hello world",
		Entities: []MessageEntity{
			{
				Type:   EntityBold,
				Offset: 6,
				Length: 5,
			},
		},
	}

	got := msg.OriginalMD()
	want := "hello *world*"

	if got != want {
		t.Fatalf("OriginalMD() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMDV2(t *testing.T) {
	msg := &Message{
		Text: "hello world",
		Entities: []MessageEntity{
			{
				Type:   EntityUnderline,
				Offset: 6,
				Length: 5,
			},
		},
	}

	got := msg.OriginalMDV2()
	want := "hello __world__"

	if got != want {
		t.Fatalf("OriginalMDV2() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalHTML(t *testing.T) {
	msg := &Message{
		Text: "hello world",
		Entities: []MessageEntity{
			{
				Type:   EntityItalic,
				Offset: 6,
				Length: 5,
			},
		},
	}

	got := msg.OriginalHTML()
	want := "hello <i>world</i>"

	if got != want {
		t.Fatalf("OriginalHTML() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalHTML_EscapesText(t *testing.T) {
	msg := &Message{
		Text: "<b>hello</b>",
	}

	got := msg.OriginalHTML()
	want := "&lt;b&gt;hello&lt;/b&gt;"

	if got != want {
		t.Fatalf("OriginalHTML() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMD_TextLink(t *testing.T) {
	msg := &Message{
		Text: "OpenAI",
		Entities: []MessageEntity{
			{
				Type:   EntityTextLink,
				Offset: 0,
				Length: 6,
				URL:    "https://openai.com",
			},
		},
	}

	got := msg.OriginalMD()
	want := "[OpenAI](https://openai.com)"

	if got != want {
		t.Fatalf("OriginalMD() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMD_TextMention(t *testing.T) {
	msg := &Message{
		Text: "John",
		Entities: []MessageEntity{
			{
				Type:   EntityTextMention,
				Offset: 0,
				Length: 4,
				User: &User{
					ID: 12345,
				},
			},
		},
	}

	got := msg.OriginalMD()
	want := "[John](tg://user?id=12345)"

	if got != want {
		t.Fatalf("OriginalMD() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMD_PreWithLanguage(t *testing.T) {
	msg := &Message{
		Text: "fmt.Println()",
		Entities: []MessageEntity{
			{
				Type:     EntityPre,
				Offset:   0,
				Length:   13,
				Language: "go",
			},
		},
	}

	got := msg.OriginalMD()
	want := "```go\nfmt.Println()```"

	if got != want {
		t.Fatalf("OriginalMD() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMD_EscapeCodeMarker(t *testing.T) {
	msg := &Message{
		Text: "`test`",
		Entities: []MessageEntity{
			{
				Type:   EntityCode,
				Offset: 0,
				Length: 6,
			},
		},
	}

	got := msg.OriginalMD()
	want := "`\\`test\\``"

	if got != want {
		t.Fatalf("OriginalMD() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalMDV2_NestedEntities(t *testing.T) {
	msg := &Message{
		Text: "hello world",
		Entities: []MessageEntity{
			{
				Type:   EntityBold,
				Offset: 0,
				Length: 11,
			},
			{
				Type:   EntityItalic,
				Offset: 6,
				Length: 5,
			},
		},
	}

	got := msg.OriginalMDV2()
	want := "*hello _world_*"

	if got != want {
		t.Fatalf("OriginalMDV2() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalHTML_NestedEntities(t *testing.T) {
	msg := &Message{
		Text: "hello world",
		Entities: []MessageEntity{
			{
				Type:   EntityBold,
				Offset: 0,
				Length: 11,
			},
			{
				Type:   EntityItalic,
				Offset: 6,
				Length: 5,
			},
		},
	}

	got := msg.OriginalHTML()
	want := "<b>hello <i>world</i></b>"

	if got != want {
		t.Fatalf("OriginalHTML() = %q, want %q", got, want)
	}
}

func TestMessage_OriginalTextMD_UsesCaptionWhenTextEmpty(t *testing.T) {
	msg := &Message{
		Caption: "caption",
		CaptionEntities: []MessageEntity{
			{
				Type:   EntityBold,
				Offset: 0,
				Length: 7,
			},
		},
	}

	got := msg.OriginalTextMD()
	want := "*caption*"

	if got != want {
		t.Fatalf("OriginalTextMD() = %q, want %q", got, want)
	}
}
