package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gotd/botapi/internal/oas"
)

var _ oas.Handler = (*BotAPI)(nil)

func TestUnimplemented(t *testing.T) {
	ctx := context.Background()
	a := assert.New(t)
	b := BotAPI{}

	{
		_, err := b.AddStickerToSet(ctx, &oas.AddStickerToSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.AnswerCallbackQuery(ctx, &oas.AnswerCallbackQuery{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.AnswerInlineQuery(ctx, &oas.AnswerInlineQuery{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.AnswerPreCheckoutQuery(ctx, &oas.AnswerPreCheckoutQuery{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.AnswerShippingQuery(ctx, &oas.AnswerShippingQuery{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.BanChatMember(ctx, &oas.BanChatMember{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.BanChatSenderChat(ctx, &oas.BanChatSenderChat{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.CopyMessage(ctx, &oas.CopyMessage{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.CreateNewStickerSet(ctx, &oas.CreateNewStickerSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.DeleteMessage(ctx, &oas.DeleteMessage{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.DeleteStickerFromSet(ctx, &oas.DeleteStickerFromSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.DeleteWebhook(ctx, oas.OptDeleteWebhook{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.EditMessageCaption(ctx, &oas.EditMessageCaption{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.EditMessageLiveLocation(ctx, &oas.EditMessageLiveLocation{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.EditMessageMedia(ctx, &oas.EditMessageMedia{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.EditMessageReplyMarkup(ctx, &oas.EditMessageReplyMarkup{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.EditMessageText(ctx, &oas.EditMessageText{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.ForwardMessage(ctx, &oas.ForwardMessage{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetChatAdministrators(ctx, &oas.GetChatAdministrators{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetChatMember(ctx, &oas.GetChatMember{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetFile(ctx, &oas.GetFile{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetGameHighScores(ctx, &oas.GetGameHighScores{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetStickerSet(ctx, &oas.GetStickerSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetUpdates(ctx, oas.OptGetUpdates{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.GetWebhookInfo(ctx)
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.PinChatMessage(ctx, &oas.PinChatMessage{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.PromoteChatMember(ctx, &oas.PromoteChatMember{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.RestrictChatMember(ctx, &oas.RestrictChatMember{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendAnimation(ctx, &oas.SendAnimation{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendAudio(ctx, &oas.SendAudio{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendDocument(ctx, &oas.SendDocument{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendMediaGroup(ctx, &oas.SendMediaGroup{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendPhoto(ctx, &oas.SendPhoto{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendSticker(ctx, &oas.SendSticker{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendVideo(ctx, &oas.SendVideo{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendVideoNote(ctx, &oas.SendVideoNote{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SendVoice(ctx, &oas.SendVoice{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetChatAdministratorCustomTitle(ctx, &oas.SetChatAdministratorCustomTitle{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetChatPermissions(ctx, &oas.SetChatPermissions{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetChatPhoto(ctx, &oas.SetChatPhoto{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetChatStickerSet(ctx, &oas.SetChatStickerSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetGameScore(ctx, &oas.SetGameScore{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetPassportDataErrors(ctx, &oas.SetPassportDataErrors{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetStickerPositionInSet(ctx, &oas.SetStickerPositionInSet{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetStickerSetThumb(ctx, &oas.SetStickerSetThumbnail{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.SetWebhook(ctx, &oas.SetWebhook{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.StopMessageLiveLocation(ctx, &oas.StopMessageLiveLocation{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.StopPoll(ctx, &oas.StopPoll{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.UnbanChatMember(ctx, &oas.UnbanChatMember{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.UnbanChatSenderChat(ctx, &oas.UnbanChatSenderChat{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.UnpinAllChatMessages(ctx, &oas.UnpinAllChatMessages{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.UnpinChatMessage(ctx, &oas.UnpinChatMessage{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}

	{
		_, err := b.UploadStickerFile(ctx, &oas.UploadStickerFile{})
		var implErr *NotImplementedError
		a.ErrorAs(err, &implErr)
	}
}
