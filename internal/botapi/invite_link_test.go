package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_RevokeChatInviteLink(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesEditExportedChatInviteRequest{
			Revoked: true,
			Peer:    &tg.InputPeerChat{ChatID: testChat().ID},
			Link:    "aboba",
		}).ThenResult(&tg.MessagesExportedChatInvite{
			Invite: tg.ChatInviteExported{
				Revoked: true,
				Link:    "aboba",
				AdminID: testUser().ID,
			},
			Users: []tg.UserClass{testUser()},
		})
		_, err := api.RevokeChatInviteLink(ctx, oas.RevokeChatInviteLink{
			ChatID:     oas.NewInt64ID(testChatID()),
			InviteLink: "aboba",
		})
		a.NoError(err)
	})
}
