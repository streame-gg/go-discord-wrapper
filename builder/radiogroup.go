package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// ── Radio group ───────────────────────────────────────────────────────────────

// RadioGroupOptionBuilder builds a single RadioGroupComponentOption.
//
//	opt := builder.NewRadioGroupOption("Yes", "yes").SetDescription("Confirm").Build()
//
// https://docs.discord.com/developers/components/reference#radio-group
type RadioGroupOptionBuilder struct {
	opt components.RadioGroupComponentOption
}

func NewRadioGroupOption(label, value string) *RadioGroupOptionBuilder {
	return &RadioGroupOptionBuilder{opt: components.RadioGroupComponentOption{Label: label, Value: value}}
}

func (b *RadioGroupOptionBuilder) SetDescription(desc string) *RadioGroupOptionBuilder {
	b.opt.Description = desc
	return b
}

func (b *RadioGroupOptionBuilder) SetDefault(def bool) *RadioGroupOptionBuilder {
	b.opt.Default = def
	return b
}

func (b *RadioGroupOptionBuilder) Build() components.RadioGroupComponentOption {
	return b.opt
}

// RadioGroupBuilder builds a RadioGroupComponent using a fluent API.
//
//	rg := builder.NewRadioGroup("size_pick").
//	    AddOptions(
//	        builder.NewRadioGroupOption("Small", "sm").Build(),
//	        builder.NewRadioGroupOption("Large", "lg").SetDefault(true).Build(),
//	    ).
//	    SetRequired(true).
//	    Build()
//
// https://docs.discord.com/developers/components/reference#radio-group
type RadioGroupBuilder struct {
	rg components.RadioGroupComponent
}

func NewRadioGroup(customID string) *RadioGroupBuilder {
	var opts []components.RadioGroupComponentOption
	return &RadioGroupBuilder{rg: components.RadioGroupComponent{CustomID: customID, Options: opts}}
}

func (b *RadioGroupBuilder) SetCustomID(id string) *RadioGroupBuilder {
	b.rg.CustomID = id
	return b
}

func (b *RadioGroupBuilder) AddOptions(opts ...components.RadioGroupComponentOption) *RadioGroupBuilder {
	b.rg.Options = append(b.rg.Options, opts...)
	return b
}

func (b *RadioGroupBuilder) SetRequired(required bool) *RadioGroupBuilder {
	b.rg.Required = required
	return b
}

func (b *RadioGroupBuilder) Build() *components.RadioGroupComponent {
	return &b.rg
}

// ── Checkbox group ────────────────────────────────────────────────────────────

// CheckboxGroupOptionBuilder builds a single CheckboxGroupComponentOption.
//
//	opt := builder.NewCheckboxGroupOption("Notifications", "notif").SetDefault(true).Build()
//
// https://docs.discord.com/developers/components/reference#checkbox-group
type CheckboxGroupOptionBuilder struct {
	opt components.CheckboxGroupComponentOption
}

func NewCheckboxGroupOption(label, value string) *CheckboxGroupOptionBuilder {
	return &CheckboxGroupOptionBuilder{opt: components.CheckboxGroupComponentOption{Label: label, Value: value}}
}

func (b *CheckboxGroupOptionBuilder) SetDescription(desc string) *CheckboxGroupOptionBuilder {
	b.opt.Description = desc
	return b
}

func (b *CheckboxGroupOptionBuilder) SetDefault(def bool) *CheckboxGroupOptionBuilder {
	b.opt.Default = def
	return b
}

func (b *CheckboxGroupOptionBuilder) Build() components.CheckboxGroupComponentOption {
	return b.opt
}

// CheckboxGroupBuilder builds a CheckboxGroupComponent using a fluent API.
//
//	cg := builder.NewCheckboxGroup("prefs").
//	    AddOptions(
//	        builder.NewCheckboxGroupOption("Dark mode", "dark").SetDefault(true).Build(),
//	        builder.NewCheckboxGroupOption("Notifications", "notif").Build(),
//	    ).
//	    SetMinValues(1).
//	    SetMaxValues(2).
//	    Build()
//
// https://docs.discord.com/developers/components/reference#checkbox-group
type CheckboxGroupBuilder struct {
	cg components.CheckboxGroupComponent
}

func NewCheckboxGroup(customID string) *CheckboxGroupBuilder {
	var opts []components.CheckboxGroupComponentOption
	return &CheckboxGroupBuilder{cg: components.CheckboxGroupComponent{CustomID: customID, Options: opts}}
}

func (b *CheckboxGroupBuilder) SetCustomID(id string) *CheckboxGroupBuilder {
	b.cg.CustomID = id
	return b
}

func (b *CheckboxGroupBuilder) AddOptions(opts ...components.CheckboxGroupComponentOption) *CheckboxGroupBuilder {
	b.cg.Options = append(b.cg.Options, opts...)
	return b
}

func (b *CheckboxGroupBuilder) SetMinValues(min int) *CheckboxGroupBuilder {
	b.cg.MinValues = &min
	return b
}

func (b *CheckboxGroupBuilder) SetMaxValues(max int) *CheckboxGroupBuilder {
	b.cg.MaxValues = &max
	return b
}

func (b *CheckboxGroupBuilder) SetRequired(required bool) *CheckboxGroupBuilder {
	b.cg.Required = required
	return b
}

func (b *CheckboxGroupBuilder) Build() *components.CheckboxGroupComponent {
	return &b.cg
}

// ── Checkbox ──────────────────────────────────────────────────────────────────

// CheckboxBuilder builds a CheckboxComponent using a fluent API.
//
//	cb := builder.NewCheckbox("agree").SetDefault(false).Build()
//
// https://docs.discord.com/developers/components/reference#checkbox
type CheckboxBuilder struct {
	cb components.CheckboxComponent
}

func NewCheckbox(customID string) *CheckboxBuilder {
	return &CheckboxBuilder{cb: components.CheckboxComponent{CustomID: customID}}
}

func (b *CheckboxBuilder) SetCustomID(id string) *CheckboxBuilder {
	b.cb.CustomID = id
	return b
}

func (b *CheckboxBuilder) SetDefault(def bool) *CheckboxBuilder {
	b.cb.Default = def
	return b
}

func (b *CheckboxBuilder) Build() *components.CheckboxComponent {
	return &b.cb
}
