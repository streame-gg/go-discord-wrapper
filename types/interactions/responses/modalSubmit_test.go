package responses

import "github.com/streame-gg/go-discord-wrapper/types/discord"

func (s *responsesSuite) TestModalSubmit() {
	component := InteractionDataModalSubmit{}
	s.EqualValues(discord.InteractionTypeModalSubmit, component.GetType())
}
