package botapi

import (
	"context"
	"crypto/rand"
	"encoding/binary"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/constant"
	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// storyConfig holds the optional parameters of the story-posting methods.
type storyConfig struct {
	caption         string
	parseMode       ParseMode
	captionEntities []MessageEntity
	postToChatPage  bool
	protectContent  bool
}

// StoryOption customizes PostStory, EditStory and RepostStory.
type StoryOption func(*storyConfig)

// WithStoryCaption sets the story caption.
func WithStoryCaption(caption string) StoryOption {
	return func(c *storyConfig) { c.caption = caption }
}

// WithStoryParseMode sets the parse mode used for the caption.
func WithStoryParseMode(mode ParseMode) StoryOption {
	return func(c *storyConfig) { c.parseMode = mode }
}

// WithStoryCaptionEntities sets explicit caption entities, overriding the parse
// mode.
func WithStoryCaptionEntities(entities []MessageEntity) StoryOption {
	return func(c *storyConfig) { c.captionEntities = entities }
}

// WithStoryPostToChatPage also posts the story to the chat's profile page.
func WithStoryPostToChatPage() StoryOption {
	return func(c *storyConfig) { c.postToChatPage = true }
}

// WithStoryProtectContent protects the story content from forwarding and
// screenshots.
func WithStoryProtectContent() StoryOption {
	return func(c *storyConfig) { c.protectContent = true }
}

// caption resolves the caption text into a (text, entities) pair: explicit
// entities take precedence over the parse mode.
func (b *Bot) storyCaption(ctx context.Context, cfg storyConfig) (string, []tg.MessageEntityClass, error) {
	if cfg.caption == "" {
		return "", nil, nil
	}

	if len(cfg.captionEntities) > 0 {
		return cfg.caption, entitiesToTg(cfg.captionEntities), nil
	}

	return b.styledMessage(ctx, cfg.caption, cfg.parseMode)
}

// PostStory posts a story on behalf of a managed business account, active for
// activePeriod seconds (one of 6, 12, 24 or 48 hours). Requires the
// can_manage_stories business bot right. Story areas are not yet supported.
func (b *Bot) PostStory(
	ctx context.Context, businessConnectionID string, content InputStoryContent, activePeriod int, opts ...StoryOption,
) (*Story, error) {
	cfg := storyConfigFrom(opts)

	media, err := b.storyMedia(ctx, content)
	if err != nil {
		return nil, err
	}

	caption, entities, err := b.storyCaption(ctx, cfg)
	if err != nil {
		return nil, err
	}

	id, err := randInt64()
	if err != nil {
		return nil, err
	}

	req := &tg.StoriesSendStoryRequest{
		Peer:         &tg.InputPeerSelf{},
		Media:        media,
		Caption:      caption,
		Entities:     entities,
		PrivacyRules: []tg.InputPrivacyRuleClass{&tg.InputPrivacyValueAllowAll{}},
		RandomID:     id,
		Period:       activePeriod,
		Pinned:       cfg.postToChatPage,
		Noforwards:   cfg.protectContent,
	}

	return b.sendStory(ctx, businessConnectionID, req)
}

// RepostStory reposts a story from another business account managed by the same
// bot. Requires the can_manage_stories business bot right for both accounts.
func (b *Bot) RepostStory(
	ctx context.Context, businessConnectionID string, fromChat ChatID, fromStoryID, activePeriod int, opts ...StoryOption,
) (*Story, error) {
	cfg := storyConfigFrom(opts)

	from, err := b.resolveInputPeer(ctx, fromChat)
	if err != nil {
		return nil, err
	}

	id, err := randInt64()
	if err != nil {
		return nil, err
	}

	req := &tg.StoriesSendStoryRequest{
		Peer: &tg.InputPeerSelf{},
		// A repost carries no new media; the content comes from the source story.
		Media:        &tg.InputMediaEmpty{},
		PrivacyRules: []tg.InputPrivacyRuleClass{&tg.InputPrivacyValueAllowAll{}},
		RandomID:     id,
		Period:       activePeriod,
		Pinned:       cfg.postToChatPage,
		Noforwards:   cfg.protectContent,
		FwdFromID:    from,
		FwdFromStory: fromStoryID,
	}

	return b.sendStory(ctx, businessConnectionID, req)
}

// EditStory edits a story previously posted by the bot on behalf of a managed
// business account. Requires the can_manage_stories business bot right. Story
// areas are not yet supported.
func (b *Bot) EditStory(
	ctx context.Context, businessConnectionID string, storyID int, content InputStoryContent, opts ...StoryOption,
) (*Story, error) {
	cfg := storyConfigFrom(opts)

	media, err := b.storyMedia(ctx, content)
	if err != nil {
		return nil, err
	}

	caption, entities, err := b.storyCaption(ctx, cfg)
	if err != nil {
		return nil, err
	}

	req := &tg.StoriesEditStoryRequest{
		Peer:     &tg.InputPeerSelf{},
		ID:       storyID,
		Media:    media,
		Caption:  caption,
		Entities: entities,
	}

	return b.sendStory(ctx, businessConnectionID, req)
}

// DeleteStory deletes a story previously posted by the bot on behalf of a
// managed business account. Requires the can_manage_stories business bot right.
func (b *Bot) DeleteStory(ctx context.Context, businessConnectionID string, storyID int) error {
	var res tg.IntVector

	err := b.invokeBusiness(ctx, businessConnectionID, &tg.StoriesDeleteStoriesRequest{
		Peer: &tg.InputPeerSelf{},
		ID:   []int{storyID},
	}, &res)
	if err != nil {
		return asAPIError(err)
	}

	return nil
}

// sendStory runs a story-mutating request over the business connection and
// extracts the resulting Story from the update list.
func (b *Bot) sendStory(ctx context.Context, connectionID string, req bin.Object) (*Story, error) {
	var res tg.UpdatesBox

	if err := b.invokeBusiness(ctx, connectionID, req, &res); err != nil {
		return nil, asAPIError(err)
	}

	story := storyFromUpdates(res.Updates)
	if story == nil {
		return nil, &Error{Code: 500, Description: "Internal Server Error: no story in response"}
	}

	return story, nil
}

// storyConfigFrom collapses the options into a config.
func storyConfigFrom(opts []StoryOption) storyConfig {
	var cfg storyConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// storyMedia builds the MTProto input media for a story content.
//
// The switch over the sealed InputStoryContent union is exhaustive.
func (b *Bot) storyMedia(ctx context.Context, content InputStoryContent) (tg.InputMediaClass, error) {
	switch c := content.(type) {
	case InputStoryContentPhoto:
		return b.storyPhotoMedia(ctx, c.Photo)
	case InputStoryContentVideo:
		return b.storyVideoMedia(ctx, c)
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// storyPhotoMedia builds the input media for a photo story.
func (b *Bot) storyPhotoMedia(ctx context.Context, photo InputFile) (tg.InputMediaClass, error) {
	switch f := photo.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		return &tg.InputMediaPhoto{ID: &tg.InputPhoto{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}, nil
	case InputFileURL:
		return &tg.InputMediaPhotoExternal{URL: string(f)}, nil
	case *InputFileUpload:
		upFile, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		return &tg.InputMediaUploadedPhoto{File: upFile}, nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// storyVideoMedia builds the input media for a video story.
func (b *Bot) storyVideoMedia(ctx context.Context, content InputStoryContentVideo) (tg.InputMediaClass, error) {
	switch f := content.Video.(type) {
	case InputFileID:
		fid, err := fileid.DecodeFileID(string(f))
		if err != nil {
			return nil, &Error{Code: 400, Description: descWrongFileID}
		}

		return &tg.InputMediaDocument{ID: &tg.InputDocument{
			ID:            fid.ID,
			AccessHash:    fid.AccessHash,
			FileReference: fid.FileReference,
		}}, nil
	case InputFileURL:
		return &tg.InputMediaDocumentExternal{URL: string(f)}, nil
	case *InputFileUpload:
		upFile, err := b.uploadInputFile(ctx, f)
		if err != nil {
			return nil, err
		}

		video := tg.DocumentAttributeVideo{
			SupportsStreaming: true,
			Duration:          content.Duration,
			Nosound:           content.IsAnimation,
		}

		return &tg.InputMediaUploadedDocument{
			File:         upFile,
			MimeType:     mimeVideoMP4,
			NosoundVideo: content.IsAnimation,
			Attributes:   []tg.DocumentAttributeClass{&video},
		}, nil
	default:
		return nil, &Error{Code: 400, Description: descInvalidFile}
	}
}

// storyFromUpdates extracts the posted/edited story from an update list.
func storyFromUpdates(resp tg.UpdatesClass) *Story {
	var (
		updates []tg.UpdateClass
		users   map[int64]*tg.User
		chats   map[int64]tg.ChatClass
	)

	switch u := resp.(type) {
	case *tg.Updates:
		updates, users, chats = u.Updates, usersByID(u.Users), chatsByID(u.Chats)
	case *tg.UpdatesCombined:
		updates, users, chats = u.Updates, usersByID(u.Users), chatsByID(u.Chats)
	default:
		return nil
	}

	for _, upd := range updates {
		us, ok := upd.(*tg.UpdateStory)
		if !ok {
			continue
		}

		item, ok := us.Story.(*tg.StoryItem)
		if !ok {
			continue
		}

		return &Story{Chat: storyChat(us.Peer, users, chats), ID: item.ID}
	}

	return nil
}

// storyChat resolves the peer that owns a story into a Bot API Chat.
func storyChat(peer tg.PeerClass, users map[int64]*tg.User, chats map[int64]tg.ChatClass) Chat {
	pu, ok := peer.(*tg.PeerUser)
	if !ok {
		return chatFromRaw(peer, chats)
	}

	var id constant.TDLibPeerID

	id.User(pu.UserID)

	c := Chat{ID: int64(id), Type: ChatTypePrivate}
	if u, ok := users[pu.UserID]; ok {
		c.FirstName = u.FirstName
		c.LastName = u.LastName
		c.Username = u.Username
	}

	return c
}

// randInt64 returns a cryptographically random int64 for use as an RPC random
// id.
func randInt64() (int64, error) {
	var buf [8]byte

	if _, err := rand.Read(buf[:]); err != nil {
		return 0, &Error{Code: 500, Description: "Internal Server Error: " + err.Error()}
	}

	return int64(binary.LittleEndian.Uint64(buf[:])), nil //nolint:gosec // wraparound is fine for a random id
}
