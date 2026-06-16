package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// ApproveSuggestedPostOption configures ApproveSuggestedPost.
type ApproveSuggestedPostOption func(*approveSuggestedPostConfig)

type approveSuggestedPostConfig struct {
	sendDate    int
	sendDateSet bool
}

// WithSuggestedPostSendDate schedules the approved post to be published at the
// given Unix time instead of immediately.
func WithSuggestedPostSendDate(date int) ApproveSuggestedPostOption {
	return func(c *approveSuggestedPostConfig) {
		c.sendDate = date
		c.sendDateSet = true
	}
}

// ApproveSuggestedPost approves a suggested post in a direct messages chat. The
// bot must have the can_post_messages administrator right in the corresponding
// channel chat.
func (b *Bot) ApproveSuggestedPost(ctx context.Context, chat ChatID, messageID int, opts ...ApproveSuggestedPostOption) error {
	var cfg approveSuggestedPostConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	req := &tg.MessagesToggleSuggestedPostApprovalRequest{
		Peer:  peer,
		MsgID: messageID,
	}
	if cfg.sendDateSet {
		req.SetScheduleDate(cfg.sendDate)
	}

	if _, err := b.raw.MessagesToggleSuggestedPostApproval(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// DeclineSuggestedPost declines a suggested post in a direct messages chat. The
// bot must have the can_manage_direct_messages administrator right in the
// corresponding channel chat. An optional comment explains the rejection.
func (b *Bot) DeclineSuggestedPost(ctx context.Context, chat ChatID, messageID int, comment string) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	req := &tg.MessagesToggleSuggestedPostApprovalRequest{
		Peer:   peer,
		MsgID:  messageID,
		Reject: true,
	}
	if comment != "" {
		req.SetRejectComment(comment)
	}

	if _, err := b.raw.MessagesToggleSuggestedPostApproval(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}
