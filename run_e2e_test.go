package botapi

import (
	"context"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/stretchr/testify/require"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"

	"github.com/gotd/teled/teledtest"
)

// TestRunEndToEnd boots a Bot against an in-process teled server (a real
// MTProto+RPC implementation, not hand-stubbed handlers): it logs in as a bot
// by token, fetches self, starts gap recovery, and — from OnStart — sends a
// real message over the encrypted transport to a freshly signed-up human, whose
// own session then reads it back. This covers the whole Run lifecycle against
// the actual server behavior teledtest provides.
//
// The test is skipped on hosts without container support (teledtest backs the
// server with a throwaway PostgreSQL container).
func TestRunEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	srv := teledtest.New(t) // skips the test when containers are unavailable.

	// Sign up a human recipient first, so the bot has a real peer to message.
	// teled access hashes are global, so the one returned here is valid for the
	// bot to address the user with. Keep the session so we can read the message
	// back from the recipient's side afterwards.
	recipientStorage := &session.StorageMemory{}
	var recipient *tg.User
	require.NoError(t, srv.Run(ctx, recipientStorage, func(api *tg.Client) error {
		recipient = signUpUser(ctx, t, api, "+1234500000", "Recipient")
		return nil
	}))

	// Build a Bot pointed at the in-process server. The resolver/publicKeys/
	// dcList seams replace Telegram's real DCs with teledtest's listener; the
	// token auto-provisions a bot account on first login.
	bot, err := New("424242:secret-bot-token", Options{
		AppID:      telegram.TestAppID,
		AppHash:    telegram.TestAppHash,
		resolver:   dcs.Plain(dcs.PlainOptions{}),
		publicKeys: srv.Keys,
		dcList: dcs.List{Options: []tg.DCOption{
			{ID: srv.DC, IPAddress: srv.Addr.IP.String(), Port: srv.Addr.Port},
		}},
		// No OnCommand handlers are registered, so the command set is empty and
		// registration is a no-op; leaving it enabled exercises that path too.
	})
	require.NoError(t, err)

	var sent *Message
	runCtx, runCancel := context.WithCancel(ctx)
	defer runCancel()
	bot.onStart = func(ctx context.Context) {
		m, err := bot.SendMessage(ctx, userRef(recipient.ID, recipient.AccessHash), "hello from OnStart")
		if err != nil {
			t.Errorf("send in OnStart: %v", err)
		}
		sent = m
		runCancel()
	}

	if err := bot.Run(runCtx); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("run: %v", err)
	}

	// The bot learned its identity during Run.
	require.NotNil(t, bot.Self(), "self")
	require.True(t, bot.Self().Bot, "self is a bot")

	// OnStart sent the message and got a real Message back.
	require.NotNil(t, sent, "no message was sent")
	require.Equal(t, "hello from OnStart", sent.Text)

	// The recipient's own session reads the message back from its history. Its
	// input peer to the bot uses the bot's (global) access hash from Self.
	require.NoError(t, srv.Run(ctx, recipientStorage, func(api *tg.Client) error {
		hist, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer:  &tg.InputPeerUser{UserID: bot.Self().ID, AccessHash: bot.Self().AccessHash},
			Limit: 1,
		})
		require.NoError(t, err)
		msgs := hist.(*tg.MessagesMessages).Messages
		require.NotEmpty(t, msgs)
		require.Equal(t, "hello from OnStart", msgs[0].(*tg.Message).Message)
		return nil
	}))
}

// signUpUser registers a fresh account via teled's dev auth flow and returns
// self.
func signUpUser(ctx context.Context, t *testing.T, api *tg.Client, phone, first string) *tg.User {
	t.Helper()
	sent, err := api.AuthSendCode(ctx, &tg.AuthSendCodeRequest{
		PhoneNumber: phone,
		APIID:       telegram.TestAppID,
		APIHash:     telegram.TestAppHash,
		Settings:    tg.CodeSettings{},
	})
	require.NoError(t, err)
	code := sent.(*tg.AuthSentCode)
	authResp, err := api.AuthSignUp(ctx, &tg.AuthSignUpRequest{
		PhoneNumber:   phone,
		PhoneCodeHash: code.PhoneCodeHash,
		FirstName:     first,
	})
	require.NoError(t, err)
	return authResp.(*tg.AuthAuthorization).User.(*tg.User)
}
