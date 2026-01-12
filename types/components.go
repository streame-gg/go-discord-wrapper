package types

type DiscordComponentType int

type AnyComponent interface {
	Type() DiscordComponentType
	UnmarshalJSON(data []byte) error
}
