package pool

import (
	"context"
	"sync"
	"time"

	"github.com/gotd/td/telegram"

	"github.com/gotd/botapi/internal/botapi"
)

type client struct {
	ctx    context.Context
	cancel context.CancelFunc

	mux      sync.Mutex
	api      *botapi.BotAPI
	client   *telegram.Client
	token    Token
	lastUsed time.Time
}

func (c *client) Kill() {
	c.cancel()
}

func (c *client) Deadline(deadline time.Time) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	return !c.lastUsed.IsZero() && c.lastUsed.Before(deadline)
}

func (c *client) Use(t time.Time) {
	c.mux.Lock()
	c.lastUsed = t
	c.mux.Unlock()
}
