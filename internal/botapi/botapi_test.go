package botapi

import (
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"
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
		}.Build(raw),
		Options{
			Logger: logger.Named("botapi"),
		},
	)
}
