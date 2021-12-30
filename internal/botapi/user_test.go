package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_GetUserProfilePhotos(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		mock.ExpectCall(&tg.PhotosGetUserPhotosRequest{
			UserID: &tg.InputUserSelf{},
			Offset: 1,
			Limit:  100,
		}).ThenResult(&tg.PhotosPhotosSlice{
			Count: 10,
			Photos: []tg.PhotoClass{
				&tg.Photo{
					ID:         10,
					AccessHash: 10,
					Date:       10,
					Sizes: []tg.PhotoSizeClass{
						&tg.PhotoSize{
							Type: "m",
							W:    10,
							H:    10,
							Size: 10,
						},
					},
					DCID: 2,
				},
				&tg.Photo{
					ID:         10,
					AccessHash: 10,
					Date:       10,
					Sizes: []tg.PhotoSizeClass{
						&tg.PhotoCachedSize{
							Type:  "m",
							W:     10,
							H:     10,
							Bytes: []byte("data"),
						},
					},
					DCID: 2,
				},
			},
			Users: nil,
		})
		r, err := api.GetUserProfilePhotos(ctx, oas.GetUserProfilePhotos{
			UserID: testUser().ID,
			Offset: oas.NewOptInt(1),
		})
		a.NoError(err)
		a.True(r.Result.Set)
		val := r.Result.Value

		a.Equal(10, val.TotalCount)
		a.Len(val.Photos, 2)
		a.Len(val.Photos[0], 1)
		a.Len(val.Photos[1], 1)

		f, s := val.Photos[0][0], val.Photos[1][0]
		a.NotEmpty(f.FileID)
		a.NotEmpty(s.FileID)
		a.Equal(f.FileSize.Value, 10)
		a.Equal(s.FileSize.Value, 4)
	})
}
