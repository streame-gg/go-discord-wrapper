package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// ── Param types ───────────────────────────────────────────────────────────────

type GetPollAnswerVotersParams struct {
	After *discord.Snowflake
	Limit *int
}

func (p GetPollAnswerVotersParams) toQuery() string {
	q := url.Values{}
	if p.After != nil {
		q.Set("after", p.After.String())
	}
	if p.Limit != nil {
		q.Set("limit", strconv.Itoa(*p.Limit))
	}
	if len(q) == 0 {
		return ""
	}
	return "?" + q.Encode()
}

// ── Poll endpoints ────────────────────────────────────────────────────────────

// CreatePoll creates a poll in a channel. Basically a wrapper for CreateMessage.
func (c *RestClient) CreatePoll(ctx context.Context, channelID discord.Snowflake, poll discord.PollRequest) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	return c.CreateMessage(ctx, channelID, CreateMessageParams{
		Poll: poll,
	})
}

type GetPollAnswerVotersResponse struct {
	Users []*discord.User `json:"users"`
}

func (c *RestClient) GetPollAnswerVoters(ctx context.Context, channelID, messageID discord.Snowflake, answerID int, params GetPollAnswerVotersParams) (*GetPollAnswerVotersResponse, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/polls/" + messageID.String() + "/answers/" + strconv.Itoa(answerID) + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[GetPollAnswerVotersResponse](c, req, map[int]bool{
		http.StatusOK: true,
	})
}

func (c *RestClient) EndPoll(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}
	if err := messageID.Validate(); err != nil {
		return nil, err
	}
	path := "/channels/" + channelID.String() + "/polls/" + messageID.String() + "/expire"
	req, err := c.generateRequest(ctx, http.MethodPost, path, nil, c.WithBotAuthorization())
	if err != nil {
		return nil, err
	}

	return doRequest[discord.Message](c, req, map[int]bool{
		http.StatusOK: true,
	})
}
