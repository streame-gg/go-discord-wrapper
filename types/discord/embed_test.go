package discord

import (
	"encoding/json"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/stretchr/testify/suite"
)

type embedSuite struct {
	suite.Suite
}

func TestEmbedSuite(t *testing.T) {
	suite.Run(t, new(embedSuite))
}

func (s *embedSuite) TestEmbed() {
	embed := Embed{
		Title: "title",
		Color: util.PointerOf(0),
	}

	raw, err := json.Marshal(embed)
	s.Require().NoError(err)
	s.Require().Equal(`{"title":"title","color":0}`, string(raw))
}
