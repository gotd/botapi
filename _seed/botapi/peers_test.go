package botapi

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_resolveID(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		var badRequestError *BadRequestError

		_, err := api.resolveID(ctx, oas.NewStringID(""))
		a.ErrorAs(err, &badRequestError)

		_, err = api.resolveID(ctx, oas.NewStringID("aboba"))
		a.ErrorAs(err, &badRequestError)

		p, err := api.resolveID(ctx, oas.NewStringID(strconv.FormatInt(testChatID(), 10)))
		a.NoError(err)
		a.IsType(peers.Chat{}, p)

		mock.ExpectCall(&tg.ContactsResolveUsernameRequest{
			Username: "tdakkota",
		}).ThenRPCErr(testError())
		_, err = api.resolveID(ctx, oas.NewStringID("@tdakkota"))
		a.Error(err)

		mock.ExpectCall(&tg.ContactsResolveUsernameRequest{
			Username: "tdakkota",
		}).ThenResult(&tg.ContactsResolvedPeer{
			Peer: &tg.PeerUser{UserID: 1337},
			Users: []tg.UserClass{
				&tg.User{
					ID:         1337,
					AccessHash: 1337,
					FirstName:  "tdakkota",
					Username:   "tdakkota",
				},
			},
		})
		p, err = api.resolveID(ctx, oas.NewStringID("@tdakkota"))
		a.NoError(err)
		a.IsType(peers.User{}, p)
	})
}
