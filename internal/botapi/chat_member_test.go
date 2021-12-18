package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_GetChatMemberCount(t *testing.T) {
	ctx := context.Background()
	testWithChat(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		r, err := api.GetChatMemberCount(ctx, oas.GetChatMemberCount{
			ChatID: oas.NewInt64ID(chatID()),
		})
		a.NoError(err)
		a.Equal(oas.ResultInt{
			Result: oas.NewOptInt(10),
			Ok:     true,
		}, r)
	})
}
