package botapi

import "context"

// SendChatAction tells the chat that the bot is performing an action (e.g.
// "typing"). The status is cleared automatically after a short period or when
// the next message is sent.
//
// The switch over ChatAction is exhaustive (exhaustive lint).
func (b *Bot) SendChatAction(ctx context.Context, chat ChatID, action ChatAction) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	a := b.sender.To(peer).TypingAction()
	switch action {
	case ChatActionTyping:
		err = a.Typing(ctx)
	case ChatActionUploadPhoto:
		err = a.UploadPhoto(ctx, 0)
	case ChatActionRecordVideo:
		err = a.RecordVideo(ctx)
	case ChatActionUploadVideo:
		err = a.UploadVideo(ctx, 0)
	case ChatActionRecordVoice:
		err = a.RecordAudio(ctx)
	case ChatActionUploadVoice:
		err = a.UploadAudio(ctx, 0)
	case ChatActionUploadDocument:
		err = a.UploadDocument(ctx, 0)
	case ChatActionChooseSticker:
		err = a.ChooseSticker(ctx)
	case ChatActionFindLocation:
		err = a.GeoLocation(ctx)
	case ChatActionRecordVideoNote:
		err = a.RecordRound(ctx)
	case ChatActionUploadVideoNote:
		err = a.UploadRound(ctx, 0)
	default:
		return &Error{Code: 400, Description: "Bad Request: wrong parameter action in request"}
	}
	return asAPIError(err)
}
