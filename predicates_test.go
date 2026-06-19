package botapi

import "testing"

func TestCommandName(t *testing.T) {
	cases := map[string]struct {
		want   string
		target string
		ok     bool
	}{
		"/start":             {"start", "", true},
		"/start@mybot":       {"start", "mybot", true},
		"/help me please":    {"help", "", true},
		"/cmd@bot with args": {"cmd", "bot", true},
		"not a command":      {"", "", false},
		"/":                  {"", "", false},
	}
	for text, want := range cases {
		got, target, ok := commandName(text)
		if got != want.want || target != want.target || ok != want.ok {
			t.Fatalf("commandName(%q) = (%q, %q, %v), want (%q, %q, %v)",
				text, got, target, ok, want.want, want.target, want.ok)
		}
	}
}

func TestCommandPredicate(t *testing.T) {
	// Untargeted command always matches.
	plain := &Context{Update: &Update{Message: &Message{Text: "/start hi"}}}
	if !Command("start")(plain) || !Command("/start")(plain) {
		t.Fatal("Command should match with and without slash")
	}

	if Command("help")(plain) {
		t.Fatal("Command should not match a different command")
	}

	nonMsg := &Context{Update: &Update{CallbackQuery: &CallbackQuery{}}}
	if Command("start")(nonMsg) {
		t.Fatal("Command should not match a non-message update")
	}

	// Targeted command matches only when the @target is this bot (case-insensitive).
	// Если поле botUsername находится в Update:
	mine := &Context{Update: &Update{Message: &Message{Text: "/start@MyBot hi"}, botUsername: "mybot"}}
	if !Command("start")(mine) {
		t.Fatal("Command should match when targeted at this bot")
	}

	other := &Context{Update: &Update{Message: &Message{Text: "/start@other_bot hi"}, botUsername: "mybot"}}
	if Command("start")(other) {
		t.Fatal("Command should not match when targeted at another bot")
	}

	// Targeted command with an unknown bot username does not match.
	unknown := &Context{Update: &Update{Message: &Message{Text: "/start@mybot hi"}}}
	if Command("start")(unknown) {
		t.Fatal("Command should not match a targeted command when the bot username is unknown")
	}
}

func TestTextAndChatPredicates(t *testing.T) {
	c := &Context{Update: &Update{Message: &Message{Text: "hello world", Chat: Chat{Type: ChatTypePrivate}}}}
	if !HasPrefix("hello")(c) || !HasText()(c) || !Regex(`^hello`)(c) {
		t.Fatal("text predicates should match")
	}

	if !ChatTypeIs(ChatTypePrivate)(c) || ChatTypeIs(ChatTypeChannel)(c) {
		t.Fatal("ChatTypeIs mismatch")
	}

	if !Not(TextEquals("nope"))(c) {
		t.Fatal("Not should invert")
	}
}

func TestCallbackPredicates(t *testing.T) {
	c := &Context{Update: &Update{CallbackQuery: &CallbackQuery{Data: "vote:42"}}}
	if !CallbackPrefix("vote:")(c) || !CallbackData("vote:42")(c) {
		t.Fatal("callback predicates should match")
	}

	if !Or(CallbackData("x"), CallbackPrefix("vote:"))(c) {
		t.Fatal("Or should match when one matches")
	}

	if c.Update.Text() != "vote:42" {
		t.Fatalf("Update.Text for callback = %q", c.Update.Text())
	}
}
