package botapi

import (
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/gotd/td/tgerr"

	"github.com/gotd/td/bin"

	"github.com/gotd/td/constant"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func testBotAPI(t *testing.T) (*tgmock.Mock, *BotAPI) {
	m := tgmock.New(t)
	raw := tg.NewClient(m)
	logger := zaptest.NewLogger(t)
	return m, NewBotAPI(
		raw,
		nil,
		peers.Options{
			Logger: logger.Named("peers"),
			Cache:  new(peers.InmemoryCache),
		}.Build(raw),
		Options{
			Logger: logger.Named("botapi"),
		},
	)
}

func testError() *tgerr.Error {
	return tgerr.New(1337, "TEST_ERROR")
}

func testChatID() int64 {
	var id constant.TDLibPeerID
	id.Chat(testChat().ID)
	return int64(id)
}

func testUser() *tg.User {
	u := &tg.User{
		Self:                 true,
		Bot:                  true,
		ID:                   10,
		AccessHash:           10,
		FirstName:            "Elsa",
		LastName:             "Jean",
		Username:             "thebot",
		BotInfoVersion:       1,
		BotInlinePlaceholder: "aboba",
	}
	u.SetFlags()
	return u
}

func testChat() *tg.Chat {
	u := &tg.Chat{
		Noforwards:        true,
		ID:                10,
		Title:             "I hate mondays",
		ParticipantsCount: 10,
		Date:              int(10),
		Version:           1,
		Photo:             &tg.ChatPhotoEmpty{},
	}
	u.SetFlags()
	return u
}

func testCommands() []tg.BotCommand {
	return []tg.BotCommand{
		{
			Command:     "freeburger",
			Description: "trolling",
		},
	}
}

func testCommandsBotAPI() []oas.BotCommand {
	return []oas.BotCommand{
		{
			Command:     "freeburger",
			Description: "trolling",
		},
	}
}

func setFlags(b bin.Object) {
	if v, ok := b.(interface {
		SetFlags()
	}); ok {
		v.SetFlags()
	}
}
