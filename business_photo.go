package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// InputProfilePhoto is a sealed union describing a profile photo to set on a
// managed business account.
//
// Concrete variants: InputProfilePhotoStatic, InputProfilePhotoAnimated.
type InputProfilePhoto interface {
	isInputProfilePhoto()
}

// InputProfilePhotoStatic is a static profile photo.
type InputProfilePhotoStatic struct {
	// Photo is the static profile photo, an uploaded file (file_id and URL are
	// not accepted by Telegram for profile photos).
	Photo InputFile
}

// InputProfilePhotoAnimated is an animated profile photo (a short video).
type InputProfilePhotoAnimated struct {
	// Animation is the animated profile photo, an uploaded file.
	Animation InputFile
	// MainFrameTimestamp is the timestamp in seconds of the frame used as a
	// static profile photo. Defaults to the start of the animation.
	MainFrameTimestamp float64
}

func (InputProfilePhotoStatic) isInputProfilePhoto()   {}
func (InputProfilePhotoAnimated) isInputProfilePhoto() {}

// uploadProfileFile uploads a profile photo file. Profile photos must be raw
// uploads; a file_id or URL references an existing document and cannot back the
// InputFile the profile-photo RPCs need.
func (b *Bot) uploadProfileFile(ctx context.Context, f InputFile) (tg.InputFileClass, error) {
	up, ok := f.(*InputFileUpload)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: profile photo must be an uploaded file"}
	}

	return b.uploadInputFile(ctx, up)
}

// buildProfilePhotoUpload turns an InputProfilePhoto into a populated
// photos.uploadProfilePhoto request, uploading the underlying file.
//
// The switch over the sealed InputProfilePhoto union is exhaustive
// (gochecksumtype).
func (b *Bot) buildProfilePhotoUpload(ctx context.Context, photo InputProfilePhoto) (*tg.PhotosUploadProfilePhotoRequest, error) {
	req := &tg.PhotosUploadProfilePhotoRequest{}

	switch p := photo.(type) {
	case InputProfilePhotoStatic:
		file, err := b.uploadProfileFile(ctx, p.Photo)
		if err != nil {
			return nil, err
		}

		req.SetFile(file)
	case InputProfilePhotoAnimated:
		video, err := b.uploadProfileFile(ctx, p.Animation)
		if err != nil {
			return nil, err
		}

		req.SetVideo(video)

		if p.MainFrameTimestamp != 0 {
			req.SetVideoStartTs(p.MainFrameTimestamp)
		}
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: invalid profile photo"}
	}

	return req, nil
}

// SetBusinessAccountProfilePhoto changes the profile photo of a managed business
// account. When isPublic is true the photo is set as the public (fallback) photo,
// shown to users who cannot see the account's main photo. The bot must have the
// can_edit_profile_photo business bot right.
//
// The switch over the sealed InputProfilePhoto union is exhaustive
// (gochecksumtype).
func (b *Bot) SetBusinessAccountProfilePhoto(
	ctx context.Context, businessConnectionID string, photo InputProfilePhoto, isPublic bool,
) error {
	req, err := b.buildProfilePhotoUpload(ctx, photo)
	if err != nil {
		return err
	}

	if isPublic {
		req.SetFallback(true)
	}

	return b.invokeBusiness(ctx, businessConnectionID, req, &tg.PhotosPhoto{})
}

// RemoveBusinessAccountProfilePhoto removes the profile photo of a managed
// business account. When isPublic is true the public (fallback) photo is removed
// instead of the main one. The bot must have the can_edit_profile_photo business
// bot right.
func (b *Bot) RemoveBusinessAccountProfilePhoto(ctx context.Context, businessConnectionID string, isPublic bool) error {
	req := &tg.PhotosUpdateProfilePhotoRequest{ID: &tg.InputPhotoEmpty{}}

	if isPublic {
		req.SetFallback(true)
	}

	return b.invokeBusiness(ctx, businessConnectionID, req, &tg.PhotosPhoto{})
}
