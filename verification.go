package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// VerificationOption configures a verify call.
type VerificationOption func(*verificationConfig)

type verificationConfig struct {
	description string
}

// WithVerificationDescription sets a custom description for the verification
// shown instead of the bot's default verification description.
func WithVerificationDescription(text string) VerificationOption {
	return func(c *verificationConfig) { c.description = text }
}

// setCustomVerification enables or disables the bot's custom verification of a
// peer.
func (b *Bot) setCustomVerification(ctx context.Context, peer tg.InputPeerClass, enabled bool, opts []VerificationOption) error {
	var cfg verificationConfig

	for _, o := range opts {
		o(&cfg)
	}

	req := &tg.BotsSetCustomVerificationRequest{Peer: peer}
	req.SetEnabled(enabled)

	if enabled && cfg.description != "" {
		req.SetCustomDescription(cfg.description)
	}

	if _, err := b.raw.BotsSetCustomVerification(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// VerifyUser verifies a user on behalf of the organization represented by the
// bot.
func (b *Bot) VerifyUser(ctx context.Context, userID int64, opts ...VerificationOption) error {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	return b.setCustomVerification(ctx, userToInputPeer(user), true, opts)
}

// VerifyChat verifies a chat on behalf of the organization represented by the
// bot.
func (b *Bot) VerifyChat(ctx context.Context, chat ChatID, opts ...VerificationOption) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	return b.setCustomVerification(ctx, peer, true, opts)
}

// RemoveUserVerification removes the verification from a user that is verified on
// behalf of the organization represented by the bot.
func (b *Bot) RemoveUserVerification(ctx context.Context, userID int64) error {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	return b.setCustomVerification(ctx, userToInputPeer(user), false, nil)
}

// RemoveChatVerification removes the verification from a chat that is verified on
// behalf of the organization represented by the bot.
func (b *Bot) RemoveChatVerification(ctx context.Context, chat ChatID) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	return b.setCustomVerification(ctx, peer, false, nil)
}
