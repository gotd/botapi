package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// GetChat returns up-to-date information about a chat.
func (b *Bot) GetChat(ctx context.Context, chat ChatID) (*Chat, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	c := chatFromPeer(p)

	return &c, nil
}

// GetUserProfilePhotosOption configures a GetUserProfilePhotos call.
type GetUserProfilePhotosOption func(*userPhotosConfig)

type userPhotosConfig struct {
	offset int
	limit  int
}

// WithProfilePhotosOffset sets the number of photos to skip.
func WithProfilePhotosOffset(offset int) GetUserProfilePhotosOption {
	return func(c *userPhotosConfig) { c.offset = offset }
}

// WithProfilePhotosLimit caps the number of photos returned (1-100).
func WithProfilePhotosLimit(limit int) GetUserProfilePhotosOption {
	return func(c *userPhotosConfig) { c.limit = limit }
}

// GetUserProfilePhotos returns a user's profile photos.
func (b *Bot) GetUserProfilePhotos(ctx context.Context, userID int64, opts ...GetUserProfilePhotosOption) (*UserProfilePhotos, error) {
	cfg := userPhotosConfig{limit: 100}
	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.PhotosGetUserPhotos(ctx, &tg.PhotosGetUserPhotosRequest{
		UserID: user,
		Offset: cfg.offset,
		Limit:  cfg.limit,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	var (
		photos []tg.PhotoClass
		total  int
	)

	switch p := res.(type) {
	case *tg.PhotosPhotos:
		photos = p.Photos
		total = len(p.Photos)
	case *tg.PhotosPhotosSlice:
		photos = p.Photos
		total = p.Count
	}

	out := &UserProfilePhotos{TotalCount: total}

	for _, ph := range photos {
		if photo, ok := ph.(*tg.Photo); ok {
			out.Photos = append(out.Photos, photoSizesFromTg(photo))
		}
	}

	return out, nil
}
