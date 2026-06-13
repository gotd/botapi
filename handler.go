package botapi

import (
	"context"
	"sync"

	"github.com/gotd/log"
)

// Context is passed to a Handler. It embeds the request context (so it can be
// passed straight to Bot methods) and carries the Bot and the Update.
type Context struct {
	context.Context

	Bot    *Bot
	Update *Update
}

// Handler processes a single update.
type Handler func(c *Context) error

// Predicate reports whether a Handler should run for an update. A Handler runs
// only if all of its predicates return true.
type Predicate func(u *Update) bool

// Middleware wraps a Handler, returning a new one. Middleware registered with
// Bot.Use runs for every handled update, outermost first.
type Middleware func(next Handler) Handler

type route struct {
	handler    Handler
	predicates []Predicate
	mws        []Middleware
}

func (r route) matches(u *Update) bool {
	for _, p := range r.predicates {
		if !p(u) {
			return false
		}
	}
	return true
}

// router holds the registered routes and global middleware. The zero value is
// ready to use and safe for concurrent registration and routing.
type router struct {
	mu     sync.RWMutex
	routes []route
	mws    []Middleware
}

// Use registers global middleware applied to every handled update. Middleware
// runs outermost-first in registration order. Call before Run.
func (b *Bot) Use(mws ...Middleware) {
	b.router.mu.Lock()
	defer b.router.mu.Unlock()
	b.router.mws = append(b.router.mws, mws...)
}

// on registers a handler guarded by the given predicates.
func (b *Bot) on(handler Handler, predicates ...Predicate) {
	b.onWith(handler, nil, predicates)
}

// onWith registers a handler with route-scoped middleware (applied inside the
// global middleware) and predicates.
func (b *Bot) onWith(handler Handler, mws []Middleware, predicates []Predicate) {
	b.router.mu.Lock()
	defer b.router.mu.Unlock()
	b.router.routes = append(b.router.routes, route{handler: handler, predicates: predicates, mws: mws})
}

// route dispatches an update to the first matching handler, wrapped in the
// global middleware. Handler errors are logged, not propagated to the update
// loop, so one failing handler does not tear down the bot.
func (b *Bot) route(ctx context.Context, u *Update) {
	b.router.mu.RLock()
	routes := b.router.routes
	mws := b.router.mws
	b.router.mu.RUnlock()

	for _, r := range routes {
		if !r.matches(u) {
			continue
		}
		h := r.handler
		for i := len(r.mws) - 1; i >= 0; i-- {
			h = r.mws[i](h)
		}
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}
		c := &Context{Context: ctx, Bot: b, Update: u}
		if err := h(c); err != nil {
			b.logger().Error(ctx, "Handler error", log.Error(err))
		}
		return
	}
}
