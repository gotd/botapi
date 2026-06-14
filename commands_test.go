package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestBotCommandScopeResolve(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		scope BotCommandScope
		want  tg.BotCommandScopeClass
	}{
		{BotCommandScopeDefault(), &tg.BotCommandScopeDefault{}},
		{BotCommandScopeAllPrivateChats(), &tg.BotCommandScopeUsers{}},
		{BotCommandScopeAllGroupChats(), &tg.BotCommandScopeChats{}},
		{BotCommandScopeAllChatAdministrators(), &tg.BotCommandScopeChatAdmins{}},
	}

	for _, c := range cases {
		// The non-targeted variants ignore the *Bot receiver, so nil is fine.
		got, err := c.scope.resolve(ctx, nil)
		if err != nil {
			t.Fatalf("resolve %T: %v", c.scope, err)
		}

		if got.TypeID() != c.want.TypeID() {
			t.Fatalf("resolve %T: got %T, want %T", c.scope, got, c.want)
		}
	}
}

func TestCommandConfigDefaultScope(t *testing.T) {
	var cfg commandConfig

	scope, err := cfg.resolveScope(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := scope.(*tg.BotCommandScopeDefault); !ok {
		t.Fatalf("default scope: got %T", scope)
	}
}
