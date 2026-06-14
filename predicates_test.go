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
	plain := &Update{Message: &Message{Text: "/start hi"}}
	if !Command("start")(plain) || !Command("/start")(plain) {
		t.Fatal("Command should match with and without slash")
	}

	if Command("help")(plain) {
		t.Fatal("Command should not match a different command")
	}

	if Command("start")(&Update{CallbackQuery: &CallbackQuery{}}) {
		t.Fatal("Command should not match a non-message update")
	}

	// Targeted command matches only when the @target is this bot (case-insensitive).
	mine := &Update{Message: &Message{Text: "/start@MyBot hi"}, botUsername: "mybot"}
	if !Command("start")(mine) {
		t.Fatal("Command should match when targeted at this bot")
	}

	other := &Update{Message: &Message{Text: "/start@other_bot hi"}, botUsername: "mybot"}
	if Command("start")(other) {
		t.Fatal("Command should not match when targeted at another bot")
	}

	// Targeted command with an unknown bot username does not match.
	unknown := &Update{Message: &Message{Text: "/start@mybot hi"}}
	if Command("start")(unknown) {
		t.Fatal("Command should not match a targeted command when the bot username is unknown")
	}
}

func TestTextAndChatPredicates(t *testing.T) {
	u := &Update{Message: &Message{Text: "hello world", Chat: Chat{Type: ChatTypePrivate}}}
	if !HasPrefix("hello")(u) || !HasText()(u) || !Regex(`^hello`)(u) {
		t.Fatal("text predicates should match")
	}

	if !ChatTypeIs(ChatTypePrivate)(u) || ChatTypeIs(ChatTypeChannel)(u) {
		t.Fatal("ChatTypeIs mismatch")
	}

	if !Not(TextEquals("nope"))(u) {
		t.Fatal("Not should invert")
	}
}

func TestCallbackPredicates(t *testing.T) {
	u := &Update{CallbackQuery: &CallbackQuery{Data: "vote:42"}}
	if !CallbackPrefix("vote:")(u) || !CallbackData("vote:42")(u) {
		t.Fatal("callback predicates should match")
	}

	if !Or(CallbackData("x"), CallbackPrefix("vote:"))(u) {
		t.Fatal("Or should match when one matches")
	}

	if u.Text() != "vote:42" {
		t.Fatalf("Update.Text for callback = %q", u.Text())
	}
}
