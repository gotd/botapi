package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/botapi/internal/oas"
)

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
