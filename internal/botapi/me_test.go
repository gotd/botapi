package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_GetMe(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	mock, api := testBotAPI(t)

	user := testUser()

	mock.ExpectCall(&tg.UsersGetUsersRequest{
		ID: []tg.InputUserClass{&tg.InputUserSelf{}},
	}).ThenResult(&tg.UserClassVector{Elems: []tg.UserClass{user}})
	result, err := api.GetMe(ctx)
	a.NoError(err)
	a.Equal(&oas.ResultUser{
		Result: oas.OptUser{
			Value: oas.User{
				ID:                      user.ID,
				IsBot:                   user.Bot,
				FirstName:               user.FirstName,
				LastName:                oas.NewOptString(user.LastName),
				Username:                oas.NewOptString(user.Username),
				LanguageCode:            oas.OptString{},
				CanJoinGroups:           oas.NewOptBool(true),
				CanReadAllGroupMessages: oas.NewOptBool(false),
				SupportsInlineQueries:   oas.NewOptBool(user.BotInlinePlaceholder != ""),
			},
			Set: true,
		},
		Ok: true,
	}, result)
}

func TestBotAPI_LogOut(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	mock, api := testBotAPI(t)

	mock.ExpectCall(&tg.AuthLogOutRequest{}).ThenRPCErr(testError())
	_, err := api.LogOut(ctx)
	a.Error(err)

	mock.ExpectCall(&tg.AuthLogOutRequest{}).ThenResult(&tg.AuthLoggedOut{})
	_, err = api.LogOut(ctx)
	a.NoError(err)
}
