package main

import (
	"flag"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	token := flag.String("token", "", "bot token")
	flag.Parse()

	b, err := tb.NewBot(tb.Settings{
		URL: "http://localhost:8081",

		Token:  *token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	_ = b
}
