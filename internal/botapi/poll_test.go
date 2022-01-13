package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func Test_pollOption(t *testing.T) {
	require.Equal(t, []byte{'0'}, pollOption(0))
}

func TestBotAPI_SendPoll(t *testing.T) {
	ctx := context.Background()
	testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
		poll := tg.Poll{
			Closed:         true,
			PublicVoters:   false,
			MultipleChoice: true,
			Quiz:           true,
			Question:       "question",
			Answers: []tg.PollAnswer{
				{
					Option: pollOption(0),
					Text:   "0",
				},
				{
					Option: pollOption(1),
					Text:   "1",
				},
			},
			ClosePeriod: 10,
			CloseDate:   0,
		}
		poll.SetFlags()
		testSentMedia(a, mock, &tg.InputMediaPoll{
			Poll: poll,
			CorrectAnswers: [][]byte{
				pollOption(0),
			},
			Solution: "solution",
			SolutionEntities: []tg.MessageEntityClass{
				&tg.MessageEntityBold{
					Offset: 0,
					Length: len("solution"),
				},
			},
		})

		_, err := api.SendPoll(ctx, oas.SendPoll{
			ChatID:                   oas.NewInt64ID(testChatID()),
			Question:                 "question",
			Options:                  []string{"0", "1"},
			IsAnonymous:              oas.NewOptBool(true),
			Type:                     oas.NewOptString("quiz"),
			AllowsMultipleAnswers:    oas.NewOptBool(true),
			CorrectOptionID:          oas.NewOptInt(0),
			Explanation:              oas.NewOptString(`<b>solution</b>`),
			ExplanationParseMode:     oas.NewOptString(`HTML`),
			OpenPeriod:               oas.NewOptInt(10),
			CloseDate:                oas.NewOptInt(1337),
			IsClosed:                 oas.NewOptBool(true),
			DisableNotification:      oas.OptBool{},
			ProtectContent:           oas.OptBool{},
			ReplyToMessageID:         oas.OptInt{},
			AllowSendingWithoutReply: oas.OptBool{},
			ReplyMarkup:              oas.OptSendReplyMarkup{},
		})
		a.NoError(err)
	})
}
