package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_convertToBotCommandScopeClass(t *testing.T) {
	newOASscope := func(typ oas.BotCommandScopeType) (r oas.OptBotCommandScope) {
		r.SetTo(oas.BotCommandScope{
			Type:                                 typ,
			BotCommandScopeDefault:               oas.BotCommandScopeDefault{},
			BotCommandScopeAllPrivateChats:       oas.BotCommandScopeAllPrivateChats{},
			BotCommandScopeAllGroupChats:         oas.BotCommandScopeAllGroupChats{},
			BotCommandScopeAllChatAdministrators: oas.BotCommandScopeAllChatAdministrators{},
			BotCommandScopeChat:                  oas.BotCommandScopeChat{},
			BotCommandScopeChatAdministrators:    oas.BotCommandScopeChatAdministrators{},
			BotCommandScopeChatMember:            oas.BotCommandScopeChatMember{},
		})
		return r
	}

	tests := []struct {
		name    string
		input   oas.OptBotCommandScope
		want    tg.BotCommandScopeClass
		wantErr bool
	}{
		{
			"Nil",
			oas.OptBotCommandScope{},
			&tg.BotCommandScopeDefault{},
			false,
		},
		{
			"",
			newOASscope(oas.BotCommandScopeDefaultBotCommandScope),
			&tg.BotCommandScopeDefault{},
			false,
		},
		{
			"",
			newOASscope(oas.BotCommandScopeAllPrivateChatsBotCommandScope),
			&tg.BotCommandScopeUsers{},
			false,
		},
		{
			"",
			newOASscope(oas.BotCommandScopeAllGroupChatsBotCommandScope),
			&tg.BotCommandScopeChats{},
			false,
		},
		{
			"",
			newOASscope(oas.BotCommandScopeAllChatAdministratorsBotCommandScope),
			&tg.BotCommandScopeChatAdmins{},
			false,
		},
		{
			"",
			oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeChatBotCommandScope,
				BotCommandScopeChat: oas.BotCommandScopeChat{
					ChatID: oas.NewInt64ID(testChatID()),
				},
			}),
			&tg.BotCommandScopePeer{Peer: testChat().AsInputPeer()},
			false,
		},
		{
			"",
			oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeChatAdministratorsBotCommandScope,
				BotCommandScopeChatAdministrators: oas.BotCommandScopeChatAdministrators{
					ChatID: oas.NewInt64ID(testChatID()),
				},
			}),
			&tg.BotCommandScopePeerAdmins{Peer: testChat().AsInputPeer()},
			false,
		},
		{
			"",
			oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeChatMemberBotCommandScope,
				BotCommandScopeChatMember: oas.BotCommandScopeChatMember{
					ChatID: oas.NewInt64ID(testChatID()),
					UserID: testUser().ID,
				},
			}),
			&tg.BotCommandScopePeerUser{
				Peer:   testChat().AsInputPeer(),
				UserID: &tg.InputUserSelf{},
			},
			false,
		},
		{
			"UnknownType",
			newOASscope("aboba"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		if tt.name == "" {
			if tt.wantErr {
				tt.name = "Error"
			} else {
				tt.name = string(tt.input.Value.Type)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
				v, err := api.convertToBotCommandScopeClass(ctx, tt.input)
				if tt.wantErr {
					a.Error(err)
					return
				}
				a.NoError(err)
				a.Equal(tt.want, v)
			})
		})
	}
}

func TestBotAPI_GetMyCommands(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.BotsGetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeDefault{},
			LangCode: "ru",
		}).ThenResult(&tg.BotCommandVector{Elems: testCommands()})
		commands, err := api.GetMyCommands(ctx, oas.NewOptGetMyCommands(oas.GetMyCommands{
			Scope:        oas.OptBotCommandScope{},
			LanguageCode: oas.NewOptString("ru"),
		}))
		a.NoError(err)
		a.Equal(testCommandsBotAPI(), commands.Result)

		mock.ExpectCall(&tg.BotsGetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeUsers{},
			LangCode: "ru",
		}).ThenResult(&tg.BotCommandVector{Elems: testCommands()})
		commands, err = api.GetMyCommands(ctx, oas.NewOptGetMyCommands(oas.GetMyCommands{
			Scope: oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeAllPrivateChatsBotCommandScope,
			}),
			LanguageCode: oas.NewOptString("ru"),
		}))
		a.NoError(err)
		a.Equal(testCommandsBotAPI(), commands.Result)
	})
}

func TestBotAPI_SetMyCommands(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.BotsSetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeDefault{},
			LangCode: "ru",
			Commands: testCommands(),
		}).ThenTrue()
		_, err := api.SetMyCommands(ctx, &oas.SetMyCommands{
			Scope:        oas.OptBotCommandScope{},
			LanguageCode: oas.NewOptString("ru"),
			Commands:     testCommandsBotAPI(),
		})
		a.NoError(err)

		mock.ExpectCall(&tg.BotsSetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeUsers{},
			LangCode: "ru",
			Commands: testCommands(),
		}).ThenTrue()
		_, err = api.SetMyCommands(ctx, &oas.SetMyCommands{
			Scope: oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeAllPrivateChatsBotCommandScope,
			}),
			LanguageCode: oas.NewOptString("ru"),
			Commands:     testCommandsBotAPI(),
		})
		a.NoError(err)
	})
}

func P[V any](v V) *V {
	return &v
}

func TestBotAPI_DeleteMyCommands(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.BotsResetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeDefault{},
			LangCode: "ru",
		}).ThenTrue()
		_, err := api.DeleteMyCommands(ctx, oas.NewOptDeleteMyCommands(oas.DeleteMyCommands{
			Scope:        oas.OptBotCommandScope{},
			LanguageCode: oas.NewOptString("ru"),
		}))
		a.NoError(err)

		mock.ExpectCall(&tg.BotsResetBotCommandsRequest{
			Scope:    &tg.BotCommandScopeUsers{},
			LangCode: "ru",
		}).ThenTrue()
		_, err = api.DeleteMyCommands(ctx, oas.NewOptDeleteMyCommands(oas.DeleteMyCommands{
			Scope: oas.NewOptBotCommandScope(oas.BotCommandScope{
				Type: oas.BotCommandScopeAllPrivateChatsBotCommandScope,
			}),
			LanguageCode: oas.NewOptString("ru"),
		}))
		a.NoError(err)
	})
}
