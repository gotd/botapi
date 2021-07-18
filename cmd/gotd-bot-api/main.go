// Binary gotd-bot-api implements telegram-bot-api server using gotd.
package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram"

	"github.com/gotd/botapi/pool"
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
	p, err := pool.NewPool(pool.Options{
		AppID:   *appID,
		AppHash: *appHash,
		Log:     log.Named("pool"),
	})
	if err != nil {
		panic(err)
	}

	// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
	r := chi.NewRouter()
	r.Route("/bot{token}", func(r chi.Router) {
		r.Post("/getMe", func(w http.ResponseWriter, r *http.Request) {
			token, err := pool.ParseToken(chi.URLParam(r, "token"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			if err := p.Do(ctx, token, func(client *telegram.Client) error {
				res, err := client.Self(ctx)
				if err != nil {
					return err
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

				return nil
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})

	if err := http.ListenAndServe(*addr, r); err != nil {
		panic(err)
	}
}
