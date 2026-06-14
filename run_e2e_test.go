package botapi

import (
	"context"
	"testing"
	"time"

	"github.com/go-faster/errors"

	"github.com/gotd/td/tdsync"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgtest"
	"github.com/gotd/td/tgtest/cluster"
)

// TestRunEndToEnd boots a Bot against an in-process tgtest cluster: it logs in
// as a bot, fetches self, starts gap recovery, and (from OnStart) sends a real
// message over the encrypted MTProto transport — covering the Run lifecycle.
func TestRunEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := cluster.NewCluster(cluster.Options{})
	// Create the primary DC the client dials; handlers live on Common, which
	// every DC falls back to.
	c.Dispatch(2, "dc2")

	botUser := &tg.User{ID: 1, AccessHash: 1, Bot: true, Username: "test_bot"}
	common := c.Common()
	common.HandleFunc(tg.AuthImportBotAuthorizationRequestTypeID,
		func(s *tgtest.Server, req *tgtest.Request) error {
			return s.SendResult(req, &tg.AuthAuthorization{User: botUser})
		})
	common.Vector(tg.UsersGetUsersRequestTypeID, botUser)
	common.HandleFunc(tg.UpdatesGetStateRequestTypeID,
		func(s *tgtest.Server, req *tgtest.Request) error {
			return s.SendResult(req, &tg.UpdatesState{Pts: 1, Qts: 1, Seq: 1, Date: 1})
		})
	common.HandleFunc(tg.UpdatesGetDifferenceRequestTypeID,
		func(s *tgtest.Server, req *tgtest.Request) error {
			return s.SendResult(req, &tg.UpdatesDifferenceEmpty{Seq: 1, Date: 1})
		})
	// Be permissive about any other bootstrap RPC the client issues.
	common.Fallback(tgtest.HandlerFunc(func(s *tgtest.Server, req *tgtest.Request) error {
		return s.SendResult(req, &tg.BoolTrue{})
	}))

	sent := make(chan *tg.MessagesSendMessageRequest, 1)
	common.HandleFunc(tg.MessagesSendMessageRequestTypeID,
		func(s *tgtest.Server, req *tgtest.Request) error {
			m := &tg.MessagesSendMessageRequest{}
			if err := m.Decode(req.Buf); err != nil {
				return err
			}
			select {
			case sent <- m:
			default:
			}
			return s.SendResult(req, &tg.Updates{
				Updates: []tg.UpdateClass{&tg.UpdateNewMessage{
					Message: &tg.Message{ID: 100, Message: m.Message, PeerID: &tg.PeerUser{UserID: 10}},
				}},
				Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 20}},
			})
		})

	var bot *Bot
	g := tdsync.NewCancellableGroup(ctx)
	g.Go(c.Up)
	g.Go(func(ctx context.Context) error {
		select {
		case <-c.Ready():
		case <-ctx.Done():
			return ctx.Err()
		}

		// Build the bot only once the cluster is up: List() carries the DC
		// addresses that Up fills in.
		var err error
		bot, err = New("123:token", Options{
			AppID:                      1,
			AppHash:                    "hash",
			DisableCommandRegistration: true,
			resolver:                   c.Resolver(),
			publicKeys:                 c.Keys(),
			dcList:                     c.List(),
		})
		if err != nil {
			return err
		}
		bot.onStart = func(ctx context.Context) {
			if _, err := bot.SendMessage(ctx, userRef(10, 20), "hello from OnStart"); err != nil {
				t.Errorf("send in OnStart: %v", err)
			}
			cancel()
		}
		return bot.Run(ctx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("run: %v", err)
	}

	select {
	case m := <-sent:
		if m.Message != "hello from OnStart" {
			t.Fatalf("sent message = %q", m.Message)
		}
	default:
		t.Fatal("no message was sent")
	}

	// Sanity: the bot learned its identity during Run.
	if bot.Self() == nil || bot.Self().Username != "test_bot" {
		t.Fatalf("self = %#v", bot.Self())
	}
}
