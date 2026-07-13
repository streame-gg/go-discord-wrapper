package responses

import "github.com/streame-gg/go-discord-wrapper/types/discord"

func (s *responsesSuite) TestComponent() {
	component := InteractionDataMessageComponent{}
	s.EqualValues(discord.InteractionTypeMessageComponent, component.GetType())
}
