package discord

import "time"

// https://docs.discord.com/developers/resources/poll#layout-type
type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = 1
)

// https://docs.discord.com/developers/resources/poll#poll-media-object
type PollQuestion struct {
	// Text should always be non-null for both questions and answers, but please do not depend on that in the future ~API Docs
	Text  string `json:"text,omitempty"`
	Emoji *Emoji `json:"emoji,omitempty"`
}

// https://docs.discord.com/developers/resources/poll#poll-answer-object
type PollAnswer struct {
	AnswerID  *int         `json:"answer_id,omitempty"`
	PollMedia PollQuestion `json:"poll_media"`
}

// https://docs.discord.com/developers/resources/poll#poll-results-object-poll-answer-count-object-structure
type PollResultsAnswerCounts struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

// https://docs.discord.com/developers/resources/poll#poll-results-object
type PollResults struct {
	IsFinalized  bool                      `json:"is_finalized"`
	AnswerCounts []PollResultsAnswerCounts `json:"answer_counts"`
}

// https://docs.discord.com/developers/resources/poll#poll-object
type Poll struct {
	Question         PollQuestion   `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           *time.Time     `json:"expiry"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
	Results          *PollResults   `json:"results,omitempty"`
}

// https://docs.discord.com/developers/resources/poll#poll-create-request-object
type PollRequest struct {
	Question         PollQuestion           `json:"question"`
	Answers          []PollAnswer           `json:"answers"`
	Duration         Option[int]            `json:"duration,omitempty"`
	AllowMultiselect Option[bool]           `json:"allow_multiselect,omitempty"`
	LayoutType       Option[PollLayoutType] `json:"layout_type,omitempty"`
}
