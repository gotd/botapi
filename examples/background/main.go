// Command background demonstrates proactive, restart-surviving sends with
// github.com/gotd/botapi: messages that are NOT replies to an incoming update.
//
// The trick is a serializable PeerRef. Sending to a chat needs its MTProto
// access hash; a PeerRef captures the chat id and that hash in a small
// JSON-serializable value. This bot persists one PeerRef per subscriber to a
// JSON file, so it can address them later — from a background goroutine, and
// even after a full restart — with no re-resolution and no live update to react
// to.
//
//	/subscribe    capture this chat's PeerRef and save it to disk
//	/unsubscribe  forget this chat
//
// While running, a background ticker (driven by Bot.Background, not a per-update
// handler context) broadcasts the time to every saved subscriber every 30s. On
// startup the bot reloads the file, deserializes each PeerRef, and sends a
// "back online" message — proving the references survive a restart.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token; SUBS_FILE overrides the store path (default ./subscribers.json):
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/background
package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

func main() {
	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("APP_ID must be a number (see https://my.telegram.org)", zap.Error(err))
	}

	// Persist session, peers and update state so the bot resumes across restarts.
	sess, err := storage.Open("background-session.bbolt")
	if err != nil {
		log.Fatal("Open storage", zap.Error(err))
	}
	defer func() { _ = sess.Close() }()

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:   appID,
		AppHash: os.Getenv("APP_HASH"),
		Logger:  logzap.New(log),
		Storage: sess,
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}
	bot.Use(botapi.Recover(), botapi.Logging())

	path := os.Getenv("SUBS_FILE")
	if path == "" {
		path = "subscribers.json"
	}
	store, err := loadStore(path)
	if err != nil {
		log.Fatal("Load subscribers", zap.Error(err))
	}

	bot.OnCommand("subscribe", "Receive background broadcasts", func(c *botapi.Context) error {
		chat, ok := c.Chat()
		if !ok {
			return nil
		}
		// Resolve the chat to a PeerRef once, here, while we have it — this is
		// what captures the access hash needed to message it later.
		ref, err := c.Bot.PeerRef(c, chat)
		if err != nil {
			return err
		}
		if store.add(ref) {
			if err := store.save(); err != nil {
				return err
			}
		}
		_, err = c.Reply("Subscribed. You'll get a broadcast every 30s, and a hello after each restart.")
		return err
	})

	bot.OnCommand("unsubscribe", "Stop background broadcasts", func(c *botapi.Context) error {
		chat, ok := c.Chat()
		if !ok {
			return nil
		}
		ref, err := c.Bot.PeerRef(c, chat)
		if err != nil {
			return err
		}
		if store.remove(ref) {
			if err := store.save(); err != nil {
				return err
			}
		}
		_, err = c.Reply("Unsubscribed.")
		return err
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Run the bot in the background so we can use its run-lifetime context for
	// proactive sends from this goroutine.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting background bot", zap.Int("subscribers", store.len()))
		if err := bot.Run(ctx); err != nil {
			log.Error("Run", zap.Error(err))
			cancel()
		}
	}()

	// Bot.Background blocks until Run has connected (or ctx is done), then yields
	// the run-lifetime context. Using it (not a handler context) is what makes
	// these sends "background": they outlive any single update.
	bg := waitReady(ctx, bot)
	if bg.Err() != nil {
		wg.Wait()
		return
	}

	// On startup, greet every saved subscriber straight from its deserialized
	// PeerRef — no incoming update, no re-resolution. This is the restart proof.
	for _, ref := range store.refs() {
		if _, err := bot.SendMessage(bg, botapi.Peer(ref), "👋 back online"); err != nil {
			log.Warn("Greet subscriber", zap.Error(err), zap.Int64("id", ref.ID))
		}
	}

	// A periodic broadcast driven by a ticker, not by any update.
	go broadcastLoop(bg, bot, store, log)

	wg.Wait()
}

// broadcastLoop sends the time to every subscriber every 30s until ctx is done.
func broadcastLoop(ctx context.Context, bot *botapi.Bot, store *store, log *zap.Logger) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			msg := "⏰ " + t.Format(time.RFC1123)
			for _, ref := range store.refs() {
				if _, err := bot.SendMessage(ctx, botapi.Peer(ref), msg); err != nil {
					log.Warn("Broadcast", zap.Error(err), zap.Int64("id", ref.ID))
				}
			}
		}
	}
}

// waitReady returns the bot's run-lifetime context once Run has connected, or a
// canceled context if ctx is done first.
func waitReady(ctx context.Context, bot *botapi.Bot) context.Context {
	for {
		if bg := bot.Background(); bg.Err() == nil {
			return bg
		}
		select {
		case <-ctx.Done():
			return ctx
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// store is a tiny JSON-file-backed set of subscriber PeerRefs. A real bot would
// use a database; the point here is only that PeerRef is serializable.
type store struct {
	path string
	mu   sync.Mutex
	byID map[int64]botapi.PeerRef
}

// loadStore reads the subscriber file, or starts empty if it does not exist.
func loadStore(path string) (*store, error) {
	s := &store{path: path, byID: map[int64]botapi.PeerRef{}}
	data, err := os.ReadFile(path) //nolint:gosec // example: path is an operator-provided env var
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var refs []botapi.PeerRef
	if err := json.Unmarshal(data, &refs); err != nil {
		return nil, err
	}
	for _, ref := range refs {
		s.byID[ref.ID] = ref
	}
	return s, nil
}

// save atomically writes the current subscriber set back to disk.
func (s *store) save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.MarshalIndent(s.list(), "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil { //nolint:gosec // example: operator-provided path
		return err
	}
	return os.Rename(tmp, s.path) //nolint:gosec // example: operator-provided path
}

// add records a subscriber, reporting whether it was newly added.
func (s *store) add(ref botapi.PeerRef) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[ref.ID]; ok {
		return false
	}
	s.byID[ref.ID] = ref
	return true
}

// remove drops a subscriber, reporting whether it was present.
func (s *store) remove(ref botapi.PeerRef) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[ref.ID]; !ok {
		return false
	}
	delete(s.byID, ref.ID)
	return true
}

// refs returns a snapshot of the subscriber references.
func (s *store) refs() []botapi.PeerRef {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.list()
}

// len reports the number of subscribers.
func (s *store) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.byID)
}

// list copies the references; callers must hold the lock.
func (s *store) list() []botapi.PeerRef {
	refs := make([]botapi.PeerRef, 0, len(s.byID))
	for _, ref := range s.byID {
		refs = append(refs, ref)
	}
	return refs
}
