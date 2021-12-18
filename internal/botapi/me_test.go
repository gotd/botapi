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

	user := &tg.User{
		Self:                 true,
		Bot:                  true,
		ID:                   10,
		AccessHash:           10,
		FirstName:            "Elsa",
		LastName:             "Jean",
		Username:             "thebot",
		BotInfoVersion:       1,
		BotInlinePlaceholder: "aboba",
	}
	user.SetFlags()

	mock.ExpectCall(&tg.UsersGetUsersRequest{
		ID: []tg.InputUserClass{&tg.InputUserSelf{}},
	}).ThenResult(&tg.UserClassVector{Elems: []tg.UserClass{user}})
	result, err := api.GetMe(ctx)
	a.NoError(err)
	a.Equal(oas.ResultUser{
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
