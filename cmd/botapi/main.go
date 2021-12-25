// Binary botapi implements telegram-bot-api server using gotd.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-faster/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/gotd/td/constant"

	"github.com/gotd/botapi/internal/botapi"
	"github.com/gotd/botapi/internal/oas"
	"github.com/gotd/botapi/internal/pool"
)

func listen(ctx context.Context, addr string, h http.Handler, logger *zap.Logger) error {
	grp, ctx := errgroup.WithContext(ctx)

	listenCfg := net.ListenConfig{}
	l, err := listenCfg.Listen(ctx, "tcp", addr)
	if err != nil {
		return errors.Errorf("bind %q: %w", addr, err)
	}
	logger.Info("Listen",
		zap.String("addr", addr),
	)

	srv := &http.Server{
		Addr:    addr,
		Handler: h,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	grp.Go(func() error {
		<-ctx.Done()

		// TODO: make it configurable
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return errors.Errorf("shutdown: %w", err)
		}

		return nil
	})
	grp.Go(func() error {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Errorf("serve %q: %w", l.Addr(), err)
		}
		return nil
	})

	if err := grp.Wait(); err != nil {
		return fmt.Errorf("http: %w", err)
	}

	return nil
}

func run(ctx context.Context) error {
	var (
		appID     = flag.Int("api-id", constant.TestAppID, "The api_id of application")
		appHash   = flag.String("api-hash", constant.TestAppHash, "The api_hash of application")
		addr      = flag.String("addr", "localhost:8081", "http listen addr")
		keepalive = flag.Duration("keepalive", 5*time.Minute, "client keepalive")
		statePath = flag.String("state", "state", "path to state file (json)")
		debug     = flag.Bool("debug", false, "enables debug mode")
	)
	flag.Parse()

	level := zap.InfoLevel
	if *debug {
		level = zap.DebugLevel
	}
	log, err := zap.NewDevelopment(zap.IncreaseLevel(level))
	if err != nil {
		return errors.Errorf("create logger: %w", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Creating pool",
		zap.Duration("keep_alive", *keepalive),
		zap.String("storage", *statePath),
		zap.Bool("debug", *debug),
	)
	p, err := pool.NewPool(*statePath, pool.Options{
		AppID:   *appID,
		AppHash: *appHash,
		Log:     log.Named("pool"),
		Debug:   *debug,
	})
	if err != nil {
		panic(err)
	}
	go p.RunGC(*keepalive)

	// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
	r := chi.NewRouter()
	r.Post("/bot{token}/{method}", func(w http.ResponseWriter, r *http.Request) {
		token, err := pool.ParseToken(chi.URLParam(r, "token"))
		if err != nil {
			botapi.NotFound(w, r)
			return
		}
		method := chi.URLParam(r, "method")

		log := log.With(zap.Int("bot_id", token.ID), zap.String("method", method))

		log.Info("New request")
		if err := p.Do(r.Context(), token, func(client *botapi.BotAPI) error {
			r.URL.Path = botapi.CorrectMethod(method)
			oas.NewServer(client).ServeHTTP(w, r)
			return nil
		}); err != nil {
			log.Warn("Do error", zap.Error(err))
		}
	})

	return listen(ctx, *addr, r, log.Named("http"))
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(2)
	}
}
