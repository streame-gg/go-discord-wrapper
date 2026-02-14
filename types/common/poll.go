package common

import "time"

type PollLayoutType int

const (
	PollLayoutTypeDefault PollLayoutType = 0
)

type PollQuestion struct {
	Text  *string `json:"text"`
	Emoji *Emoji  `json:"emoji,omitempty"`
}

type PollAnswer struct {
	AnswerID  int          `json:"answer_id"`
	PollMedia PollQuestion `json:"poll_media"`
}

type PollResultsAnswerCounts struct {
	ID      int  `json:"id"`
	Count   int  `json:"count"`
	MeVoted bool `json:"me_voted"`
}

type PollResults struct {
	IsFinalized  bool                      `json:"is_finalized"`
	AnswerCounts []PollResultsAnswerCounts `json:"answer_counts"`
}

type Poll struct {
	Question         *PollQuestion  `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           *time.Time     `json:"expiry,omitempty"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
	Results          *PollResults   `json:"results,omitempty"`
}

type PollRequest struct {
	Question         *PollQuestion  `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Duration         *time.Time     `json:"duration,omitempty"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
}
