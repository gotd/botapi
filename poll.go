package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// pollFromTg converts an MTProto poll and its results into a Bot API Poll.
func pollFromTg(poll *tg.Poll, results *tg.PollResults) *Poll {
	out := &Poll{
		ID:                    strconv.FormatInt(poll.ID, 10),
		Question:              poll.Question.Text,
		IsClosed:              poll.Closed,
		IsAnonymous:           !poll.PublicVoters,
		Type:                  PollRegular,
		AllowsMultipleAnswers: poll.MultipleChoice,
		OpenPeriod:            poll.ClosePeriod,
	}
	if poll.Quiz {
		out.Type = PollQuiz
	}

	// Index vote counts and the correct option by answer bytes.
	type stat struct {
		voters  int
		correct bool
	}

	byOption := make(map[string]stat, len(results.Results))

	for _, r := range results.Results {
		byOption[string(r.Option)] = stat{voters: r.Voters, correct: r.Correct}
	}

	out.Options = make([]PollOption, 0, len(poll.Answers))
	for i, a := range poll.Answers {
		ans, ok := a.(*tg.PollAnswer)
		if !ok {
			continue
		}

		s := byOption[string(ans.Option)]

		out.Options = append(out.Options, PollOption{Text: ans.Text.Text, VoterCount: s.voters})

		if s.correct {
			out.CorrectOptionID = i
		}
	}

	out.TotalVoterCount = results.TotalVoters
	if s, ok := results.GetSolution(); ok {
		out.Explanation = s
	}

	out.ExplanationEntities = entitiesFromTg(results.SolutionEntities)

	return out
}

// StopPoll stops an open poll and returns its final state. The bot must be the
// poll's author.
func (b *Bot) StopPoll(ctx context.Context, chat ChatID, messageID int, markup ReplyMarkup) (*Poll, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	poll, results, err := b.fetchPoll(ctx, p, messageID)
	if err != nil {
		return nil, err
	}

	closed := *poll

	closed.Closed = true

	req := &tg.MessagesEditMessageRequest{
		Peer: p.InputPeer(),
		ID:   messageID,
	}
	req.SetMedia(&tg.InputMediaPoll{Poll: closed})

	if markup != nil {
		mkp, err := replyMarkupToTg(markup)
		if err != nil {
			return nil, err
		}

		req.SetReplyMarkup(mkp)
	}

	if _, err := b.raw.MessagesEditMessage(ctx, req); err != nil {
		return nil, asAPIError(err)
	}

	// Closing does not change the tally, so the fetched results are final.
	return pollFromTg(&closed, results), nil
}

// fetchPoll loads the poll media of a message via the appropriate getMessages
// call for the peer kind.
func (b *Bot) fetchPoll(ctx context.Context, p peers.Peer, messageID int) (*tg.Poll, *tg.PollResults, error) {
	ids := []tg.InputMessageClass{&tg.InputMessageID{ID: messageID}}

	var (
		res tg.MessagesMessagesClass
		err error
	)

	if ch, ok := p.(peers.Channel); ok {
		res, err = b.raw.ChannelsGetMessages(ctx, &tg.ChannelsGetMessagesRequest{
			Channel: ch.InputChannel(),
			ID:      ids,
		})
	} else {
		res, err = b.raw.MessagesGetMessages(ctx, ids)
	}

	if err != nil {
		return nil, nil, asAPIError(err)
	}

	msgs, ok := res.AsModified()
	if !ok {
		return nil, nil, &Error{Code: 400, Description: "Bad Request: message to stop not found"}
	}

	for _, m := range msgs.GetMessages() {
		msg, ok := m.(*tg.Message)
		if !ok || msg.ID != messageID {
			continue
		}

		media, ok := msg.Media.(*tg.MessageMediaPoll)
		if !ok {
			return nil, nil, &Error{Code: 400, Description: "Bad Request: message is not a poll"}
		}

		return &media.Poll, &media.Results, nil
	}

	return nil, nil, &Error{Code: 400, Description: "Bad Request: message to stop not found"}
}
