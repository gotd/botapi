// Package pool implements client pool.
package pool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
)

type fileStorage struct {
	path string
	mux  sync.Mutex
}

type sessionFile struct {
	Data map[int][]byte `json:"data"`
}

func (s *fileStorage) Store(ctx context.Context, id int, data []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	var decoded sessionFile

	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) || len(b) == 0 {
		// Blank initial session.
	} else if err == nil {
		if err := json.Unmarshal(b, &decoded); err != nil {
			return xerrors.Errorf("unmarshal session file: %w", err)
		}
	}
	if decoded.Data == nil {
		decoded.Data = map[int][]byte{}
	}

	decoded.Data[id] = data

	if b, err = json.Marshal(&decoded); err != nil {
		return err
	}

	return os.WriteFile(s.path, b, 0600)
}

func (s *fileStorage) Load(ctx context.Context, id int) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) || len(data) == 0 {
		return nil, session.ErrNotFound
	}

	var decoded sessionFile
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}

	if len(decoded.Data) == 0 {
		return nil, session.ErrNotFound
	}

	return decoded.Data[id], nil
}

type clientStorage struct {
	storage StateStorage
	id      int
}

func (c clientStorage) LoadSession(ctx context.Context) ([]byte, error) {
	data, err := c.storage.Load(ctx, c.id)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, session.ErrNotFound
	}
	return data, nil
}

func (c clientStorage) StoreSession(ctx context.Context, data []byte) error {
	return c.storage.Store(ctx, c.id, data)
}

// Token represents bot token, like 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
type Token struct {
	ID     int    // 123456
	Secret string // ABC-DEF1234ghIkl-zyx57W2v1u123ew11
}

func ParseToken(s string) (Token, error) {
	if s == "" {
		return Token{}, errors.New("blank")
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return Token{}, errors.New("invalid token")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return Token{}, err
	}
	return Token{
		ID:     id,
		Secret: parts[1],
	}, err
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%s", t.ID, t.Secret)
}

type client struct {
	ctx    context.Context
	cancel context.CancelFunc

	mux      sync.Mutex
	telegram *telegram.Client
	token    Token
	lastUsed time.Time
}

func (c *client) Deadline(deadline time.Time) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lastUsed.Before(deadline)
}

func (c *client) Use(t time.Time) {
	c.mux.Lock()
	c.lastUsed = t
	c.mux.Unlock()
}

// Pool of clients.
type Pool struct {
	appID   int
	appHash string
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
			c.cancel()
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

// Do acquires telegram client by token.
//
// Returns error if token is invalid. Block until client is available,
// authentication error or context cancelled.
func (p *Pool) Do(ctx context.Context, token Token, fn func(client *telegram.Client) error) error {
	p.clientsMux.Lock()
	c, ok := p.clients[token]
	p.clientsMux.Unlock()

	if ok {
		// Happy path.
		c.Use(p.now())
		return fn(c.telegram)
	}

	options := telegram.Options{
		Logger: p.log.Named("client").With(zap.Int("id", token.ID)),
	}
	if p.storage != nil {
		options.SessionStorage = clientStorage{
			id:      token.ID,
			storage: p.storage,
		}
	}
	tgClient := telegram.NewClient(p.appID, p.appHash, options)

	tgContext, tgCancel := context.WithCancel(context.Background())
	c = &client{
		ctx:      tgContext,
		cancel:   tgCancel,
		telegram: tgClient,
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
			if ok && c.telegram == tgClient {
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

		return fn(c.telegram)
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
}

type StateStorage interface {
	Store(ctx context.Context, id int, data []byte) error
	Load(ctx context.Context, id int) ([]byte, error)
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
		log:     opt.Log,
		clients: map[Token]*client{},
		storage: opt.Storage,
	}
	return p, nil
}
