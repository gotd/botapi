package botapi

import "testing"

func TestSendOptions(t *testing.T) {
	var cfg sendConfig

	q := "x"

	for _, o := range []SendOption{
		Silent(),
		DisableWebPagePreview(),
		ProtectContent(),
		ReplyTo(99),
		WithParseMode(ParseModeHTML),
		WithReplyMarkup(&InlineKeyboardMarkup{InlineKeyboard: [][]InlineKeyboardButton{{{Text: "t", SwitchInlineQuery: &q}}}}),
	} {
		o(&cfg)
	}

	if !cfg.silent || !cfg.disableWebPreview || !cfg.protect {
		t.Fatalf("bool options not applied: %+v", cfg)
	}

	if cfg.replyTo != 99 || cfg.parseMode != ParseModeHTML || cfg.markup == nil {
		t.Fatalf("value options not applied: %+v", cfg)
	}
}

func TestStyledText(t *testing.T) {
	for _, mode := range []ParseMode{ParseModeNone, ParseModeHTML, ParseModeMarkdownV2, ParseModeMarkdown} {
		opts, err := styledText("hello", mode, nil)
		if err != nil || len(opts) != 1 {
			t.Fatalf("mode %q: got %d opts, err %v", mode, len(opts), err)
		}
	}

	if _, err := styledText("x", ParseMode("weird"), nil); err == nil {
		t.Fatal("unknown parse mode should error")
	}
}
