package botapi

import "testing"

func TestCommandName(t *testing.T) {
	cases := map[string]struct {
		want string
		ok   bool
	}{
		"/start":             {"start", true},
		"/start@mybot":       {"start", true},
		"/help me please":    {"help", true},
		"/cmd@bot with args": {"cmd", true},
		"not a command":      {"", false},
		"/":                  {"", false},
	}
	for text, want := range cases {
		got, ok := commandName(text)
		if got != want.want || ok != want.ok {
			t.Fatalf("commandName(%q) = (%q, %v), want (%q, %v)", text, got, ok, want.want, want.ok)
		}
	}
}

func TestCommandPredicate(t *testing.T) {
	msg := &Update{Message: &Message{Text: "/start@mybot hi"}}
	if !Command("start")(msg) || !Command("/start")(msg) {
		t.Fatal("Command should match with and without slash")
	}
	if Command("help")(msg) {
		t.Fatal("Command should not match a different command")
	}
	if Command("start")(&Update{CallbackQuery: &CallbackQuery{}}) {
		t.Fatal("Command should not match a non-message update")
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
