package main

import (
	"strings"

	glog "github.com/gotd/log"

	"github.com/gotd/botapi"
)

// registerInline wires inline mode (enable it in @BotFather first). Type
// "@yourbot some text" in any chat and the bot offers transformed variants;
// picking one fires OnChosenInlineResult.
func registerInline(bot *botapi.Bot) {
	bot.OnInlineQuery(func(c *botapi.Context) error {
		q := strings.TrimSpace(c.Update.InlineQuery.Query)
		if q == "" {
			// An empty answer clears the inline results list.
			return c.AnswerInline(nil)
		}

		results := []botapi.InlineQueryResult{
			article("upper", "UPPERCASE", strings.ToUpper(q)),
			article("lower", "lowercase", strings.ToLower(q)),
			article("reverse", "Reversed", reverse(q)),
		}

		// IsPersonal + a short cache time keep results per-user and fresh.
		return c.AnswerInline(results,
			botapi.WithInlineCacheTime(1),
			botapi.WithInlinePersonal(),
		)
	})

	bot.OnChosenInlineResult(func(c *botapi.Context) error {
		chosen := c.Update.ChosenInlineResult
		glog.For(c.Bot.Logger()).Info(c, "inline result chosen",
			glog.String("result_id", chosen.ResultID),
			glog.String("query", chosen.Query),
		)

		return nil
	})
}

// article builds a text-content inline result with an attached inline keyboard.
func article(id, title, text string) botapi.InlineQueryResult {
	return &botapi.InlineQueryResultArticle{
		ID:                  id,
		Title:               title,
		Description:         text,
		InputMessageContent: &botapi.InputTextMessageContent{MessageText: text},
		ReplyMarkup: botapi.InlineKeyboard(botapi.InlineRow(
			botapi.InlineButtonURL("gotd/td", "https://github.com/gotd/td"),
		)),
	}
}
