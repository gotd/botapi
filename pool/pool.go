// Package pool runs and multiplexes many bots by token over a single process,
// lazily starting a Bot per token and garbage-collecting idle ones.
//
// It is the multi-bot front end for github.com/gotd/botapi: a local Bot API
// server, or any service that serves many bots, can hold one Pool and call Do
// with a token to borrow a connected Bot.
package pool

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/gotd/log"
	"go.etcd.io/bbolt"

	"github.com/gotd/botapi"
	"github.com/gotd/botapi/storage"
)

// Options configures a Pool.
type Options struct {
	// AppID and AppHash are the shared MTProto app identity used for every bot
	// (https://my.telegram.org). Required.
	AppID   int
	AppHash string

	// Logger is the base structured logger (github.com/gotd/log port).
	// Defaults to a no-op logger.
	Logger log.Logger

	// StateDir, when set, is the directory holding each bot's persistent session
	// file (<id>.bbolt). When empty, bots run with in-memory storage and nothing
	// survives a restart.
	StateDir string

	// IdleTimeout shuts down a bot that has not been borrowed for this long.
	// Zero disables idle collection (RunGC then does nothing).
	IdleTimeout time.Duration
}

// Pool is a concurrency-safe set of bots keyed by token.
type Pool struct {
	opt Options
	log log.Logger

	mu   sync.Mutex
	bots map[string]*managed
}

// New constructs an empty Pool. It performs no network I/O.
func New(opt Options) (*Pool, error) {
	if opt.AppID == 0 || opt.AppHash == "" {
		return nil, errors.New("AppID and AppHash are required")
	}
	opt.Logger = log.OrNop(opt.Logger)
	return &Pool{
		opt:  opt,
		log:  opt.Logger,
		bots: map[string]*managed{},
	}, nil
}

// Do borrows the running bot for token, starting and authorizing it on first
// use, and invokes fn with it. It blocks until the bot is ready, fn returns, or
// ctx is canceled. A failure to start the bot is returned to every concurrent
// caller waiting on it.
func (p *Pool) Do(ctx context.Context, token string, fn func(*botapi.Bot) error) error {
	tok, err := ParseToken(token)
	if err != nil {
		return err
	}

	m, err := p.acquire(tok)
	if err != nil {
		return err
	}

	select {
	case <-m.ready:
		if m.startErr != nil {
			return m.startErr
		}
		m.use()
		return fn(m.bot)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// acquire returns the managed bot for the token, creating and starting it if it
// does not exist yet.
func (p *Pool) acquire(tok Token) (*managed, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if m, ok := p.bots[tok.String()]; ok {
		return m, nil
	}

	m, err := p.start(tok)
	if err != nil {
		return nil, err
	}
	p.bots[tok.String()] = m
	return m, nil
}

// start builds and runs a bot for the token in the background.
func (p *Pool) start(tok Token) (*managed, error) {
	botLog := log.With(log.Named(p.log, "bot"), log.Int("id", tok.ID))

	var (
		store botapi.Storage
		db    *bbolt.DB
	)
	if p.opt.StateDir != "" {
		path := filepath.Join(p.opt.StateDir, fmt.Sprintf("%d.bbolt", tok.ID))
		opened, err := bbolt.Open(path, 0o666, bbolt.DefaultOptions)
		if err != nil {
			return nil, errors.Wrap(err, "open state db")
		}
		db = opened
		store = storage.NewBBoltStorage(opened)
	}

	m := &managed{ready: make(chan struct{}), db: db}

	bot, err := botapi.New(tok.String(), botapi.Options{
		AppID:   p.opt.AppID,
		AppHash: p.opt.AppHash,
		Logger:  botLog,
		Storage: store,
		OnStart: func(context.Context) { m.markReady(nil) },
	})
	if err != nil {
		if db != nil {
			_ = db.Close()
		}
		return nil, err
	}
	m.bot = bot

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	go func() {
		runErr := bot.Run(ctx)
		// If Run returns before OnStart fired, this is a startup failure that
		// must reach every waiter; markReady is a no-op once the bot is ready.
		if runErr != nil && !errors.Is(runErr, context.Canceled) {
			m.markReady(runErr)
			log.For(botLog).Warn(context.Background(), "Bot stopped", log.Error(runErr))
		} else {
			// Unblock any startup waiter on a clean shutdown too.
			m.markReady(errStopped)
		}
		p.drop(tok.String(), m)
		if m.db != nil {
			_ = m.db.Close()
		}
	}()

	return m, nil
}

// errStopped marks a bot that shut down before (or without) becoming ready.
var errStopped = errors.New("bot stopped before becoming ready")

// drop removes m from the pool if it is still the current entry for token and
// kills it.
func (p *Pool) drop(token string, m *managed) {
	p.mu.Lock()
	if cur, ok := p.bots[token]; ok && cur == m {
		delete(p.bots, token)
	}
	p.mu.Unlock()
	m.kill()
}

// Kill stops and removes the bot for the token, if present.
func (p *Pool) Kill(token string) {
	p.mu.Lock()
	m, ok := p.bots[token]
	if ok {
		delete(p.bots, token)
	}
	p.mu.Unlock()
	if ok {
		m.kill()
	}
}

// Close stops every bot in the pool.
func (p *Pool) Close() {
	p.mu.Lock()
	bots := p.bots
	p.bots = map[string]*managed{}
	p.mu.Unlock()
	for _, m := range bots {
		m.kill()
	}
}

// RunGC periodically reaps bots idle longer than IdleTimeout until ctx is
// canceled. It is a no-op when IdleTimeout is zero.
func (p *Pool) RunGC(ctx context.Context) {
	if p.opt.IdleTimeout <= 0 {
		<-ctx.Done()
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			p.reap(now.Add(-p.opt.IdleTimeout))
		}
	}
}

// reap kills every bot whose last use is before the deadline.
func (p *Pool) reap(deadline time.Time) {
	p.mu.Lock()
	var dead []*managed
	for token, m := range p.bots {
		if m.idleBefore(deadline) {
			dead = append(dead, m)
			delete(p.bots, token)
		}
	}
	p.mu.Unlock()
	for _, m := range dead {
		m.kill()
	}
}

// managed is one bot under the pool's control.
type managed struct {
	bot    *botapi.Bot
	cancel context.CancelFunc
	db     *bbolt.DB

	// ready is closed exactly once when the bot has either become ready or
	// failed to start; startErr (set before the close) carries the outcome.
	ready     chan struct{}
	readyOnce sync.Once
	startErr  error

	mu       sync.Mutex
	lastUsed time.Time
}

// markReady latches the startup outcome and unblocks every waiter. The first
// call wins; later calls (e.g. shutdown after a successful start) are no-ops.
func (m *managed) markReady(err error) {
	m.readyOnce.Do(func() {
		m.startErr = err
		close(m.ready)
	})
}

func (m *managed) use() {
	m.mu.Lock()
	m.lastUsed = time.Now()
	m.mu.Unlock()
}

func (m *managed) idleBefore(deadline time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return !m.lastUsed.IsZero() && m.lastUsed.Before(deadline)
}

func (m *managed) kill() {
	if m.cancel != nil {
		m.cancel()
	}
}
