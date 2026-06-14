package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestPollFromTg(t *testing.T) {
	poll := &tg.Poll{
		ID:             99,
		Closed:         true,
		Quiz:           true,
		MultipleChoice: false,
		Question:       tg.TextWithEntities{Text: "2+2?"},
		ClosePeriod:    30,
		Answers: []tg.PollAnswerClass{
			&tg.PollAnswer{Text: tg.TextWithEntities{Text: "3"}, Option: []byte{0}},
			&tg.PollAnswer{Text: tg.TextWithEntities{Text: "4"}, Option: []byte{1}},
		},
	}
	poll.SetPublicVoters(false)

	results := &tg.PollResults{
		TotalVoters: 5,
		Results: []tg.PollAnswerVoters{
			{Option: []byte{0}, Voters: 2, Correct: false},
			{Option: []byte{1}, Voters: 3, Correct: true},
		},
	}
	results.SetSolution("it's four")

	got := pollFromTg(poll, results)
	if got.ID != "99" || got.Question != "2+2?" || !got.IsClosed || got.Type != PollQuiz {
		t.Fatalf("poll header: %#v", got)
	}

	if !got.IsAnonymous {
		t.Fatal("non-public poll should be anonymous")
	}

	if len(got.Options) != 2 || got.Options[0].Text != "3" || got.Options[1].VoterCount != 3 {
		t.Fatalf("options: %#v", got.Options)
	}

	if got.CorrectOptionID != 1 {
		t.Fatalf("correct option: %d", got.CorrectOptionID)
	}

	if got.TotalVoterCount != 5 || got.Explanation != "it's four" || got.OpenPeriod != 30 {
		t.Fatalf("tally/explanation: %#v", got)
	}
}
