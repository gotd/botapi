package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestEditMessageMediaPhotoURL(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)

	media := &InputMediaPhoto{Media: FileURL("https://e/p.jpg"), Caption: "cap"}
	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, media); err != nil {
		t.Fatalf("EditMessageMedia: %v", err)
	}

	var req tg.MessagesEditMessageRequest

	inv.decode(t, tg.MessagesEditMessageRequestTypeID, &req)

	if _, ok := req.Media.(*tg.InputMediaPhotoExternal); !ok {
		t.Fatalf("media = %#v", req.Media)
	}
}

func TestEditMessageMediaDocumentURL(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesEditMessageRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)

	media := &InputMediaDocument{Media: FileURL("https://e/f.pdf")}
	if _, err := b.EditMessageMedia(context.Background(), userRef(10, 20), 5, media); err != nil {
		t.Fatalf("EditMessageMedia doc: %v", err)
	}
}

func TestLeaveChat(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsLeaveChannelRequestTypeID, okUpdates())

	b := newMockBot(inv)

	if err := b.LeaveChat(context.Background(), tdlibChannel(50)); err != nil {
		t.Fatalf("LeaveChat: %v", err)
	}

	if !inv.called(tg.ChannelsLeaveChannelRequestTypeID) {
		t.Fatal("expected channels.leaveChannel")
	}
}

func TestSetMyCommandsScopes(t *testing.T) {
	scopes := []BotCommandScope{
		BotCommandScopeDefault(),
		BotCommandScopeAllPrivateChats(),
		BotCommandScopeAllGroupChats(),
		BotCommandScopeAllChatAdministrators(),
		BotCommandScopeChat(tdlibChannel(50)),
		BotCommandScopeChatAdministrators(tdlibChannel(50)),
		BotCommandScopeChatMember(tdlibChannel(50), 99),
	}
	for _, scope := range scopes {
		inv := newMockInvoker()
		inv.reply(tg.BotsSetBotCommandsRequestTypeID, &tg.BoolTrue{})

		b := newMockBot(inv)
		cmds := []BotCommand{{Command: "start", Description: "Start"}}

		err := b.SetMyCommands(context.Background(), cmds,
			WithCommandScope(scope), WithLanguageCode("en"))
		if err != nil {
			t.Fatalf("SetMyCommands(%T): %v", scope, err)
		}

		var req tg.BotsSetBotCommandsRequest

		inv.decode(t, tg.BotsSetBotCommandsRequestTypeID, &req)

		if req.LangCode != "en" {
			t.Fatalf("lang = %q", req.LangCode)
		}
	}
}

// TestOptionBuildersApply exercises the option constructors that other tests do
// not cover, by checking their effect on the produced MTProto request.
func TestCallbackOptionBuilders(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetBotCallbackAnswerRequestTypeID, &tg.BoolTrue{})

	b := newMockBot(inv)

	err := b.AnswerCallbackQuery(context.Background(), "12",
		WithCallbackURL("https://t.me/x"), WithCallbackCacheTime(30))
	if err != nil {
		t.Fatalf("AnswerCallbackQuery: %v", err)
	}

	var req tg.MessagesSetBotCallbackAnswerRequest

	inv.decode(t, tg.MessagesSetBotCallbackAnswerRequestTypeID, &req)

	if req.URL != "https://t.me/x" || req.CacheTime != 30 {
		t.Fatalf("req = %#v", req)
	}
}

func TestGetUserProfilePhotosOptions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PhotosGetUserPhotosRequestTypeID, &tg.PhotosPhotos{
		Users: []tg.UserClass{&tg.User{ID: 99, AccessHash: 1}},
	})

	b := newMockBot(inv)

	_, err := b.GetUserProfilePhotos(context.Background(), 99,
		WithProfilePhotosOffset(5), WithProfilePhotosLimit(10))
	if err != nil {
		t.Fatalf("GetUserProfilePhotos: %v", err)
	}

	var req tg.PhotosGetUserPhotosRequest

	inv.decode(t, tg.PhotosGetUserPhotosRequestTypeID, &req)

	if req.Offset != 5 || req.Limit != 10 {
		t.Fatalf("req = %#v", req)
	}
}

func TestBanChatMemberOptions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.ChannelsEditBannedRequestTypeID, okUpdates())
	inv.reply(tg.ChannelsDeleteParticipantHistoryRequestTypeID, &tg.MessagesAffectedHistory{})

	b := newMockBot(inv)

	err := b.BanChatMember(context.Background(), tdlibChannel(50), 99,
		WithBanUntil(1700000000), WithRevokeMessages())
	if err != nil {
		t.Fatalf("BanChatMember: %v", err)
	}

	var req tg.ChannelsEditBannedRequest

	inv.decode(t, tg.ChannelsEditBannedRequestTypeID, &req)

	if req.BannedRights.UntilDate != 1700000000 {
		t.Fatalf("until = %d", req.BannedRights.UntilDate)
	}
}

func TestSetGameScoreOptions(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.MessagesSetGameScoreRequestTypeID, editUpdates(&tg.Message{ID: 5, PeerID: &tg.PeerUser{UserID: 10}}))

	b := newMockBot(inv)

	_, err := b.SetGameScore(context.Background(), userRef(10, 20), 5, 99, 1000,
		WithForceScore(), WithoutEditMessage())
	if err != nil {
		t.Fatalf("SetGameScore: %v", err)
	}

	var req tg.MessagesSetGameScoreRequest

	inv.decode(t, tg.MessagesSetGameScoreRequestTypeID, &req)

	if !req.Force || req.EditMessage {
		t.Fatalf("req = %#v", req)
	}
}
