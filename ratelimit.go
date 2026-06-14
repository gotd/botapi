package botapi

import (
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/telegram"
	updhook "github.com/gotd/td/telegram/updates/hook"
	"golang.org/x/time/rate"
)

// buildMiddlewares assembles the client invoker middleware chain from the
// options. Order is outer-to-inner: optional flood-wait retry wraps the
// optional rate limiter, which wraps the update hook that harvests access
// hashes from RPC responses.
func buildMiddlewares(opt Options, h telegram.UpdateHandler) []telegram.Middleware {
	var mw []telegram.Middleware

	if opt.FloodWait {
		w := floodwait.NewSimpleWaiter()
		if opt.MaxFloodWaitRetries > 0 {
			w = w.WithMaxRetries(uint(opt.MaxFloodWaitRetries))
		}

		mw = append(mw, w)
	}

	if opt.RequestsPerSecond > 0 {
		burst := opt.RequestBurst
		if burst <= 0 {
			burst = 1
		}

		mw = append(mw, ratelimit.New(rate.Limit(opt.RequestsPerSecond), burst))
	}

	mw = append(mw, updhook.UpdateHook(h.Handle))

	return mw
}
