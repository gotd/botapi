package botapi

import "context"

// CopyMessage copies a message: it re-sends the message's content to another
// chat as a new message, without the "forwarded from" header. Unlike
// ForwardMessage, the result is not linked to the original.
func (b *Bot) CopyMessage(ctx context.Context, to, from ChatID, messageID int, opts ...SendOption) (*Message, error) {
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

	resp, err := builder.ForwardIDs(fromPeer, messageID).DropAuthor().Send(ctx)

	return b.sentMessage(ctx, toPeer, resp, err)
}
