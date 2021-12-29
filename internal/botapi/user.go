package botapi

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

// GetUserProfilePhotos implements oas.Handler.
func (b *BotAPI) GetUserProfilePhotos(ctx context.Context, req oas.GetUserProfilePhotos) (oas.ResultUserProfilePhotos, error) {
	userID, err := b.resolveUserID(ctx, req.UserID)
	if err != nil {
		return oas.ResultUserProfilePhotos{}, errors.Wrap(err, "resolve userID")
	}

	response, err := b.raw.PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: userID.InputUser(),
		Offset: req.Offset.Value,
		Limit:  req.Limit.Or(100),
	})
	if err != nil {
		return oas.ResultUserProfilePhotos{}, errors.Wrap(err, "get photos")
	}
	var totalCount int
	switch response := response.(type) {
	case *tg.PhotosPhotos:
		totalCount = len(response.Photos)
	case *tg.PhotosPhotosSlice:
		totalCount = response.Count
	}

	var photos [][]oas.PhotoSize
	for _, p := range response.MapPhotos() {
		photos = append(photos, b.convertToBotAPIPhotoSizes(p))
	}

	return oas.ResultUserProfilePhotos{
		Result: oas.NewOptUserProfilePhotos(oas.UserProfilePhotos{
			TotalCount: totalCount,
			Photos:     photos,
		}),
		Ok: true,
	}, nil
}
