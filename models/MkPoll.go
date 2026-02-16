package models

import "time"

type MkPoll struct {
	Multiple  bool           `json:"multiple"`
	ExpiresAt *string        `json:"expiresAt"`
	Choices   []MkPollChoice `json:"choices"`
}

type MkPollChoice struct {
	Text    string `json:"text"`
	Votes   int    `json:"votes"`
	IsVoted bool   `json:"isVoted"`
}

func (p *MkPoll) ToPoll(noteID string) *Poll {
	if p == nil {
		return nil
	}
	poll := &Poll{
		ID:       noteID,
		Multiple: p.Multiple,
		Options:  []PollOption{},
		OwnVotes: []int{},
		Emojis:   []struct{}{},
	}
	if p.ExpiresAt != nil {
		poll.ExpiresAt = p.ExpiresAt
		t, err := time.Parse(time.RFC3339, *p.ExpiresAt)
		if err == nil {
			poll.Expired = t.Before(time.Now())
		}
	}
	var totalVotes int
	var hasVoted bool
	for i, c := range p.Choices {
		poll.Options = append(poll.Options, PollOption{
			Title:      c.Text,
			VotesCount: c.Votes,
		})
		totalVotes += c.Votes
		if c.IsVoted {
			hasVoted = true
			poll.OwnVotes = append(poll.OwnVotes, i)
		}
	}
	poll.VotesCount = totalVotes
	poll.Voted = &hasVoted
	return poll
}
