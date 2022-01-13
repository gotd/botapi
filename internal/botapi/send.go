package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/peers"

	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) sentMessage(
	ctx context.Context,
	p peers.Peer,
	resp tg.UpdatesClass, err error,
) (oas.ResultMessage, error) {
	m, err := unpack.MessageClass(resp, err)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "send")
	}

	msg, ok := m.(*tg.Message)
	if !ok {
		return oas.ResultMessage{
			Ok: true,
		}, nil
	}
	if msg.PeerID == nil {
		switch p := p.InputPeer().(type) {
		case *tg.InputPeerChat:
			msg.PeerID = &tg.PeerChat{ChatID: p.ChatID}
		case *tg.InputPeerUser:
			msg.PeerID = &tg.PeerUser{UserID: p.UserID}
		case *tg.InputPeerChannel:
			msg.PeerID = &tg.PeerChannel{ChannelID: p.ChannelID}
		}
	}

	resultMsg, err := b.convertPlainMessage(ctx, msg)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "get message")
	}
	return oas.ResultMessage{
		Result: oas.NewOptMessage(resultMsg),
		Ok:     true,
	}, nil
}

type sendOpts struct {
	To                       oas.ID
	DisableWebPagePreview    oas.OptBool
	DisableNotification      oas.OptBool
	ProtectContent           oas.OptBool
	ReplyToMessageID         oas.OptInt
	AllowSendingWithoutReply oas.OptBool
	ReplyMarkup              oas.OptSendReplyMarkup
}

func (b *BotAPI) prepareSend(
	ctx context.Context,
	req sendOpts,
) (*message.Builder, peers.Peer, error) {
	p, err := b.resolveID(ctx, req.To)
	if err != nil {
		return nil, nil, errors.Wrap(err, "resolve chatID")
	}
	s := &b.sender.To(p.InputPeer()).Builder

	if v := req.DisableWebPagePreview.Or(false); v {
		s = s.NoWebpage()
	}
	if v := req.DisableNotification.Or(false); v {
		s = s.Silent()
	}
	if v := req.ProtectContent.Or(false); v {
		s = s.NoForwards()
	}
	// TODO(tdakkota): check allow_sending_without_reply
	if v, ok := req.ReplyToMessageID.Get(); ok {
		s = s.Reply(v)
	}
	if m, ok := req.ReplyMarkup.Get(); ok {
		mkp, err := b.convertToTelegramReplyMarkup(ctx, m)
		if err != nil {
			return nil, nil, errors.Wrap(err, "convert markup")
		}
		s = s.Markup(mkp)
	}
	return s, p, nil
}

// SendMessage implements oas.Handler.
func (b *BotAPI) SendMessage(ctx context.Context, req oas.SendMessage) (oas.ResultMessage, error) {
	parseMode, isParseModeSet := req.ParseMode.Get()
	if isParseModeSet && parseMode != "HTML" {
		return oas.ResultMessage{}, &NotImplementedError{Message: "only HTML formatting is supported"}
	}

	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableWebPagePreview:    req.DisableWebPagePreview,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}

	var resp tg.UpdatesClass
	if isParseModeSet {
		// FIXME(tdakkota): random_id unpacking.
		resp, err = s.StyledText(ctx, html.String(b.peers.UserResolveHook(ctx), req.Text))
	} else {
		// FIXME(tdakkota): get entities from request.
		resp, err = s.Text(ctx, req.Text)
	}

	return b.sentMessage(ctx, p, resp, err)
}
