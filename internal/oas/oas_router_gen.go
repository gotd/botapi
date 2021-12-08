// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/otelogen"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// No-op definition for keeping imports.
var (
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = sort.Ints
	_ = http.MethodGet
	_ = io.Copy
	_ = json.Marshal
	_ = bytes.NewReader
	_ = strconv.ParseInt
	_ = time.Time{}
	_ = conv.ToInt32
	_ = uuid.UUID{}
	_ = uri.PathEncoder{}
	_ = url.URL{}
	_ = math.Mod
	_ = validate.Int{}
	_ = ht.NewRequest
	_ = net.IP{}
	_ = otelogen.Version
	_ = trace.TraceIDFromHex
	_ = otel.GetTracerProvider
	_ = metric.NewNoopMeterProvider
	_ = regexp.MustCompile
	_ = jx.Null
	_ = sync.Pool{}
)

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func skipSlash(p []byte) []byte {
	if len(p) > 0 && p[0] == '/' {
		return p[1:]
	}
	return p
}

// nextElem return next path element from p and forwarded p.
func nextElem(p []byte) (elem, next []byte) {
	p = skipSlash(p)
	idx := bytes.IndexByte(p, '/')
	if idx < 0 {
		idx = len(p)
	}
	return p[:idx], p[idx:]
}

