package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/streame-gg/go-discord-wrapper/types/common"
)

// ── Param types ───────────────────────────────────────────────────────────────

type GetPollAnswerVotersParams struct {
	After *common.Snowflake
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

func (c *RestClient) GetPollAnswerVoters(ctx context.Context, channelID, messageID common.Snowflake, answerID int, params GetPollAnswerVotersParams) ([]*common.User, error) {
	if err := channelID.Validate(); err != nil {
		return nil, err
	}

	if err := messageID.Validate(); err != nil {
		return nil, err
	}

	path := "/channels/" + channelID.String() + "/polls/" + messageID.String() + "/answers/" + strconv.Itoa(answerID) + params.toQuery()
	req, err := c.generateRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Users []*common.User `json:"users"`
	}
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return result.Users, nil
}

func (c *RestClient) EndPoll(ctx context.Context, channelID, messageID common.Snowflake) (*common.Message, error) {
	path := "/channels/" + channelID.String() + "/polls/" + messageID.String() + "/expire"
	req, err := c.generateRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	var result common.Message
	if _, err := c.do(req, http.StatusOK, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
