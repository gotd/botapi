package internal

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

var _ interface {
	CreateChatInviteLink(context.Context, oas.CreateChatInviteLink) (oas.ResultChatInviteLink, error)
	EditChatInviteLink(context.Context, oas.EditChatInviteLink) (oas.ResultChatInviteLink, error)
	ExportChatInviteLink(context.Context, oas.ExportChatInviteLink) (oas.ResultString, error)
	ForwardMessage(context.Context, oas.ForwardMessage) (oas.ResultMessage, error)
	GetChat(context.Context, oas.GetChat) (oas.ResultChat, error)
	GetChatAdministrators(context.Context, oas.GetChatAdministrators) (oas.ResultArrayOfChatMember, error)
	GetChatMember(context.Context, oas.GetChatMember) (oas.ResultChatMember, error)
	GetChatMemberCount(context.Context, oas.GetChatMemberCount) (oas.ResultInt, error)
	GetGameHighScores(context.Context, oas.GetGameHighScores) (oas.ResultArrayOfGameHighScore, error)
	GetMe(context.Context) (oas.ResultUser, error)
	GetMyCommands(context.Context, oas.GetMyCommands) (oas.ResultArrayOfBotCommand, error)
	GetUpdates(context.Context, oas.GetUpdates) (oas.ResultArrayOfUpdate, error)
	GetUserProfilePhotos(context.Context, oas.GetUserProfilePhotos) (oas.ResultUserProfilePhotos, error)
	GetWebhookInfo(context.Context) (oas.ResultWebhookInfo, error)
	RevokeChatInviteLink(context.Context, oas.RevokeChatInviteLink) (oas.ResultChatInviteLink, error)
	SendAnimation(context.Context, oas.SendAnimation) (oas.ResultMessage, error)
	SendAudio(context.Context, oas.SendAudio) (oas.ResultMessage, error)
	SendContact(context.Context, oas.SendContact) (oas.ResultMessage, error)
	SendDice(context.Context, oas.SendDice) (oas.ResultMessage, error)
	SendDocument(context.Context, oas.SendDocument) (oas.ResultMessage, error)
	SendGame(context.Context, oas.SendGame) (oas.ResultMessage, error)
	SendInvoice(context.Context, oas.SendInvoice) (oas.ResultMessage, error)
	SendLocation(context.Context, oas.SendLocation) (oas.ResultMessage, error)
	SendMediaGroup(context.Context, oas.SendMediaGroup) (oas.ResultArrayOfMessage, error)
	SendMessage(context.Context, oas.SendMessage) (oas.ResultMessage, error)
	SendPhoto(context.Context, oas.SendPhoto) (oas.ResultMessage, error)
	SendPoll(context.Context, oas.SendPoll) (oas.ResultMessage, error)
	SendSticker(context.Context, oas.SendSticker) (oas.ResultMessage, error)
	SendVenue(context.Context, oas.SendVenue) (oas.ResultMessage, error)
	SendVideo(context.Context, oas.SendVideo) (oas.ResultMessage, error)
	SendVideoNote(context.Context, oas.SendVideoNote) (oas.ResultMessage, error)
	SendVoice(context.Context, oas.SendVoice) (oas.ResultMessage, error)
	StopPoll(context.Context, oas.StopPoll) (oas.ResultPoll, error)
} = (*oas.Client)(nil)
