package discord

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/suite"
)

type messageTestSuite struct {
	suite.Suite
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(messageTestSuite))
}

func (s *messageTestSuite) TestMarshalRoundtrip() {
	lol := Message{
		ID:              123456,
		ChannelID:       Snowflake(1),
		EditedTimestamp: util.PointerOf(time.Now().Round(0)),
	}

	raw, err := json.Marshal(&lol)
	s.Require().NoError(err)
	s.Require().NotEmpty(string(raw))

	lol2 := &Message{}
	s.Require().NoError(json.Unmarshal(raw, lol2))
	s.Equal(lol, *lol2)
}

func (s *messageTestSuite) TestEditedTimestampUnmarshal_EditedTimestampIsNull() {
	raw := "{\"author\":{\"avatar\":null,\"discriminator\":\"\",\"flags\":0,\"id\":\"0\",\"mfa_enabled\":false,\"public_flags\":0,\"username\":\"\"},\"channel_id\":\"1\",\"channel_type\":0,\"content\":\"\",\"edited_timestamp\":null,\"flags\":0,\"id\":\"123456\",\"mention_everyone\":false,\"mention_roles\":null,\"mentions\":null,\"pinned\":false,\"tts\":false,\"type\":0}"

	lol2 := &Message{}
	s.Require().NoError(json.Unmarshal([]byte(raw), lol2))
	s.Nil(lol2.EditedTimestamp)
}

func (s *messageTestSuite) TestIsWebhookMessage() {
	msg := Message{WebhookID: util.PointerOf(Snowflake(1))}
	s.True(msg.IsWebhookMessage())

	msg = Message{}
	s.False(msg.IsWebhookMessage())
}
