// Binary botapi implements telegram-bot-api server using gotd.
package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/gotd/botapi/internal/botapi"
	"github.com/gotd/botapi/internal/oas"
	"github.com/gotd/botapi/internal/pool"
)

func main() {
	var (
		appID     = flag.Int("api-id", 0, "The api_id of application")
		appHash   = flag.String("api-hash", "", "The api_hash of application")
		addr      = flag.String("addr", "localhost:8081", "http listen addr")
		keepalive = flag.Duration("keepalive", time.Second*5, "client keepalive")
		statePath = flag.String("state", "", "path to state file (json)")
		debug     = flag.Bool("debug", false, "enables debug mode")
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

	handler := botapi.NewBotAPI(p, *debug)
	server := oas.NewServer(handler)

	// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
	r := chi.NewRouter()
	r.Post("/bot{token}/{method}", func(w http.ResponseWriter, r *http.Request) {
		token, err := pool.ParseToken(chi.URLParam(r, "token"))
		if err != nil {
			botapi.NotFound(w, r)
			return
		}
		r.WithContext(botapi.PropagateToken(r.Context(), token))

		r.URL.Path = chi.URLParam(r, "method")
		server.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(*addr, r); err != nil {
		panic(err)
	}
}
