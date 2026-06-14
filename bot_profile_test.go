package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetMyInfoFields(t *testing.T) {
	cases := []struct {
		name string
		call func(*Bot) error
		want func(*tg.BotsSetBotInfoRequest) (string, bool)
	}{
		{
			"name",
			func(b *Bot) error { return b.SetMyName(context.Background(), "Botty") },
			func(r *tg.BotsSetBotInfoRequest) (string, bool) { return r.GetName() },
		},
		{
			"description",
			func(b *Bot) error { return b.SetMyDescription(context.Background(), "long") },
			func(r *tg.BotsSetBotInfoRequest) (string, bool) { return r.GetDescription() },
		},
		{
			"short",
			func(b *Bot) error { return b.SetMyShortDescription(context.Background(), "short") },
			func(r *tg.BotsSetBotInfoRequest) (string, bool) { return r.GetAbout() },
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inv := newMockInvoker()
			inv.reply(tg.BotsSetBotInfoRequestTypeID, &tg.BoolTrue{})

			if err := c.call(newMockBot(inv)); err != nil {
				t.Fatalf("set: %v", err)
			}

			var req tg.BotsSetBotInfoRequest

			inv.decode(t, tg.BotsSetBotInfoRequestTypeID, &req)

			if _, ok := req.Bot.(*tg.InputUserSelf); !ok {
				t.Fatalf("bot = %#v, want self", req.Bot)
			}

			if _, ok := c.want(&req); !ok {
				t.Fatal("expected field to be set")
			}
		})
	}
}

func TestSetMyNameLanguage(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetBotInfoRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SetMyName(context.Background(), "Bot", WithBotInfoLanguage("ru")); err != nil {
		t.Fatalf("SetMyName: %v", err)
	}

	var req tg.BotsSetBotInfoRequest

	inv.decode(t, tg.BotsSetBotInfoRequestTypeID, &req)

	if req.LangCode != "ru" {
		t.Fatalf("lang = %q", req.LangCode)
	}
}

func TestGetMyInfo(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsGetBotInfoRequestTypeID, &tg.BotsBotInfo{
		Name:        "Botty",
		Description: "long",
		About:       "short",
	})

	b := newMockBot(inv)
	ctx := context.Background()

	if got, err := b.GetMyName(ctx); err != nil || got != "Botty" {
		t.Fatalf("GetMyName = %q, %v", got, err)
	}

	if got, err := b.GetMyDescription(ctx); err != nil || got != "long" {
		t.Fatalf("GetMyDescription = %q, %v", got, err)
	}

	if got, err := b.GetMyShortDescription(ctx); err != nil || got != "short" {
		t.Fatalf("GetMyShortDescription = %q, %v", got, err)
	}
}

func TestSetChatMenuButtonWebApp(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetBotMenuButtonRequestTypeID, &tg.BoolTrue{})

	button := MenuButtonWebApp{Type: MenuButtonWebAppType, Text: "Open", WebApp: WebAppInfo{URL: "https://e.x"}}
	if err := newMockBot(inv).SetChatMenuButton(context.Background(), button); err != nil {
		t.Fatalf("SetChatMenuButton: %v", err)
	}

	var req tg.BotsSetBotMenuButtonRequest

	inv.decode(t, tg.BotsSetBotMenuButtonRequestTypeID, &req)

	if _, ok := req.UserID.(*tg.InputUserEmpty); !ok {
		t.Fatalf("user = %#v, want empty (default)", req.UserID)
	}

	b, ok := req.Button.(*tg.BotMenuButton)
	if !ok || b.Text != "Open" || b.URL != "https://e.x" {
		t.Fatalf("button = %#v", req.Button)
	}
}

func TestGetChatMenuButtonDefault(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsGetBotMenuButtonRequestTypeID, &tg.BotMenuButtonCommands{})

	got, err := newMockBot(inv).GetChatMenuButton(context.Background())
	if err != nil {
		t.Fatalf("GetChatMenuButton: %v", err)
	}

	if _, ok := got.(MenuButtonCommands); !ok {
		t.Fatalf("button = %#v, want commands", got)
	}
}

func TestSetMyDefaultAdministratorRights(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsSetBotBroadcastDefaultAdminRightsRequestTypeID, &tg.BoolTrue{})
	inv.reply(tg.BotsSetBotGroupDefaultAdminRightsRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)
	rights := ChatAdminRights{CanPostMessages: true, CanDeleteMessages: true}

	if err := b.SetMyDefaultAdministratorRights(context.Background(), rights, true); err != nil {
		t.Fatalf("set channels: %v", err)
	}

	var ch tg.BotsSetBotBroadcastDefaultAdminRightsRequest

	inv.decode(t, tg.BotsSetBotBroadcastDefaultAdminRightsRequestTypeID, &ch)

	if !ch.AdminRights.PostMessages {
		t.Fatalf("rights = %#v", ch.AdminRights)
	}

	if err := b.SetMyDefaultAdministratorRights(context.Background(), rights, false); err != nil {
		t.Fatalf("set groups: %v", err)
	}

	if !inv.called(tg.BotsSetBotGroupDefaultAdminRightsRequestTypeID) {
		t.Fatal("group rights not set")
	}
}

func TestGetMyDefaultAdministratorRights(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.UsersGetFullUserRequestTypeID, &tg.UsersUserFull{
		FullUser: tg.UserFull{
			ID:                      1,
			BotGroupAdminRights:     tg.ChatAdminRights{DeleteMessages: true},
			BotBroadcastAdminRights: tg.ChatAdminRights{PostMessages: true},
		},
	})

	b := newMockBot(inv)
	ctx := context.Background()

	groups, err := b.GetMyDefaultAdministratorRights(ctx, false)
	if err != nil || !groups.CanDeleteMessages || groups.CanPostMessages {
		t.Fatalf("groups = %#v, %v", groups, err)
	}

	channels, err := b.GetMyDefaultAdministratorRights(ctx, true)
	if err != nil || !channels.CanPostMessages || channels.CanDeleteMessages {
		t.Fatalf("channels = %#v, %v", channels, err)
	}
}
