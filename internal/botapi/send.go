package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

// SendAnimation implements oas.Handler.
func (b *BotAPI) SendAnimation(ctx context.Context, req oas.SendAnimation) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendAudio implements oas.Handler.
func (b *BotAPI) SendAudio(ctx context.Context, req oas.SendAudio) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendChatAction implements oas.Handler.
func (b *BotAPI) SendChatAction(ctx context.Context, req oas.SendChatAction) (oas.Result, error) {
	p, err := b.resolveID(ctx, req.ChatID)
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "resolve chatID")
	}

	s := b.sender.To(p.InputPeer()).TypingAction()
	progress := 0
	switch req.Action {
	case "cancel":
		err = s.Cancel(ctx)
	case "typing":
		err = s.Typing(ctx)
	case "record_video":
		err = s.RecordVideo(ctx)
	case "upload_video":
		err = s.UploadVideo(ctx, progress)
	case "record_audio", "record_voice":
		err = s.RecordAudio(ctx)
	case "upload_audio", "upload_voice":
		err = s.UploadVideo(ctx, progress)
	case "upload_photo":
		err = s.UploadPhoto(ctx, progress)
	case "upload_document":
		err = s.UploadDocument(ctx, progress)
	case "choose_sticker":
		err = s.ChooseSticker(ctx)
	case "pick_up_location", "find_location":
		err = s.GeoLocation(ctx)
	case "record_video_note":
		err = s.RecordRound(ctx)
	case "upload_video_note":
		err = s.UploadRound(ctx, progress)
	default:
		return oas.Result{}, &BadRequestError{"Wrong parameter action in request"}
	}
	if err != nil {
		return oas.Result{}, errors.Wrap(err, "send action")
	}

	return resultOK(true), nil
}

// SendContact implements oas.Handler.
func (b *BotAPI) SendContact(ctx context.Context, req oas.SendContact) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendDice implements oas.Handler.
func (b *BotAPI) SendDice(ctx context.Context, req oas.SendDice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendDocument implements oas.Handler.
func (b *BotAPI) SendDocument(ctx context.Context, req oas.SendDocument) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendGame implements oas.Handler.
func (b *BotAPI) SendGame(ctx context.Context, req oas.SendGame) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendInvoice implements oas.Handler.
func (b *BotAPI) SendInvoice(ctx context.Context, req oas.SendInvoice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendLocation implements oas.Handler.
func (b *BotAPI) SendLocation(ctx context.Context, req oas.SendLocation) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendMediaGroup implements oas.Handler.
func (b *BotAPI) SendMediaGroup(ctx context.Context, req oas.SendMediaGroup) (oas.ResultArrayOfMessage, error) {
	return oas.ResultArrayOfMessage{}, &NotImplementedError{}
}

// SendMessage implements oas.Handler.
func (b *BotAPI) SendMessage(ctx context.Context, req oas.SendMessage) (oas.ResultMessage, error) {
	parseMode, isParseModeSet := req.ParseMode.Get()
	if isParseModeSet && parseMode != "HTML" {
		return oas.ResultMessage{}, &NotImplementedError{Message: "only HTML formatting is supported"}
	}

	p, err := b.resolveID(ctx, req.ChatID)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "resolve chatID")
	}
	s := &b.sender.To(p.InputPeer()).Builder

	if v := req.DisableWebPagePreview.Or(false); v {
		s = s.NoWebpage()
	}
	if v := req.DisableNotification.Or(false); v {
		s = s.Silent()
	}
	if v, ok := req.ReplyToMessageID.Get(); ok {
		s = s.Reply(v)
	}
	if m := req.ReplyMarkup; m != nil {
		mkp, err := b.convertToTelegramReplyMarkup(ctx, m)
		if err != nil {
			return oas.ResultMessage{}, errors.Wrap(err, "convert markup")
		}
		s = s.Markup(mkp)
	}

	var resp tg.UpdatesClass
	if isParseModeSet {
		// FIXME(tdakkota): random_id unpacking.
		resp, err = s.StyledText(ctx, html.String(b.peers.UserResolveHook(ctx), req.Text))
	} else {
		// FIXME(tdakkota): get entities from request.
		resp, err = s.Text(ctx, req.Text)
	}

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

// SendPhoto implements oas.Handler.
func (b *BotAPI) SendPhoto(ctx context.Context, req oas.SendPhoto) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendPoll implements oas.Handler.
func (b *BotAPI) SendPoll(ctx context.Context, req oas.SendPoll) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendSticker implements oas.Handler.
func (b *BotAPI) SendSticker(ctx context.Context, req oas.SendSticker) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVenue implements oas.Handler.
func (b *BotAPI) SendVenue(ctx context.Context, req oas.SendVenue) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVideo implements oas.Handler.
func (b *BotAPI) SendVideo(ctx context.Context, req oas.SendVideo) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVideoNote implements oas.Handler.
func (b *BotAPI) SendVideoNote(ctx context.Context, req oas.SendVideoNote) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVoice implements oas.Handler.
func (b *BotAPI) SendVoice(ctx context.Context, req oas.SendVoice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}
