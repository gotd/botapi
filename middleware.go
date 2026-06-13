package botapi

import (
	"context"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

// Recover wraps a handler so a panic is recovered, logged with its stack, and
// converted into an error instead of crashing the update loop.
func Recover() Middleware {
	return func(next Handler) Handler {
		return func(c *Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					c.Bot.log.Error("Recovered from panic in handler",
						zap.Any("panic", r),
						zap.ByteString("stack", debug.Stack()),
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
				c.Bot.log.Warn("Update handled with error",
					zap.Int("update_id", c.Update.UpdateID), zap.Error(err))
			} else {
				c.Bot.log.Debug("Update handled", zap.Int("update_id", c.Update.UpdateID))
			}
			return err
		}
	}
}
