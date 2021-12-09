// Binary botapi implements telegram-bot-api server using gotd.
package main

import (
	"context"
	"flag"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"

	"github.com/gotd/botapi/internal/pool"
)

type handleContext struct {
	Method  string
	Client  *telegram.Client
	Writer  http.ResponseWriter
	Request *http.Request
}

type handler struct {
	handlers map[string]func(ctx context.Context, h handleContext) error
}

func (h handler) On(method string, f func(ctx context.Context, h handleContext) error) {
	h.handlers[strings.ToLower(method)] = f
}

func main() {
	var (
		appID     = flag.Int("api-id", 0, "The api_id of application")
		appHash   = flag.String("api-hash", "", "The api_hash of application")
		addr      = flag.String("addr", "localhost:8081", "http listen addr")
		keepalive = flag.Duration("keepalive", time.Second*5, "client keepalive")
		statePath = flag.String("state", "", "path to state file (json)")
	)
	flag.Parse()

	log, err := zap.NewDevelopment(zap.IncreaseLevel(zap.InfoLevel))
	if err != nil {
		panic(err)
	}

	var storage pool.StateStorage
	if *statePath != "" {
		storage = pool.NewFileStorage(*statePath)
	}

	log.Info("Start", zap.String("addr", *addr))
	p, err := pool.NewPool(pool.Options{
		AppID:   *appID,
		AppHash: *appHash,
		Log:     log.Named("pool"),
		Storage: storage,
	})
	if err != nil {
		panic(err)
	}
	go p.RunGC(*keepalive)

	h := handler{
		handlers: map[string]func(ctx context.Context, h handleContext) error{},
	}
	h.On("getMe", getMe)

	// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
	r := chi.NewRouter()
	r.Post("/bot{token}/{method}", func(w http.ResponseWriter, r *http.Request) {
		token, err := pool.ParseToken(chi.URLParam(r, "token"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		method := strings.ToLower(chi.URLParam(r, "method"))
		handler, ok := h.handlers[method]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := r.Context()
		if err := p.Do(ctx, token, func(client *telegram.Client) error {
			return handler(ctx, handleContext{
				Method:  method,
				Client:  client,
				Writer:  w,
				Request: r,
			})
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(*addr, r); err != nil {
		panic(err)
	}
}
