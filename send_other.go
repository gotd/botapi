package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

// sendMedia applies the send options to a builder and sends a media message.
func (b *Bot) sendMedia(ctx context.Context, chat ChatID, media message.MediaOption, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	builder := &b.sender.To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.Media(ctx, media)

	return b.sentMessage(ctx, peer, resp, err)
}

// SendDice sends an animated emoji with a random value. A zero emoji defaults
// to the die (🎲).
func (b *Bot) SendDice(ctx context.Context, chat ChatID, emoji DiceEmoji, opts ...SendOption) (*Message, error) {
	if emoji == "" {
		emoji = DiceDie
	}

	return b.sendMedia(ctx, chat, message.MediaDice(string(emoji)), opts...)
}

// SendLocation sends a point on the map.
func (b *Bot) SendLocation(ctx context.Context, chat ChatID, latitude, longitude float64, opts ...SendOption) (*Message, error) {
	return b.sendMedia(ctx, chat, message.GeoPoint(latitude, longitude, 0), opts...)
}

// SendVenue sends information about a venue.
func (b *Bot) SendVenue(
	ctx context.Context, chat ChatID, latitude, longitude float64, title, address string, opts ...SendOption,
) (*Message, error) {
	return b.sendMedia(ctx, chat, message.Venue(latitude, longitude, 0, title, address), opts...)
}

// SendContact sends a phone contact.
func (b *Bot) SendContact(ctx context.Context, chat ChatID, phoneNumber, firstName, lastName string, opts ...SendOption) (*Message, error) {
	return b.sendMedia(ctx, chat, message.Contact(tg.InputMediaContact{
		PhoneNumber: phoneNumber,
		FirstName:   firstName,
		LastName:    lastName,
	}), opts...)
}

// SendPoll sends a native poll. At least two options are required.
func (b *Bot) SendPoll(ctx context.Context, chat ChatID, question string, options []string, opts ...SendOption) (*Message, error) {
	if len(options) < 2 {
		return nil, &Error{Code: 400, Description: "Bad Request: poll must have at least 2 option"}
	}

	answers := make([]message.PollAnswerOption, len(options))
	for i, o := range options {
		answers[i] = message.PollAnswer(o)
	}

	poll := message.Poll(question, answers[0], answers[1], answers[2:]...)

	return b.sendMedia(ctx, chat, poll, opts...)
}
