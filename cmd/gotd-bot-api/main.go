// Binary gotd-bot-api implements telegram-bot-api server using gotd.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
)

func main() {
	var (
		appID   = flag.Int("api-id", 0, "The api_id of application")
		appHash = flag.String("api-hash", "", "The api_hash of application")
		addr    = flag.String("addr", "localhost:8081", "http listen addr")
	)
	flag.Parse()

	log, err := zap.NewDevelopment(zap.IncreaseLevel(zap.InfoLevel))
	if err != nil {
		panic(err)
	}

	log.Info("Start", zap.String("addr", *addr))

	var mux sync.Mutex
	clients := make(map[string]*telegram.Client)

	// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
	r := chi.NewRouter()
	r.Route("/bot{token}", func(r chi.Router) {
		r.Post("/getMe", func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "token")
			if token == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Info("Connect", zap.String("token", token))

			mux.Lock()

			client, ok := clients[token]
			if !ok {
				client = telegram.NewClient(*appID, *appHash, telegram.Options{
					Logger: log,
				})
				clients[token] = client
				initResult := make(chan error, 1)

				go func() {
					ctx := context.Background()
					if err := client.Run(ctx, func(ctx context.Context) error {
						if _, err := client.Auth().Bot(ctx, token); err != nil {
							return err
						}

						log.Info("Logged in")

						initResult <- nil

						<-ctx.Done()
						return ctx.Err()
					}); err != nil {
						log.Error("Run failed", zap.Error(err))

						select {
						case initResult <- err:
						default:
						}
					}
				}()

				<-initResult
			}

			mux.Unlock()

			res, err := client.Self(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			type User struct {
				ID int `json:"id"`

				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`

				// Returns only in getMe
				CanJoinGroups   bool `json:"can_join_groups"`
				CanReadMessages bool `json:"can_read_all_group_messages"`
				SupportsInline  bool `json:"supports_inline_queries"`
			}
			_ = json.NewEncoder(w).Encode(struct {
				Result User `json:"result"`
			}{
				Result: User{
					ID:              res.ID,
					FirstName:       res.FirstName,
					LastName:        res.LastName,
					Username:        res.Username,
					LanguageCode:    res.LangCode,
					IsBot:           res.Bot,
					CanJoinGroups:   !res.BotNochats,
					CanReadMessages: res.BotChatHistory,
					SupportsInline:  false, // ?
				},
			})
		})
	})

	if err := http.ListenAndServe(*addr, r); err != nil {
		panic(err)
	}
}
