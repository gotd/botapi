package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func testWithChat(t *testing.T, cb func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI)) {
	a := require.New(t)

	mock, api := testBotAPI(t)
	a.NoError(api.peers.Apply(context.Background(), nil, []tg.ChatClass{
		&tg.Chat{ID: 10, ParticipantsCount: 10},
	}))

	cb(a, mock, api)
}

func TestBotAPI_SetChatDescription(t *testing.T) {
	ctx := context.Background()
	testWithChat(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesEditChatAboutRequest{
			Peer:  &tg.InputPeerChat{ChatID: 10},
			About: "",
		}).ThenTrue()
		_, err := api.SetChatDescription(ctx, oas.SetChatDescription{
			ChatID:      oas.NewInt64ID(chatID()),
			Description: oas.OptString{},
		})
		a.NoError(err)

		mock.ExpectCall(&tg.MessagesEditChatAboutRequest{
			Peer:  &tg.InputPeerChat{ChatID: 10},
			About: "aboba",
		}).ThenTrue()
		_, err = api.SetChatDescription(ctx, oas.SetChatDescription{
			ChatID:      oas.NewInt64ID(chatID()),
			Description: oas.NewOptString("aboba"),
		})
		a.NoError(err)
	})
}

func TestBotAPI_SetChatTitle(t *testing.T) {
	ctx := context.Background()
	testWithChat(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesEditChatTitleRequest{
			ChatID: 10,
			Title:  "aboba",
		}).ThenResult(&tg.Updates{})
		_, err := api.SetChatTitle(ctx, oas.SetChatTitle{
			ChatID: oas.NewInt64ID(chatID()),
			Title:  "aboba",
		})
		a.NoError(err)
	})
}

func TestBotAPI_LeaveChat(t *testing.T) {
	ctx := context.Background()
	testWithChat(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesDeleteChatUserRequest{
			ChatID: 10,
			UserID: &tg.InputUserSelf{},
		}).ThenResult(&tg.Updates{})
		_, err := api.LeaveChat(ctx, oas.LeaveChat{
			ChatID: oas.NewInt64ID(chatID()),
		})
		a.NoError(err)
	})
}
