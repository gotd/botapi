package main

import (
	"sync/atomic"
	"time"

	glog "github.com/gotd/log"

	"github.com/gotd/botapi"
)

// handled counts updates that reached a handler — a tiny example of state shared
// across handler invocations from a custom middleware.
var handled atomic.Int64

// metrics is a custom Middleware: it wraps every handler, increments a counter
// and records how long the handler ran. A Middleware is just
// func(next Handler) Handler, so anything composable goes here — auth gates,
// per-user rate limits, request tracing, etc.
func metrics() botapi.Middleware {
	return func(next botapi.Handler) botapi.Handler {
		return func(c *botapi.Context) error {
			start := time.Now()
			n := handled.Add(1)

			err := next(c)

			glog.For(c.Bot.Logger()).Debug(c,
				"handled update",
				glog.Int64("seq", n),
				glog.Duration("took", time.Since(start)),
			)

			return err
		}
	}
}
