package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetManagedBotAccessSettings(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsGetAccessSettingsRequestTypeID, &tg.BotsAccessSettings{
		Restricted: true,
		AddUsers:   []tg.UserClass{&tg.User{ID: 11}, &tg.User{ID: 12}},
	})

	got, err := newMockBot(inv).GetManagedBotAccessSettings(context.Background(), 99)
	if err != nil {
		t.Fatalf("GetManagedBotAccessSettings: %v", err)
	}

	if !got.IsAccessRestricted || len(got.AddedUserIDs) != 2 {
		t.Fatalf("settings = %#v", got)
	}

	var req tg.BotsGetAccessSettingsRequest

	inv.decode(t, tg.BotsGetAccessSettingsRequestTypeID, &req)

	if _, ok := req.Bot.(*tg.InputUser); !ok {
		t.Fatalf("bot = %#v", req.Bot)
	}
}

func TestSetManagedBotAccessSettings(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsEditAccessSettingsRequestTypeID, &tg.BoolTrue{})

	if err := newMockBot(inv).SetManagedBotAccessSettings(context.Background(), 99, true, []int64{11, 12}); err != nil {
		t.Fatalf("SetManagedBotAccessSettings: %v", err)
	}

	var req tg.BotsEditAccessSettingsRequest

	inv.decode(t, tg.BotsEditAccessSettingsRequestTypeID, &req)

	if !req.Restricted || len(req.AddUsers) != 2 {
		t.Fatalf("req = %#v", req)
	}
}

func TestGetManagedBotToken(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsExportBotTokenRequestTypeID, &tg.BotsExportedBotToken{Token: "123:ABC"})

	tok, err := newMockBot(inv).GetManagedBotToken(context.Background(), 99)
	if err != nil {
		t.Fatalf("GetManagedBotToken: %v", err)
	}

	if tok != "123:ABC" {
		t.Fatalf("token = %q", tok)
	}

	var req tg.BotsExportBotTokenRequest

	inv.decode(t, tg.BotsExportBotTokenRequestTypeID, &req)

	if req.Revoke {
		t.Fatal("get token should not revoke")
	}
}

func TestReplaceManagedBotToken(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.BotsExportBotTokenRequestTypeID, &tg.BotsExportedBotToken{Token: "456:NEW"})

	tok, err := newMockBot(inv).ReplaceManagedBotToken(context.Background(), 99)
	if err != nil {
		t.Fatalf("ReplaceManagedBotToken: %v", err)
	}

	if tok != "456:NEW" {
		t.Fatalf("token = %q", tok)
	}

	var req tg.BotsExportBotTokenRequest

	inv.decode(t, tg.BotsExportBotTokenRequestTypeID, &req)

	if !req.Revoke {
		t.Fatal("replace token should revoke")
	}
}
