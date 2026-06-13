package botapi

import "testing"

func TestOnCommandCollectsCommands(t *testing.T) {
	b := newTestBot(t)
	noop := func(*Context) error { return nil }

	b.OnCommand("start", "Start the bot", noop)
	b.OnCommand("/help", "Show help", noop) // leading slash is stripped
	b.Group().OnCommand("ban", "Ban a user", noop)
	b.OnCommand("start", "duplicate", noop) // dedup by name, keep first description

	if len(b.commands) != 3 {
		t.Fatalf("want 3 commands, got %d: %#v", len(b.commands), b.commands)
	}
	want := map[string]string{
		"start": "Start the bot",
		"help":  "Show help",
		"ban":   "Ban a user",
	}
	for _, c := range b.commands {
		if want[c.Command] != c.Description {
			t.Fatalf("command %q: got %q, want %q", c.Command, c.Description, want[c.Command])
		}
	}
}

func TestRegisterCommandIgnoresEmpty(t *testing.T) {
	b := newTestBot(t)
	b.registerCommand("", "no name")
	b.registerCommand("/", "slash only")
	if len(b.commands) != 0 {
		t.Fatalf("empty command names should be ignored, got %#v", b.commands)
	}
}

func TestDisableCommandRegistration(t *testing.T) {
	b, err := New("123:abc", Options{AppID: 1, AppHash: "h", DisableCommandRegistration: true})
	if err != nil {
		t.Fatal(err)
	}
	if b.registerCommands {
		t.Fatal("registerCommands should be false when disabled")
	}
}
