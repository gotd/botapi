package botapi

import "context"

// ForwardMessage forwards a single message from one chat to another and returns
// the forwarded message. Silent and ProtectContent options apply.
func (b *Bot) ForwardMessage(ctx context.Context, to, from ChatID, messageID int, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	toPeer, err := b.resolveInputPeer(ctx, to)
	if err != nil {
		return nil, err
	}

	fromPeer, err := b.resolveInputPeer(ctx, from)
	if err != nil {
		return nil, err
	}

	builder := &b.sender.To(toPeer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.ForwardIDs(fromPeer, messageID).Send(ctx)

	return b.sentMessage(ctx, toPeer, resp, err)
}

// DeleteMessage deletes a message for everyone in the chat.
func (b *Bot) DeleteMessage(ctx context.Context, chat ChatID, messageID int) error {
	return b.DeleteMessages(ctx, chat, []int{messageID})
}

// DeleteMessages deletes several messages for everyone in the chat.
func (b *Bot) DeleteMessages(ctx context.Context, chat ChatID, messageIDs []int) error {
	if len(messageIDs) == 0 {
		return nil
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	_, err = b.sender.To(peer).Revoke().Messages(ctx, messageIDs...)

	return asAPIError(err)
}
