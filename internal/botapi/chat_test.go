package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"
	"github.com/stretchr/testify/require"

	"github.com/gotd/botapi/internal/oas"
)

func testWithCache(t *testing.T, cb func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI)) {
	a := require.New(t)

	mock, api := testBotAPI(t)
	a.NoError(api.peers.Apply(context.Background(),
		[]tg.UserClass{
			testUser(),
		},
		[]tg.ChatClass{
			testChat(),
			testChannel(),
		},
	))

	cb(a, mock, api)
}

func Test_convertToBotAPIChatPermissions(t *testing.T) {
	tests := []struct {
		name string
		p    tg.ChatBannedRights
		want oas.ChatPermissions
	}{
		{
			name: "Zero",
			want: oas.ChatPermissions{
				CanSendMessages:       oas.NewOptBool(false),
				CanSendMediaMessages:  oas.NewOptBool(false),
				CanSendPolls:          oas.NewOptBool(false),
				CanSendOtherMessages:  oas.NewOptBool(false),
				CanAddWebPagePreviews: oas.NewOptBool(false),
				CanChangeInfo:         oas.NewOptBool(false),
				CanInviteUsers:        oas.NewOptBool(false),
				CanPinMessages:        oas.NewOptBool(false),
			},
		},
		{
			name: "Full",
			p: tg.ChatBannedRights{
				ViewMessages: true,
				SendMessages: true,
				SendMedia:    true,
				SendStickers: true,
				SendGifs:     true,
				SendGames:    true,
				SendInline:   true,
				EmbedLinks:   true,
				SendPolls:    true,
				ChangeInfo:   true,
				InviteUsers:  true,
				PinMessages:  true,
			},
			want: oas.ChatPermissions{
				CanSendMessages:       oas.NewOptBool(true),
				CanSendMediaMessages:  oas.NewOptBool(true),
				CanSendPolls:          oas.NewOptBool(true),
				CanSendOtherMessages:  oas.NewOptBool(true),
				CanAddWebPagePreviews: oas.NewOptBool(true),
				CanChangeInfo:         oas.NewOptBool(true),
				CanInviteUsers:        oas.NewOptBool(true),
				CanPinMessages:        oas.NewOptBool(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetFlags()
			require.Equal(t, tt.want, convertToBotAPIChatPermissions(tt.p))
		})
	}
}

func TestBotAPI_SetChatDescription(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesEditChatAboutRequest{
			Peer:  &tg.InputPeerChat{ChatID: 10},
			About: "",
		}).ThenTrue()
		_, err := api.SetChatDescription(ctx, oas.SetChatDescription{
			ChatID:      oas.NewInt64ID(testChatID()),
			Description: oas.OptString{},
		})
		a.NoError(err)

		mock.ExpectCall(&tg.MessagesEditChatAboutRequest{
			Peer:  &tg.InputPeerChat{ChatID: 10},
			About: "aboba",
		}).ThenTrue()
		_, err = api.SetChatDescription(ctx, oas.SetChatDescription{
			ChatID:      oas.NewInt64ID(testChatID()),
			Description: oas.NewOptString("aboba"),
		})
		a.NoError(err)
	})
}

func TestBotAPI_ApproveChatJoinRequest(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesHideChatJoinRequestRequest{
			Approved: true,
			Peer:     &tg.InputPeerChat{ChatID: testChat().ID},
			UserID:   &tg.InputUserSelf{},
		}).ThenResult(&tg.Updates{})
		_, err := api.ApproveChatJoinRequest(ctx, oas.ApproveChatJoinRequest{
			ChatID: oas.NewInt64ID(testChatID()),
			UserID: testUser().ID,
		})
		a.NoError(err)
	})
}

func TestBotAPI_DeclineChatJoinRequest(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesHideChatJoinRequestRequest{
			Approved: false,
			Peer:     &tg.InputPeerChat{ChatID: testChat().ID},
			UserID:   &tg.InputUserSelf{},
		}).ThenResult(&tg.Updates{})
		_, err := api.DeclineChatJoinRequest(ctx, oas.DeclineChatJoinRequest{
			ChatID: oas.NewInt64ID(testChatID()),
			UserID: testUser().ID,
		})
		a.NoError(err)
	})
}

func TestBotAPI_DeleteChatStickerSet(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		_, err := api.DeleteChatStickerSet(ctx, oas.DeleteChatStickerSet{
			ChatID: oas.NewInt64ID(testChatID()),
		})
		a.Error(err)

		mock.ExpectCall(&tg.ChannelsSetStickersRequest{
			Channel: &tg.InputChannel{
				ChannelID:  testChannel().ID,
				AccessHash: testChannel().AccessHash,
			},
			Stickerset: &tg.InputStickerSetEmpty{},
		}).ThenTrue()
		_, err = api.DeleteChatStickerSet(ctx, oas.DeleteChatStickerSet{
			ChatID: oas.NewInt64ID(testChannelID()),
		})
		a.NoError(err)
	})
}

func TestBotAPI_GetChat(t *testing.T) {

}

func TestBotAPI_SetChatTitle(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesEditChatTitleRequest{
			ChatID: testChat().ID,
			Title:  "aboba",
		}).ThenResult(&tg.Updates{})
		_, err := api.SetChatTitle(ctx, oas.SetChatTitle{
			ChatID: oas.NewInt64ID(testChatID()),
			Title:  "aboba",
		})
		a.NoError(err)
	})
}

func TestBotAPI_LeaveChat(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.MessagesDeleteChatUserRequest{
			ChatID: testChat().ID,
			UserID: &tg.InputUserSelf{},
		}).ThenResult(&tg.Updates{})
		_, err := api.LeaveChat(ctx, oas.LeaveChat{
			ChatID: oas.NewInt64ID(testChatID()),
		})
		a.NoError(err)
	})
}

func TestBotAPI_DeleteChatPhoto(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.ChannelsEditPhotoRequest{
			Channel: &tg.InputChannel{
				ChannelID:  testChannel().ID,
				AccessHash: testChannel().AccessHash,
			},
			Photo: &tg.InputChatPhotoEmpty{},
		}).ThenTrue()
		_, err := api.DeleteChatPhoto(ctx, oas.DeleteChatPhoto{
			ChatID: oas.NewInt64ID(testChannelID()),
		})
		a.NoError(err)
	})
}
