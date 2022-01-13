package botapi

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

var _pollOptions = []byte(`0123456789`)

func pollOption(i int) []byte {
	if i < len(_pollOptions) {
		return _pollOptions[i : i+1]
	}
	return []byte{byte(i)}
}

func (b *BotAPI) convertToBotAPIPoll(ctx context.Context, media *tg.MessageMediaPoll) (oas.Poll, bool) {
	var (
		poll    = media.Poll
		results = media.Results

		typ = oas.PollTypeRegular
	)
	if a, r := len(poll.Answers), len(results.Results); a != r {
		b.logger.Warn("Got poll where len(answers) != len(results)",
			zap.Int64("poll_id", poll.ID),
			zap.Int("answers", a),
			zap.Int("results", r),
		)
		return oas.Poll{}, false
	}

	if poll.Quiz {
		typ = oas.PollTypeQuiz
	}
	resultPoll := oas.Poll{
		ID:                    strconv.FormatInt(poll.ID, 10),
		Question:              poll.Question,
		Options:               nil,
		TotalVoterCount:       results.TotalVoters,
		IsClosed:              poll.Closed,
		IsAnonymous:           !poll.PublicVoters,
		Type:                  typ,
		AllowsMultipleAnswers: poll.MultipleChoice,
		CorrectOptionID:       oas.OptInt{},
		Explanation:           optString(results.GetSolution),
		ExplanationEntities:   nil,
		OpenPeriod:            optInt(poll.GetClosePeriod),
		CloseDate:             optInt(poll.GetCloseDate),
	}

	if e := results.SolutionEntities; len(e) > 0 {
		resultPoll.ExplanationEntities = b.convertToBotAPIEntities(ctx, e)
	}

	// SAFETY: length equality checked above.
	for i, result := range results.Results {
		if result.Correct {
			resultPoll.CorrectOptionID.SetTo(i)
		}
		resultPoll.Options = append(resultPoll.Options, oas.PollOption{
			Text:       poll.Answers[i].Text,
			VoterCount: result.Voters,
		})
	}

	return resultPoll, true
}

// SendPoll implements oas.Handler.
func (b *BotAPI) SendPoll(ctx context.Context, req oas.SendPoll) (oas.ResultMessage, error) {
	s, p, err := b.prepareSend(
		ctx,
		sendOpts{
			To:                       req.ChatID,
			DisableNotification:      req.DisableNotification,
			ProtectContent:           req.ProtectContent,
			ReplyToMessageID:         req.ReplyToMessageID,
			AllowSendingWithoutReply: req.AllowSendingWithoutReply,
			ReplyMarkup:              req.ReplyMarkup,
		},
	)
	if err != nil {
		return oas.ResultMessage{}, errors.Wrap(err, "prepare send")
	}

	answers := make([]tg.PollAnswer, len(req.Options))
	for i, opt := range req.Options {
		answers[i] = tg.PollAnswer{
			Text:   opt,
			Option: pollOption(i),
		}
	}

	poll := tg.Poll{
		Closed:         req.IsClosed.Value,
		PublicVoters:   !req.IsAnonymous.Value,
		MultipleChoice: req.AllowsMultipleAnswers.Value,
		Quiz:           req.Type.Value == "quiz",
		Question:       req.Question,
		Answers:        answers,
		ClosePeriod:    0,
		CloseDate:      0,
	}
	poll.SetFlags()

	if v, ok := req.OpenPeriod.Get(); ok {
		// Prefer open_period.
		//
		// See https://github.com/tdlib/td/blob/fa8feefed70d64271945e9d5fd010b957d93c8cd/td/telegram/MessageContent.cpp#L1914-L1916.
		poll.SetClosePeriod(v)
	} else if v, ok := req.CloseDate.Get(); ok {
		poll.SetCloseDate(v)
	}

	media := &tg.InputMediaPoll{
		Poll:             poll,
		CorrectAnswers:   nil,
		Solution:         "",
		SolutionEntities: nil,
	}
	if v, ok := req.CorrectOptionID.Get(); ok {
		if v < 0 || v >= len(answers) {
			// See https://github.com/tdlib/td/blob/fa8feefed70d64271945e9d5fd010b957d93c8cd/td/telegram/MessageContent.cpp#L1898.
			return oas.ResultMessage{}, &BadRequestError{Message: "Wrong correct option ID specified"}
		}
		media.CorrectAnswers = [][]byte{answers[v].Option}
	}
	if explanation, ok := req.Explanation.Get(); ok {
		// FIXME(tdakkota): get entities from request.
		parseMode, isParseModeSet := req.ExplanationParseMode.Get()
		if isParseModeSet && parseMode != "HTML" {
			return oas.ResultMessage{}, &NotImplementedError{Message: "only HTML formatting is supported"}
		}

		if isParseModeSet {
			var builder entity.Builder
			if err := html.HTML(strings.NewReader(explanation), &builder, html.Options{
				UserResolver:          b.peers.UserResolveHook(ctx),
				DisableTelegramEscape: false,
			}); err != nil {
				return oas.ResultMessage{}, errors.Wrap(err, "parse explanation")
			}
			text, entities := builder.Complete()
			media.SetSolution(text)
			media.SetSolutionEntities(entities)
		} else {
			media.SetSolution(explanation)
		}
	}

	resp, err := s.Media(ctx, message.Media(media))
	return b.sentMessage(ctx, p, resp, err)
}

// StopPoll implements oas.Handler.
func (b *BotAPI) StopPoll(ctx context.Context, req oas.StopPoll) (oas.ResultPoll, error) {
	return oas.ResultPoll{}, &NotImplementedError{}
}
