package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// SetMyProfilePhoto changes the bot's own profile photo. The photo must be an
// uploaded file (Telegram does not accept a file_id or URL for profile photos).
func (b *Bot) SetMyProfilePhoto(ctx context.Context, photo InputProfilePhoto) error {
	req, err := b.buildProfilePhotoUpload(ctx, photo)
	if err != nil {
		return err
	}

	if _, err := b.raw.PhotosUploadProfilePhoto(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// RemoveMyProfilePhoto removes the bot's own current profile photo.
func (b *Bot) RemoveMyProfilePhoto(ctx context.Context) error {
	if _, err := b.raw.PhotosUpdateProfilePhoto(ctx, &tg.PhotosUpdateProfilePhotoRequest{
		ID: &tg.InputPhotoEmpty{},
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
