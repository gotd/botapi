// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"github.com/go-faster/jx"
	"go.opentelemetry.io/otel/trace"
)

func encodeAddStickerToSetRequestJSON(
	req AddStickerToSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeAnswerCallbackQueryRequestJSON(
	req AnswerCallbackQuery,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeAnswerInlineQueryRequestJSON(
	req AnswerInlineQuery,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeAnswerPreCheckoutQueryRequestJSON(
	req AnswerPreCheckoutQuery,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeAnswerShippingQueryRequestJSON(
	req AnswerShippingQuery,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeAnswerWebAppQueryRequestJSON(
	req AnswerWebAppQuery,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeApproveChatJoinRequestRequestJSON(
	req ApproveChatJoinRequest,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeBanChatMemberRequestJSON(
	req BanChatMember,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeBanChatSenderChatRequestJSON(
	req BanChatSenderChat,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeCopyMessageRequestJSON(
	req CopyMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeCreateChatInviteLinkRequestJSON(
	req CreateChatInviteLink,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeCreateNewStickerSetRequestJSON(
	req CreateNewStickerSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeclineChatJoinRequestRequestJSON(
	req DeclineChatJoinRequest,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeleteChatPhotoRequestJSON(
	req DeleteChatPhoto,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeleteChatStickerSetRequestJSON(
	req DeleteChatStickerSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeleteMessageRequestJSON(
	req DeleteMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeleteMyCommandsRequestJSON(
	req OptDeleteMyCommands,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeDeleteStickerFromSetRequestJSON(
	req DeleteStickerFromSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeDeleteWebhookRequestJSON(
	req OptDeleteWebhook,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeEditChatInviteLinkRequestJSON(
	req EditChatInviteLink,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeEditMessageCaptionRequestJSON(
	req EditMessageCaption,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeEditMessageLiveLocationRequestJSON(
	req EditMessageLiveLocation,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeEditMessageMediaRequestJSON(
	req EditMessageMedia,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeEditMessageReplyMarkupRequestJSON(
	req EditMessageReplyMarkup,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeEditMessageTextRequestJSON(
	req EditMessageText,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeExportChatInviteLinkRequestJSON(
	req ExportChatInviteLink,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeForwardMessageRequestJSON(
	req ForwardMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetChatRequestJSON(
	req GetChat,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetChatAdministratorsRequestJSON(
	req GetChatAdministrators,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetChatMemberRequestJSON(
	req GetChatMember,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetChatMemberCountRequestJSON(
	req GetChatMemberCount,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetChatMenuButtonRequestJSON(
	req OptGetChatMenuButton,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeGetFileRequestJSON(
	req GetFile,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetGameHighScoresRequestJSON(
	req GetGameHighScores,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetMyCommandsRequestJSON(
	req OptGetMyCommands,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeGetMyDefaultAdministratorRightsRequestJSON(
	req OptGetMyDefaultAdministratorRights,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeGetStickerSetRequestJSON(
	req GetStickerSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeGetUpdatesRequestJSON(
	req OptGetUpdates,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeGetUserProfilePhotosRequestJSON(
	req GetUserProfilePhotos,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeLeaveChatRequestJSON(
	req LeaveChat,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodePinChatMessageRequestJSON(
	req PinChatMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodePromoteChatMemberRequestJSON(
	req PromoteChatMember,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeRestrictChatMemberRequestJSON(
	req RestrictChatMember,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeRevokeChatInviteLinkRequestJSON(
	req RevokeChatInviteLink,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendAnimationRequestJSON(
	req SendAnimation,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendAudioRequestJSON(
	req SendAudio,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendChatActionRequestJSON(
	req SendChatAction,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendContactRequestJSON(
	req SendContact,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendDiceRequestJSON(
	req SendDice,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendDocumentRequestJSON(
	req SendDocument,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendGameRequestJSON(
	req SendGame,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendInvoiceRequestJSON(
	req SendInvoice,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendLocationRequestJSON(
	req SendLocation,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendMediaGroupRequestJSON(
	req SendMediaGroup,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendMessageRequestJSON(
	req SendMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendPhotoRequestJSON(
	req SendPhoto,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendPollRequestJSON(
	req SendPoll,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendStickerRequestJSON(
	req SendSticker,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendVenueRequestJSON(
	req SendVenue,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendVideoRequestJSON(
	req SendVideo,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendVideoNoteRequestJSON(
	req SendVideoNote,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSendVoiceRequestJSON(
	req SendVoice,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatAdministratorCustomTitleRequestJSON(
	req SetChatAdministratorCustomTitle,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatDescriptionRequestJSON(
	req SetChatDescription,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatMenuButtonRequestJSON(
	req OptSetChatMenuButton,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeSetChatPermissionsRequestJSON(
	req SetChatPermissions,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatPhotoRequestJSON(
	req SetChatPhoto,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatStickerSetRequestJSON(
	req SetChatStickerSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetChatTitleRequestJSON(
	req SetChatTitle,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetGameScoreRequestJSON(
	req SetGameScore,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetMyCommandsRequestJSON(
	req SetMyCommands,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetMyDefaultAdministratorRightsRequestJSON(
	req OptSetMyDefaultAdministratorRights,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()
	if req.Set {
		req.Encode(e)
	}
	return e, nil
}
func encodeSetPassportDataErrorsRequestJSON(
	req SetPassportDataErrors,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetStickerPositionInSetRequestJSON(
	req SetStickerPositionInSet,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetStickerSetThumbRequestJSON(
	req SetStickerSetThumb,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeSetWebhookRequestJSON(
	req SetWebhook,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeStopMessageLiveLocationRequestJSON(
	req StopMessageLiveLocation,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeStopPollRequestJSON(
	req StopPoll,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeUnbanChatMemberRequestJSON(
	req UnbanChatMember,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeUnbanChatSenderChatRequestJSON(
	req UnbanChatSenderChat,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeUnpinAllChatMessagesRequestJSON(
	req UnpinAllChatMessages,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeUnpinChatMessageRequestJSON(
	req UnpinChatMessage,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
func encodeUploadStickerFileRequestJSON(
	req UploadStickerFile,
	span trace.Span,
) (
	data *jx.Encoder,

	rerr error,
) {
	e := jx.GetEncoder()

	req.Encode(e)
	return e, nil
}
