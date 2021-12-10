// Package pool implements client pool.
package pool

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/botapi"
	"github.com/gotd/botapi/internal/peers"
)

// Pool of clients.
type Pool struct {
	appID   int
	appHash string
	debug   bool
	log     *zap.Logger

	storage    StateStorage
	clientsMux sync.Mutex
	clients    map[Token]*client
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
	log := p.log.Named("client").With(zap.Int("id", token.ID))

	peerStorage := peers.NewInmemoryStorage()
	gaps := updates.New(updates.Config{
		Handler: telegram.UpdateHandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {
			return nil
		}),
		AccessHasher: peers.AccessHasher{
			Storage: peerStorage,
		},
		Logger: log.Named("gaps"),
	})
	options := telegram.Options{
		Logger:        log,
		UpdateHandler: peers.UpdateHook(peerStorage, gaps),
	}
	if p.storage != nil {
		options.SessionStorage = clientStorage{
			id:      fmt.Sprintf("%x:%x", token.ID, sha256.Sum256([]byte(token.Secret))),
			storage: p.storage,
		}
	}
	tgClient := telegram.NewClient(p.appID, p.appHash, options)

	tgContext, tgCancel := context.WithCancel(context.Background())
	c = &client{
		ctx:      tgContext,
		cancel:   tgCancel,
		api:      botapi.NewBotAPI(tgClient, peerStorage, p.debug),
		token:    token,
		lastUsed: p.now(),
	}

	// Wait for initialization.
	initializationResult := make(chan error, 1)
	go func() {
		defer close(initializationResult)
		defer tgCancel()

		defer func() {
			// Removing client from client list on close.
			p.clientsMux.Lock()
			c, ok := p.clients[token]
			if ok && c.api.Client() == tgClient {
				delete(p.clients, token)
			}
			p.clientsMux.Unlock()
		}()

		if err := tgClient.Run(c.ctx, func(ctx context.Context) error {
			status, err := tgClient.Auth().Status(ctx)
			if err != nil {
				return err
			}
			if status.Authorized {
				// Ok.
			} else {
				if _, err := tgClient.Auth().Bot(ctx, token.String()); err != nil {
					return err
				}
			}

			// Done.
			select {
			case initializationResult <- nil:
			default:
			}

			<-ctx.Done()
			return ctx.Err()
		}); err != nil {
			// Failed.
			select {
			case initializationResult <- err:
			default:
			}
		}
	}()

	// Waiting for initialization.
	select {
	case err := <-initializationResult:
		if err != nil {
			return err
		}

		p.clientsMux.Lock()
		conflictingClient, ok := p.clients[token]
		if ok {
			// Existing conflicting client, so stopping current client.
			tgCancel()
			c = conflictingClient
		} else {
			p.clients[token] = c
		}
		p.clientsMux.Unlock()

		return fn(c.api)
	case <-ctx.Done():
		return ctx.Err()
	case <-tgContext.Done():
		return tgContext.Err()
	}
}

func (p *Pool) RunGC(timeout time.Duration) {
	timer := time.NewTicker(time.Second)
	for now := range timer.C {
		deadline := now.Add(-timeout)
		p.tick(deadline)
	}
}

type Options struct {
	AppID   int
	AppHash string
	Log     *zap.Logger
	Storage StateStorage
	Debug   bool
}

type StateStorage interface {
	Store(ctx context.Context, id string, data []byte) error
	Load(ctx context.Context, id string) ([]byte, error)
}

func NewFileStorage(path string) StateStorage {
	return &fileStorage{
		path: path,
	}
}

func NewPool(opt Options) (*Pool, error) {
	p := &Pool{
		appID:   opt.AppID,
		appHash: opt.AppHash,
		debug:   opt.Debug,
		log:     opt.Log,
		clients: map[Token]*client{},
		storage: opt.Storage,
	}
	return p, nil
}
