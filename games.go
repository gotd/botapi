package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

// Game represents a game. Use BotFather to set up a game for your bot.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation      `json:"animation,omitempty"`
}

// GameHighScore represents one row of a game's high-score table.
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}

// SendGame sends a game identified by its short name (configured via BotFather).
func (b *Bot) SendGame(ctx context.Context, chat ChatID, gameShortName string, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	media := message.Media(&tg.InputMediaGame{
		ID: &tg.InputGameShortName{BotID: &tg.InputUserSelf{}, ShortName: gameShortName},
	})

	builder := &b.sender.To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.Media(ctx, media)

	return b.sentMessage(ctx, peer, resp, err)
}

// SetGameScoreOption configures a SetGameScore call.
type SetGameScoreOption func(*gameScoreConfig)

type gameScoreConfig struct {
	force              bool
	disableEditMessage bool
}

// WithForceScore updates the score even if it is lower than the user's current
// best.
func WithForceScore() SetGameScoreOption {
	return func(c *gameScoreConfig) { c.force = true }
}

// WithoutEditMessage leaves the game message unchanged instead of updating it
// with the new score.
func WithoutEditMessage() SetGameScoreOption {
	return func(c *gameScoreConfig) { c.disableEditMessage = true }
}

// SetGameScore sets a user's score in the game contained in the given message.
func (b *Bot) SetGameScore(
	ctx context.Context, chat ChatID, messageID int, userID int64, score int, opts ...SetGameScoreOption,
) (*Message, error) {
	var cfg gameScoreConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp, err := b.raw.MessagesSetGameScore(ctx, &tg.MessagesSetGameScoreRequest{
		EditMessage: !cfg.disableEditMessage,
		Force:       cfg.force,
		Peer:        peer,
		ID:          messageID,
		UserID:      user,
		Score:       score,
	})

	return b.sentMessage(ctx, peer, resp, err)
}

// GetGameHighScores returns the high scores of the game in the given message for
// the user and their close neighbors.
func (b *Bot) GetGameHighScores(ctx context.Context, chat ChatID, messageID int, userID int64) ([]GameHighScore, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.MessagesGetGameHighScores(ctx, &tg.MessagesGetGameHighScoresRequest{
		Peer:   peer,
		ID:     messageID,
		UserID: user,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	users := usersByID(res.Users)
	out := make([]GameHighScore, 0, len(res.Scores))

	for _, s := range res.Scores {
		hs := GameHighScore{Position: s.Pos, Score: s.Score, User: User{ID: s.UserID}}
		if u, ok := users[s.UserID]; ok {
			hs.User = userFromTgUser(u)
		}

		out = append(out, hs)
	}

	return out, nil
}
