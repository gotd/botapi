package botapi

import (
	"context"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// UserProfileAudios is the result of GetUserProfileAudios: the audios a user has
// added to their profile.
type UserProfileAudios struct {
	// TotalCount is the total number of profile audios the user has.
	TotalCount int `json:"total_count"`
	// Audios are the requested profile audios.
	Audios []Audio `json:"audios"`
}

// GetUserProfileAudiosOption configures GetUserProfileAudios.
type GetUserProfileAudiosOption func(*profileAudiosConfig)

type profileAudiosConfig struct {
	offset int
	limit  int
}

// WithProfileAudiosOffset sets the number of audios to skip.
func WithProfileAudiosOffset(offset int) GetUserProfileAudiosOption {
	return func(c *profileAudiosConfig) { c.offset = offset }
}

// WithProfileAudiosLimit caps the number of audios returned (1-100).
func WithProfileAudiosLimit(limit int) GetUserProfileAudiosOption {
	return func(c *profileAudiosConfig) { c.limit = limit }
}

// audioFromDocument builds a Bot API Audio from an MTProto audio document,
// reusing the shared document classifier.
func audioFromDocument(d *tg.Document) Audio {
	var m Message

	setDocumentMedia(d, &m)

	if m.Audio != nil {
		return *m.Audio
	}

	fileID, uniqueID := encodeFileID(fileid.FromDocument(d))

	return Audio{FileID: fileID, FileUniqueID: uniqueID, MIMEType: d.MimeType, FileSize: d.Size}
}

// GetUserProfileAudios returns the audios a user has added to their profile.
func (b *Bot) GetUserProfileAudios(ctx context.Context, userID int64, opts ...GetUserProfileAudiosOption) (*UserProfileAudios, error) {
	cfg := profileAudiosConfig{limit: 100}
	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.UsersGetSavedMusic(ctx, &tg.UsersGetSavedMusicRequest{
		ID:     user,
		Offset: cfg.offset,
		Limit:  cfg.limit,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	music, ok := res.(*tg.UsersSavedMusic)
	if !ok {
		return &UserProfileAudios{}, nil
	}

	out := &UserProfileAudios{TotalCount: music.Count, Audios: make([]Audio, 0, len(music.Documents))}

	for _, doc := range music.Documents {
		if d, ok := doc.(*tg.Document); ok {
			out.Audios = append(out.Audios, audioFromDocument(d))
		}
	}

	return out, nil
}
