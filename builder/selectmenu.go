package builder

import (
	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// SelectOptionBuilder builds a single StringSelectMenuComponentOption.
//
//	opt := builder.NewSelectOption("Red", "red").SetDescription("The colour red").Build()
//
// https://docs.discord.com/developers/components/reference#string-select-select-option-structure
type SelectOptionBuilder struct {
	opt components.StringSelectMenuComponentOption
}

func NewSelectOption(label, value string) *SelectOptionBuilder {
	return &SelectOptionBuilder{opt: components.StringSelectMenuComponentOption{Label: label, Value: value}}
}

func (b *SelectOptionBuilder) SetDescription(desc string) *SelectOptionBuilder {
	b.opt.Description = desc
	return b
}

func (b *SelectOptionBuilder) SetEmoji(emoji *discord.Emoji) *SelectOptionBuilder {
	b.opt.Emoji = emoji
	return b
}

func (b *SelectOptionBuilder) SetDefault(def bool) *SelectOptionBuilder {
	b.opt.Default = def
	return b
}

func (b *SelectOptionBuilder) Build() components.StringSelectMenuComponentOption {
	return b.opt
}

// ── String select ─────────────────────────────────────────────────────────────

// StringSelectMenuBuilder builds a StringSelectMenu using a fluent API.
//
//	menu := builder.NewStringSelectMenu().
//	    SetCustomID("color_pick").
//	    SetPlaceholder("Pick a colour").
//	    AddOptions(
//	        builder.NewSelectOption("Red", "red").Build(),
//	        builder.NewSelectOption("Blue", "blue").Build(),
//	    ).
//	    Build()
//
// https://docs.discord.com/developers/components/reference#string-select
type StringSelectMenuBuilder struct {
	menu components.StringSelectMenu
}

func NewStringSelectMenu() *StringSelectMenuBuilder {
	var opts []components.StringSelectMenuComponentOption
	return &StringSelectMenuBuilder{menu: components.StringSelectMenu{Options: opts}}
}

func (b *StringSelectMenuBuilder) SetCustomID(id string) *StringSelectMenuBuilder {
	b.menu.CustomID = id
	return b
}

func (b *StringSelectMenuBuilder) SetPlaceholder(placeholder string) *StringSelectMenuBuilder {
	b.menu.Placeholder = placeholder
	return b
}

func (b *StringSelectMenuBuilder) SetMinValues(min int) *StringSelectMenuBuilder {
	b.menu.MinValues = &min
	return b
}

func (b *StringSelectMenuBuilder) SetMaxValues(max int) *StringSelectMenuBuilder {
	b.menu.MaxValues = &max
	return b
}

func (b *StringSelectMenuBuilder) SetDisabled(disabled bool) *StringSelectMenuBuilder {
	b.menu.Disabled = disabled
	return b
}

func (b *StringSelectMenuBuilder) AddOptions(opts ...components.StringSelectMenuComponentOption) *StringSelectMenuBuilder {
	b.menu.Options = append(b.menu.Options, opts...)
	return b
}

func (b *StringSelectMenuBuilder) Build() *components.StringSelectMenu {
	return &b.menu
}

// ── User select ───────────────────────────────────────────────────────────────

// UserSelectMenuBuilder builds a UserSelectMenu using a fluent API.
// https://docs.discord.com/developers/components/reference#user-select
type UserSelectMenuBuilder struct {
	menu components.UserSelectMenu
}

func NewUserSelectMenu() *UserSelectMenuBuilder { return &UserSelectMenuBuilder{} }

func (b *UserSelectMenuBuilder) SetCustomID(id string) *UserSelectMenuBuilder {
	b.menu.CustomID = id
	return b
}

func (b *UserSelectMenuBuilder) SetPlaceholder(placeholder string) *UserSelectMenuBuilder {
	b.menu.Placeholder = placeholder
	return b
}

func (b *UserSelectMenuBuilder) SetMinValues(min int) *UserSelectMenuBuilder {
	b.menu.MinValues = &min
	return b
}

func (b *UserSelectMenuBuilder) SetMaxValues(max int) *UserSelectMenuBuilder {
	b.menu.MaxValues = &max
	return b
}

func (b *UserSelectMenuBuilder) SetDisabled(disabled bool) *UserSelectMenuBuilder {
	b.menu.Disabled = disabled
	return b
}

func (b *UserSelectMenuBuilder) AddDefaultValues(values ...components.SelectDefaultValue) *UserSelectMenuBuilder {
	if b.menu.DefaultValues == nil {
		dv := []components.SelectDefaultValue{}
		b.menu.DefaultValues = dv
	}
	b.menu.DefaultValues = append(b.menu.DefaultValues, values...)
	return b
}

func (b *UserSelectMenuBuilder) Build() *components.UserSelectMenu {
	return &b.menu
}

// ── Role select ───────────────────────────────────────────────────────────────

// RoleSelectMenuBuilder builds a RoleSelectMenu using a fluent API.
// https://docs.discord.com/developers/components/reference#role-select
type RoleSelectMenuBuilder struct {
	menu components.RoleSelectMenu
}

func NewRoleSelectMenu() *RoleSelectMenuBuilder { return &RoleSelectMenuBuilder{} }

func (b *RoleSelectMenuBuilder) SetCustomID(id string) *RoleSelectMenuBuilder {
	b.menu.CustomID = id
	return b
}

func (b *RoleSelectMenuBuilder) SetPlaceholder(placeholder string) *RoleSelectMenuBuilder {
	b.menu.Placeholder = placeholder
	return b
}

func (b *RoleSelectMenuBuilder) SetMinValues(min int) *RoleSelectMenuBuilder {
	b.menu.MinValues = &min
	return b
}

func (b *RoleSelectMenuBuilder) SetMaxValues(max int) *RoleSelectMenuBuilder {
	b.menu.MaxValues = &max
	return b
}

func (b *RoleSelectMenuBuilder) SetDisabled(disabled bool) *RoleSelectMenuBuilder {
	b.menu.Disabled = disabled
	return b
}

func (b *RoleSelectMenuBuilder) AddDefaultValues(values ...components.SelectDefaultValue) *RoleSelectMenuBuilder {
	if b.menu.DefaultValues == nil {
		dv := []components.SelectDefaultValue{}
		b.menu.DefaultValues = dv
	}
	b.menu.DefaultValues = append(b.menu.DefaultValues, values...)
	return b
}

func (b *RoleSelectMenuBuilder) Build() *components.RoleSelectMenu {
	return &b.menu
}

// ── Channel select ────────────────────────────────────────────────────────────

// ChannelSelectMenuBuilder builds a ChannelSelectMenu using a fluent API.
// https://docs.discord.com/developers/components/reference#channel-select
type ChannelSelectMenuBuilder struct {
	menu components.ChannelSelectMenu
}

func NewChannelSelectMenu() *ChannelSelectMenuBuilder { return &ChannelSelectMenuBuilder{} }

func (b *ChannelSelectMenuBuilder) SetCustomID(id string) *ChannelSelectMenuBuilder {
	b.menu.CustomID = id
	return b
}

func (b *ChannelSelectMenuBuilder) SetPlaceholder(placeholder string) *ChannelSelectMenuBuilder {
	b.menu.Placeholder = placeholder
	return b
}

func (b *ChannelSelectMenuBuilder) SetMinValues(min int) *ChannelSelectMenuBuilder {
	b.menu.MinValues = &min
	return b
}

func (b *ChannelSelectMenuBuilder) SetMaxValues(max int) *ChannelSelectMenuBuilder {
	b.menu.MaxValues = &max
	return b
}

func (b *ChannelSelectMenuBuilder) SetDisabled(disabled bool) *ChannelSelectMenuBuilder {
	b.menu.Disabled = disabled
	return b
}

func (b *ChannelSelectMenuBuilder) AddDefaultValues(values ...components.SelectDefaultValue) *ChannelSelectMenuBuilder {
	if b.menu.DefaultValues == nil {
		dv := []components.SelectDefaultValue{}
		b.menu.DefaultValues = dv
	}
	b.menu.DefaultValues = append(b.menu.DefaultValues, values...)
	return b
}

func (b *ChannelSelectMenuBuilder) Build() *components.ChannelSelectMenu {
	return &b.menu
}

// ── Mentionable select ────────────────────────────────────────────────────────

// MentionableSelectMenuBuilder builds a MentionableSelectMenu using a fluent API.
// https://docs.discord.com/developers/components/reference#mentionable-select
type MentionableSelectMenuBuilder struct {
	menu components.MentionableSelectMenu
}

func NewMentionableSelectMenu() *MentionableSelectMenuBuilder {
	return &MentionableSelectMenuBuilder{}
}

func (b *MentionableSelectMenuBuilder) SetCustomID(id string) *MentionableSelectMenuBuilder {
	b.menu.CustomID = id
	return b
}

func (b *MentionableSelectMenuBuilder) SetPlaceholder(placeholder string) *MentionableSelectMenuBuilder {
	b.menu.Placeholder = placeholder
	return b
}

func (b *MentionableSelectMenuBuilder) SetMinValues(min int) *MentionableSelectMenuBuilder {
	b.menu.MinValues = &min
	return b
}

func (b *MentionableSelectMenuBuilder) SetMaxValues(max int) *MentionableSelectMenuBuilder {
	b.menu.MaxValues = &max
	return b
}

func (b *MentionableSelectMenuBuilder) SetDisabled(disabled bool) *MentionableSelectMenuBuilder {
	b.menu.Disabled = disabled
	return b
}

func (b *MentionableSelectMenuBuilder) AddDefaultValues(values ...components.SelectDefaultValue) *MentionableSelectMenuBuilder {
	if b.menu.DefaultValues == nil {
		dv := []components.SelectDefaultValue{}
		b.menu.DefaultValues = dv
	}
	b.menu.DefaultValues = append(b.menu.DefaultValues, values...)
	return b
}

func (b *MentionableSelectMenuBuilder) Build() *components.MentionableSelectMenu {
	return &b.menu
}
