package botapi

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/gotd/log"
)

// Recover wraps a handler so a panic is recovered, logged with its stack, and
// converted into an error instead of crashing the update loop.
func Recover() Middleware {
	return func(next Handler) Handler {
		return func(c *Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					c.Bot.logger().Error(c, "Recovered from panic in handler",
						log.Any("panic", r),
						log.String("stack", string(debug.Stack())),
					)
					err = &Error{Code: 500, Description: "Internal Server Error: handler panicked"}
				}
			}()
			return next(c)
		}
	}
}

// Timeout wraps a handler so its context is canceled after d.
func Timeout(d time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(c *Context) error {
			ctx, cancel := context.WithTimeout(c.Context, d)
			defer cancel()
			scoped := *c
			scoped.Context = ctx
			return next(&scoped)
		}
	}
}

// Logging wraps a handler to log each handled update at debug level, and at warn
// level when the handler returns an error.
func Logging() Middleware {
	return func(next Handler) Handler {
		return func(c *Context) error {
			err := next(c)
			if err != nil {
				c.Bot.logger().Warn(c, "Update handled with error",
					log.Int("update_id", c.Update.UpdateID), log.Error(err))
			} else {
				c.Bot.logger().Debug(c, "Update handled", log.Int("update_id", c.Update.UpdateID))
			}
			return err
		}
	}
}