// ServeHTTP serves http request as defined by OpenAPI v3 specification,
// calling handler that matches the path or returning not found error.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := []byte(r.URL.Path)
	if len(p) == 0 {
		s.notFound(w, r)
		return
	}

	var (
		elem []byte            // current element, without slashes
		args map[string]string // lazily initialized
	)

	// Static code generated router with unwrapped path search.
	switch r.Method {
	case "POST":
		// Root edge.
		elem, p = nextElem(p)
		switch string(elem) {
		case "addStickerToSet": // -> 1
			// POST /addStickerToSet
			s.handleAddStickerToSetRequest(args, w, r)
			return
		case "answerCallbackQuery": // -> 2
			// POST /answerCallbackQuery
			s.handleAnswerCallbackQueryRequest(args, w, r)
			return
		case "answerInlineQuery": // -> 3
			// POST /answerInlineQuery
			s.handleAnswerInlineQueryRequest(args, w, r)
			return
		case "answerPreCheckoutQuery": // -> 4
			// POST /answerPreCheckoutQuery
			s.handleAnswerPreCheckoutQueryRequest(args, w, r)
			return
		case "answerShippingQuery": // -> 5
			// POST /answerShippingQuery
			s.handleAnswerShippingQueryRequest(args, w, r)
			return
		case "approveChatJoinRequest": // -> 6
			// POST /approveChatJoinRequest
			s.handleApproveChatJoinRequestRequest(args, w, r)
			return
		case "banChatMember": // -> 7
			// POST /banChatMember
			s.handleBanChatMemberRequest(args, w, r)
			return
		case "banChatSenderChat": // -> 8
			// POST /banChatSenderChat
			s.handleBanChatSenderChatRequest(args, w, r)
			return
		case "copyMessage": // -> 9
			// POST /copyMessage
			s.handleCopyMessageRequest(args, w, r)
			return
		case "createChatInviteLink": // -> 10
			// POST /createChatInviteLink
			s.handleCreateChatInviteLinkRequest(args, w, r)
			return
		case "createNewStickerSet": // -> 11
			// POST /createNewStickerSet
			s.handleCreateNewStickerSetRequest(args, w, r)
			return
		case "declineChatJoinRequest": // -> 12
			// POST /declineChatJoinRequest
			s.handleDeclineChatJoinRequestRequest(args, w, r)
			return
		case "deleteChatPhoto": // -> 13
			// POST /deleteChatPhoto
			s.handleDeleteChatPhotoRequest(args, w, r)
			return
		case "deleteChatStickerSet": // -> 14
			// POST /deleteChatStickerSet
			s.handleDeleteChatStickerSetRequest(args, w, r)
			return
		case "deleteMessage": // -> 15
			// POST /deleteMessage
			s.handleDeleteMessageRequest(args, w, r)
			return
		case "deleteMyCommands": // -> 16
			// POST /deleteMyCommands
			s.handleDeleteMyCommandsRequest(args, w, r)
			return
		case "deleteStickerFromSet": // -> 17
			// POST /deleteStickerFromSet
			s.handleDeleteStickerFromSetRequest(args, w, r)
			return
		case "deleteWebhook": // -> 18
			// POST /deleteWebhook
			s.handleDeleteWebhookRequest(args, w, r)
			return
		case "editChatInviteLink": // -> 19
			// POST /editChatInviteLink
			s.handleEditChatInviteLinkRequest(args, w, r)
			return
		case "editMessageCaption": // -> 20
			// POST /editMessageCaption
			s.handleEditMessageCaptionRequest(args, w, r)
			return
		case "editMessageLiveLocation": // -> 21
			// POST /editMessageLiveLocation
			s.handleEditMessageLiveLocationRequest(args, w, r)
			return
		case "editMessageMedia": // -> 22
			// POST /editMessageMedia
			s.handleEditMessageMediaRequest(args, w, r)
			return
		case "editMessageReplyMarkup": // -> 23
			// POST /editMessageReplyMarkup
			s.handleEditMessageReplyMarkupRequest(args, w, r)
			return
		case "editMessageText": // -> 24
			// POST /editMessageText
			s.handleEditMessageTextRequest(args, w, r)
			return
		case "exportChatInviteLink": // -> 25
			// POST /exportChatInviteLink
			s.handleExportChatInviteLinkRequest(args, w, r)
			return
		case "forwardMessage": // -> 26
			// POST /forwardMessage
			s.handleForwardMessageRequest(args, w, r)
			return
		case "getChat": // -> 27
			// POST /getChat
			s.handleGetChatRequest(args, w, r)
			return
		case "getChatAdministrators": // -> 28
			// POST /getChatAdministrators
			s.handleGetChatAdministratorsRequest(args, w, r)
			return
		case "getChatMember": // -> 29
			// POST /getChatMember
			s.handleGetChatMemberRequest(args, w, r)
			return
		case "getChatMemberCount": // -> 30
			// POST /getChatMemberCount
			s.handleGetChatMemberCountRequest(args, w, r)
			return
		case "getFile": // -> 31
			// POST /getFile
			s.handleGetFileRequest(args, w, r)
			return
		case "getGameHighScores": // -> 32
			// POST /getGameHighScores
			s.handleGetGameHighScoresRequest(args, w, r)
			return
		case "getMe": // -> 33
			// POST /getMe
			s.handleGetMeRequest(args, w, r)
			return
		case "getMyCommands": // -> 34
			// POST /getMyCommands
			s.handleGetMyCommandsRequest(args, w, r)
			return
		case "getStickerSet": // -> 35
			// POST /getStickerSet
			s.handleGetStickerSetRequest(args, w, r)
			return
		case "getUpdates": // -> 36
			// POST /getUpdates
			s.handleGetUpdatesRequest(args, w, r)
			return
		case "getUserProfilePhotos": // -> 37
			// POST /getUserProfilePhotos
			s.handleGetUserProfilePhotosRequest(args, w, r)
			return
		case "leaveChat": // -> 38
			// POST /leaveChat
			s.handleLeaveChatRequest(args, w, r)
			return
		case "pinChatMessage": // -> 39
			// POST /pinChatMessage
			s.handlePinChatMessageRequest(args, w, r)
			return
		case "promoteChatMember": // -> 40
			// POST /promoteChatMember
			s.handlePromoteChatMemberRequest(args, w, r)
			return
		case "restrictChatMember": // -> 41
			// POST /restrictChatMember
			s.handleRestrictChatMemberRequest(args, w, r)
			return
		case "revokeChatInviteLink": // -> 42
			// POST /revokeChatInviteLink
			s.handleRevokeChatInviteLinkRequest(args, w, r)
			return
		case "sendAnimation": // -> 43
			// POST /sendAnimation
			s.handleSendAnimationRequest(args, w, r)
			return
		case "sendAudio": // -> 44
			// POST /sendAudio
			s.handleSendAudioRequest(args, w, r)
			return
		case "sendChatAction": // -> 45
			// POST /sendChatAction
			s.handleSendChatActionRequest(args, w, r)
			return
		case "sendContact": // -> 46
			// POST /sendContact
			s.handleSendContactRequest(args, w, r)
			return
		case "sendDice": // -> 47
			// POST /sendDice
			s.handleSendDiceRequest(args, w, r)
			return
		case "sendDocument": // -> 48
			// POST /sendDocument
			s.handleSendDocumentRequest(args, w, r)
			return
		case "sendGame": // -> 49
			// POST /sendGame
			s.handleSendGameRequest(args, w, r)
			return
		case "sendInvoice": // -> 50
			// POST /sendInvoice
			s.handleSendInvoiceRequest(args, w, r)
			return
		case "sendLocation": // -> 51
			// POST /sendLocation
			s.handleSendLocationRequest(args, w, r)
			return
		case "sendMediaGroup": // -> 52
			// POST /sendMediaGroup
			s.handleSendMediaGroupRequest(args, w, r)
			return
		case "sendMessage": // -> 53
			// POST /sendMessage
			s.handleSendMessageRequest(args, w, r)
			return
		case "sendPhoto": // -> 54
			// POST /sendPhoto
			s.handleSendPhotoRequest(args, w, r)
			return
		case "sendPoll": // -> 55
			// POST /sendPoll
			s.handleSendPollRequest(args, w, r)
			return
		case "sendSticker": // -> 56
			// POST /sendSticker
			s.handleSendStickerRequest(args, w, r)
			return
		case "sendVenue": // -> 57
			// POST /sendVenue
			s.handleSendVenueRequest(args, w, r)
			return
		case "sendVideo": // -> 58
			// POST /sendVideo
			s.handleSendVideoRequest(args, w, r)
			return
		case "sendVideoNote": // -> 59
			// POST /sendVideoNote
			s.handleSendVideoNoteRequest(args, w, r)
			return
		case "sendVoice": // -> 60
			// POST /sendVoice
			s.handleSendVoiceRequest(args, w, r)
			return
		case "setChatAdministratorCustomTitle": // -> 61
			// POST /setChatAdministratorCustomTitle
			s.handleSetChatAdministratorCustomTitleRequest(args, w, r)
			return
		case "setChatDescription": // -> 62
			// POST /setChatDescription
			s.handleSetChatDescriptionRequest(args, w, r)
			return
		case "setChatPermissions": // -> 63
			// POST /setChatPermissions
			s.handleSetChatPermissionsRequest(args, w, r)
			return
		case "setChatPhoto": // -> 64
			// POST /setChatPhoto
			s.handleSetChatPhotoRequest(args, w, r)
			return
		case "setChatStickerSet": // -> 65
			// POST /setChatStickerSet
			s.handleSetChatStickerSetRequest(args, w, r)
			return
		case "setChatTitle": // -> 66
			// POST /setChatTitle
			s.handleSetChatTitleRequest(args, w, r)
			return
		case "setGameScore": // -> 67
			// POST /setGameScore
			s.handleSetGameScoreRequest(args, w, r)
			return
		case "setMyCommands": // -> 68
			// POST /setMyCommands
			s.handleSetMyCommandsRequest(args, w, r)
			return
		case "setPassportDataErrors": // -> 69
			// POST /setPassportDataErrors
			s.handleSetPassportDataErrorsRequest(args, w, r)
			return
		case "setStickerPositionInSet": // -> 70
			// POST /setStickerPositionInSet
			s.handleSetStickerPositionInSetRequest(args, w, r)
			return
		case "setStickerSetThumb": // -> 71
			// POST /setStickerSetThumb
			s.handleSetStickerSetThumbRequest(args, w, r)
			return
		case "setWebhook": // -> 72
			// POST /setWebhook
			s.handleSetWebhookRequest(args, w, r)
			return
		case "stopMessageLiveLocation": // -> 73
			// POST /stopMessageLiveLocation
			s.handleStopMessageLiveLocationRequest(args, w, r)
			return
		case "stopPoll": // -> 74
			// POST /stopPoll
			s.handleStopPollRequest(args, w, r)
			return
		case "unbanChatMember": // -> 75
			// POST /unbanChatMember
			s.handleUnbanChatMemberRequest(args, w, r)
			return
		case "unbanChatSenderChat": // -> 76
			// POST /unbanChatSenderChat
			s.handleUnbanChatSenderChatRequest(args, w, r)
			return
		case "unpinAllChatMessages": // -> 77
			// POST /unpinAllChatMessages
			s.handleUnpinAllChatMessagesRequest(args, w, r)
			return
		case "unpinChatMessage": // -> 78
			// POST /unpinChatMessage
			s.handleUnpinChatMessageRequest(args, w, r)
			return
		case "uploadStickerFile": // -> 79
			// POST /uploadStickerFile
			s.handleUploadStickerFileRequest(args, w, r)
			return
		default:
			s.notFound(w, r)
			return
		}
	default:
		s.notFound(w, r)
		return
	}
}