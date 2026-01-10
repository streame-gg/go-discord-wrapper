package types

import "encoding/json"

type AnyGuildWrapper struct {
	Guild AnyGuild
}

func (ag *AnyGuildWrapper) UnmarshalJSON(data []byte) error {
	var probe struct {
		Unavailable *bool `json:"unavailable"`
	}

	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	if probe.Unavailable != nil && *probe.Unavailable {
		var ug UnavailableGuild
		if err := json.Unmarshal(data, &ug); err != nil {
			return err
		}
		ag.Guild = ug
		return nil
	}

	var g Guild
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}
	ag.Guild = g
	return nil
}

type AnyGuild interface {
	IsAvailable() bool
	GetID() DiscordSnowflake
}

type Guild struct {
	ID   DiscordSnowflake `json:"id"`
	Name string           `json:"name"`
}

func (g Guild) IsAvailable() bool {
	return true
}

func (g Guild) GetID() DiscordSnowflake {
	return g.ID
}

type UnavailableGuild struct {
	ID          DiscordSnowflake `json:"id"`
	Unavailable *bool            `json:"unavailable"`
	Name        string           `json:"name,omitempty"`
}

func (ug UnavailableGuild) IsAvailable() bool {
	return !*ug.Unavailable
}

func (ug UnavailableGuild) GetID() DiscordSnowflake {
	return ug.ID
}
