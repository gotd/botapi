// Package pool implements client pool.
package pool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"go.etcd.io/bbolt"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/botapi"
	"github.com/gotd/botapi/internal/botstorage"
)

// Pool of clients.
type Pool struct {
	appID   int
	appHash string
	debug   bool
	log     *zap.Logger

	statePath string

	clients    map[Token]*client
	clientsMux sync.Mutex
}

func (p *Pool) tick(deadline time.Time) {
	p.clientsMux.Lock()
	var toRemove []Token
	for token, c := range p.clients {
		if c.Deadline(deadline) {
			toRemove = append(toRemove, token)
			c.Kill()
		}
	}
	for _, token := range toRemove {
		delete(p.clients, token)
	}
	p.clientsMux.Unlock()
}

func (p *Pool) now() time.Time {
	return time.Now()
}

// Kill shutdowns client by token.
func (p *Pool) Kill(token Token) {
	p.clientsMux.Lock()
	defer p.clientsMux.Unlock()

	c, ok := p.clients[token]
	if !ok {
		return
	}
	c.Kill()
	delete(p.clients, token)
}

// Do acquires telegram client by token.
//
// Returns error if token is invalid. Block until client is available,
// authentication error or context cancelled.
func (p *Pool) Do(ctx context.Context, token Token, fn func(client *botapi.BotAPI) error) error {
	p.clientsMux.Lock()
	c, ok := p.clients[token]
	p.clientsMux.Unlock()

	if ok {
		// Happy path.
		c.Use(p.now())
		return fn(c.api)
	}

	initializationResult := make(chan error, 1)
	c, err := p.createClient(token, initializationResult)
	if err != nil {
		return errors.Wrap(err, "init")
	}

	// Waiting for initialization.
	select {
	case err, ok := <-initializationResult:
		if !ok {
			if c.ctx != nil {
				return c.ctx.Err()
			}
			return errors.New("unknown initialization error")
		}
		if err != nil {
			return err
		}

		p.clientsMux.Lock()
		conflictingClient, ok := p.clients[token]
		if ok {
			// Existing conflicting client, so stopping current client.
			c.Kill()
			c = conflictingClient
		} else {
			p.clients[token] = c
		}
		p.clientsMux.Unlock()

		return fn(c.api)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pool) RunGC(timeout time.Duration) {
	timer := time.NewTicker(time.Second)
	for now := range timer.C {
		deadline := now.Add(-timeout)
		p.tick(deadline)
	}
}

func (p *Pool) createClient(token Token, initializationResult chan<- error) (_ *client, rErr error) {
	log := p.log.Named("client").With(zap.Int("id", token.ID))

	dbPath := filepath.Join(p.statePath, fmt.Sprintf("%d.bbolt", token.ID))
	db, err := bbolt.Open(dbPath, 0o666, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rErr != nil {
			multierr.AppendInto(&rErr, db.Close())
		}
	}()
	storage := botstorage.NewBBoltStorage(db)

	var handler telegram.UpdateHandlerFunc = func(ctx context.Context, u tg.UpdatesClass) error {
		return nil
	}
	pClient := new(tg.Client)
	peerManager := peers.Options{
		Storage: storage,
		Cache:   storage,
		Logger:  log.Named("peers"),
	}.Build(pClient)
	gaps := updates.New(updates.Config{
		Handler: handler,
		OnChannelTooLong: func(channelID int64) {
			log.Warn("Got channel too long", zap.Int64("channel_id", channelID))
		},
		Storage:      storage,
		AccessHasher: peerManager,
		Logger:       log.Named("gaps"),
	})
	h := peerManager.UpdateHook(gaps)
	options := telegram.Options{
		Logger:         log.Named("client"),
		UpdateHandler:  h,
		SessionStorage: storage,
		Middlewares: []telegram.Middleware{
			updhook.UpdateHook(h.Handle),
		},
	}
	tgClient := telegram.NewClient(p.appID, p.appHash, options)
	// FIXME(tdakkota): fix this
	*pClient = *tgClient.API()

	tgContext, tgCancel := context.WithCancel(context.Background())
	c := &client{
		ctx:    tgContext,
		cancel: tgCancel,
		api: botapi.NewBotAPI(tgClient.API(), gaps, peerManager, botapi.Options{
			Debug:  p.debug,
			Logger: log.Named("botapi"),
		}),
		client:   tgClient,
		token:    token,
		lastUsed: time.Time{},
	}

	// Wait for initialization.
	go func() {
		defer func() {
			// Removing client from client list on close.
			p.clientsMux.Lock()
			found, ok := p.clients[token]
			if ok && found.client == c.client {
				delete(p.clients, token)
			}
			p.clientsMux.Unlock()
			// Kill client.
			c.Kill()
			_ = db.Close()
			// Stop waiting for result.
			close(initializationResult)
		}()

		if err := c.client.Run(c.ctx, func(ctx context.Context) error {
			status, err := c.client.Auth().Status(ctx)
			if err != nil {
				return err
			}
			if status.Authorized {
				// Ok.
			} else {
				if _, err := c.client.Auth().Bot(ctx, token.String()); err != nil {
					return err
				}
			}

			if err := c.api.Init(ctx); err != nil {
				return errors.Wrap(err, "init BotAPI")
			}
			defer func() {
				_ = gaps.Logout()
			}()

			// Done.
			select {
			case initializationResult <- nil:
				// Update lastUsed, because it is zero during initialization.
				c.Use(p.now())
			default:
			}

			<-ctx.Done()
			return ctx.Err()
		}); err != nil {
			// Failed.
			select {
			case initializationResult <- err:
				log.Warn("Initialize", zap.Error(err))
			default:
			}
		}
	}()

	return c, nil
}

type Options struct {
	AppID   int
	AppHash string
	Log     *zap.Logger
	Debug   bool
}

func NewPool(statePath string, opt Options) (*Pool, error) {
	p := &Pool{
		appID:     opt.AppID,
		appHash:   opt.AppHash,
		debug:     opt.Debug,
		log:       opt.Log,
		clients:   map[Token]*client{},
		statePath: statePath,
	}
	if err := os.MkdirAll(statePath, 0o750); err != nil {
		return nil, errors.Wrap(err, "create state dir")
	}
	return p, nil
}
