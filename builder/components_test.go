package builder

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/types/components"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Button ────────────────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestButtonBuilder_AllSetters() {
	t := su.T()
	emoji := &discord.Emoji{Name: "🔥"}
	sku := discord.Snowflake(123)

	btn := NewButton().
		SetStyle(components.ButtonStylePrimary).
		SetLabel("Click me").
		SetCustomID("my_button").
		SetURL("https://example.com").
		SetEmoji(emoji).
		SetSKUID(sku).
		SetDisabled(true).
		Build()

	assert.Equal(t, components.ButtonStylePrimary, btn.Style)
	assert.Equal(t, "Click me", btn.Label)
	assert.Equal(t, "my_button", btn.CustomID)
	assert.Equal(t, "https://example.com", btn.URL)
	assert.Same(t, emoji, btn.Emoji)
	require.NotNil(t, btn.SkuID)
	assert.Equal(t, sku, btn.SkuID)
	assert.True(t, btn.Disabled)
}

// ── ActionRow ─────────────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestActionRowBuilder_AddComponentsAppends() {
	t := su.T()
	b := NewActionRow()
	assert.Empty(t, b.Build().Components)

	b.AddComponents(NewButton().SetCustomID("a").Build())
	b.AddComponents(NewButton().SetCustomID("b").Build(), NewButton().SetCustomID("c").Build())

	assert.Len(t, b.Build().Components, 3)
}

// ── String select ─────────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestSelectOptionBuilder_AllSetters() {
	t := su.T()
	emoji := &discord.Emoji{Name: "x"}
	opt := NewSelectOption("Red", "red").
		SetDescription("The colour red").
		SetEmoji(emoji).
		SetDefault(true).
		Build()

	assert.Equal(t, "Red", opt.Label)
	assert.Equal(t, "red", opt.Value)
	assert.Equal(t, "The colour red", opt.Description)
	assert.Same(t, emoji, opt.Emoji)
	assert.True(t, opt.Default)
}

func (su *builderComponentsSuite) TestStringSelectMenuBuilder_AllSetters() {
	t := su.T()
	menu := NewStringSelectMenu().
		SetCustomID("color_pick").
		SetPlaceholder("Pick a colour").
		SetMinValues(1).
		SetMaxValues(2).
		SetDisabled(true).
		AddOptions(NewSelectOption("Red", "red").Build()).
		AddOptions(NewSelectOption("Blue", "blue").Build()).
		Build()

	assert.Equal(t, "color_pick", menu.CustomID)
	assert.Equal(t, "Pick a colour", menu.Placeholder)
	require.NotNil(t, menu.MinValues)
	assert.Equal(t, 1, *menu.MinValues)
	require.NotNil(t, menu.MaxValues)
	assert.Equal(t, 2, *menu.MaxValues)
	assert.True(t, menu.Disabled)
	require.NotNil(t, menu.Options)
	assert.Len(t, menu.Options, 2)
}

// ── Entity selects (user/role/channel/mentionable) ────────────────────────────
//
// These share the AddDefaultValues nil-init-then-append branch; exercise it on
// each so both the nil and non-nil paths run.

func defaultVal(id uint64, typ components.SelectDefaultValueType) components.SelectDefaultValue {
	return components.SelectDefaultValue{ID: discord.Snowflake(id), Type: typ}
}

func (su *builderComponentsSuite) TestUserSelectMenuBuilder_AllSetters() {
	t := su.T()
	menu := NewUserSelectMenu().
		SetCustomID("u").
		SetPlaceholder("pick user").
		SetMinValues(0).
		SetMaxValues(5).
		SetDisabled(true).
		AddDefaultValues(defaultVal(1, components.SelectDefaultValueTypeUser)).
		AddDefaultValues(defaultVal(2, components.SelectDefaultValueTypeUser)).
		Build()

	assert.Equal(t, "u", menu.CustomID)
	assert.Equal(t, "pick user", menu.Placeholder)
	require.NotNil(t, menu.MinValues)
	assert.Equal(t, 0, *menu.MinValues)
	require.NotNil(t, menu.MaxValues)
	assert.Equal(t, 5, *menu.MaxValues)
	assert.True(t, menu.Disabled)
	require.NotNil(t, menu.DefaultValues)
	assert.Len(t, menu.DefaultValues, 2)
}

func (su *builderComponentsSuite) TestRoleSelectMenuBuilder_AllSetters() {
	t := su.T()
	menu := NewRoleSelectMenu().
		SetCustomID("r").
		SetPlaceholder("pick role").
		SetMinValues(1).
		SetMaxValues(3).
		SetDisabled(false).
		AddDefaultValues(defaultVal(1, components.SelectDefaultValueTypeRole)).
		AddDefaultValues(defaultVal(2, components.SelectDefaultValueTypeRole)).
		Build()

	assert.Equal(t, "r", menu.CustomID)
	require.NotNil(t, menu.DefaultValues)
	assert.Len(t, menu.DefaultValues, 2)
	assert.False(t, menu.Disabled)
}

func (su *builderComponentsSuite) TestChannelSelectMenuBuilder_AllSetters() {
	t := su.T()
	menu := NewChannelSelectMenu().
		SetCustomID("c").
		SetPlaceholder("pick channel").
		SetMinValues(1).
		SetMaxValues(1).
		SetDisabled(true).
		AddDefaultValues(defaultVal(1, components.SelectDefaultValueTypeChannel)).
		AddDefaultValues(defaultVal(2, components.SelectDefaultValueTypeChannel)).
		Build()

	assert.Equal(t, "c", menu.CustomID)
	require.NotNil(t, menu.DefaultValues)
	assert.Len(t, menu.DefaultValues, 2)
}

func (su *builderComponentsSuite) TestMentionableSelectMenuBuilder_AllSetters() {
	t := su.T()
	menu := NewMentionableSelectMenu().
		SetCustomID("m").
		SetPlaceholder("pick mention").
		SetMinValues(1).
		SetMaxValues(2).
		SetDisabled(true).
		AddDefaultValues(defaultVal(1, components.SelectDefaultValueTypeUser)).
		AddDefaultValues(defaultVal(2, components.SelectDefaultValueTypeRole)).
		Build()

	assert.Equal(t, "m", menu.CustomID)
	require.NotNil(t, menu.DefaultValues)
	assert.Len(t, menu.DefaultValues, 2)
}

// ── Text input + Label ────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestTextInputBuilder_AllSetters() {
	t := su.T()
	in := NewTextInput().
		SetCustomID("name").
		SetStyle(components.TextInputStyleShort).
		SetMinLength(2).
		SetMaxLength(10).
		SetRequired(true).
		SetValue("preset").
		SetPlaceholder("Your name").
		Build()

	assert.Equal(t, "name", in.CustomID)
	assert.Equal(t, components.TextInputStyleShort, in.Style)
	require.NotNil(t, in.MinLength)
	assert.Equal(t, 2, *in.MinLength)
	require.NotNil(t, in.MaxLength)
	assert.Equal(t, 10, *in.MaxLength)
	require.NotNil(t, in.Required)
	assert.True(t, in.Required)
	assert.Equal(t, "preset", in.Value)
	assert.Equal(t, "Your name", in.Placeholder)
}

func (su *builderComponentsSuite) TestLabelBuilder_AllSetters() {
	t := su.T()
	child := NewTextInput().SetCustomID("name").Build()
	lbl := NewLabel().
		SetLabel("Your name").
		SetDescription("Enter your full name").
		SetComponent(child).
		Build()

	assert.Equal(t, "Your name", lbl.Label)
	assert.Equal(t, "Enter your full name", lbl.Description)
	assert.Equal(t, child, lbl.Component)
}

// ── Section + Thumbnail ───────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestSectionBuilder_AllSetters() {
	t := su.T()
	thumb := NewThumbnail().SetURL("https://example.com/img.png").Build()
	s := NewSection().
		AddComponents(NewTextDisplay().SetContent("Hello").Build()).
		AddComponents(NewTextDisplay().SetContent("World").Build()).
		SetAccessory(thumb).
		Build()

	assert.Len(t, s.Components, 2)
	assert.Equal(t, thumb, s.Accessory)
}

func (su *builderComponentsSuite) TestThumbnailBuilder_AllSetters() {
	t := su.T()
	thumb := NewThumbnail().
		SetURL("https://example.com/img.png").
		SetDescription("a thumb").
		SetSpoiler(true).
		Build()

	require.NotNil(t, thumb.Media)
	assert.Equal(t, "https://example.com/img.png", thumb.Media.URL)
	assert.Equal(t, "a thumb", thumb.Description)
	assert.True(t, thumb.Spoiler)
}

// ── Container, TextDisplay, Separator ─────────────────────────────────────────

func (su *builderComponentsSuite) TestContainerBuilder_AllSetters() {
	t := su.T()
	c := NewContainer().
		SetAccentColor(0x5865F2).
		SetSpoiler(true).
		AddComponents(NewTextDisplay().SetContent("Hello!").Build()).
		AddComponents(NewTextDisplay().SetContent("Again").Build()).
		Build()

	assert.Equal(t, 0x5865F2, *c.AccentColor)
	assert.True(t, c.Spoiler)
	assert.Len(t, c.Components, 2)
}

func (su *builderComponentsSuite) TestTextDisplayBuilder_SetContent() {
	t := su.T()
	td := NewTextDisplay().SetContent("## Hello world").Build()
	assert.Equal(t, "## Hello world", td.Content)
}

func (su *builderComponentsSuite) TestSeparatorBuilder_AllSetters() {
	t := su.T()
	sep := NewSeparator().
		SetDivider(true).
		SetSpacing(components.SeparatorComponentSpacingLarge).
		Build()

	assert.True(t, sep.Divider)
	assert.Equal(t, components.SeparatorComponentSpacingLarge, sep.Spacing)
}

// ── File display + upload ─────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestFileDisplayBuilder_AllSetters() {
	t := su.T()
	f := NewFileDisplay("https://cdn.example.com/report.pdf").
		SetName("report.pdf").
		SetSpoiler(true).
		Build()

	require.NotNil(t, f.File)
	assert.Equal(t, "https://cdn.example.com/report.pdf", f.File.URL)
	assert.Equal(t, "report.pdf", f.Name)
	assert.True(t, f.Spoiler)
}

func (su *builderComponentsSuite) TestFileUploadBuilder_AllSetters() {
	t := su.T()
	fu := NewFileUpload("attachment").
		SetCustomID("attachment2").
		SetRequired(true).
		SetMinValues(1).
		SetMaxValues(3).
		Build()

	assert.Equal(t, "attachment2", fu.CustomID)
	require.NotNil(t, fu.Required)
	assert.True(t, fu.Required)
	require.NotNil(t, fu.MinValues)
	assert.Equal(t, 1, *fu.MinValues)
	require.NotNil(t, fu.MaxValues)
	assert.Equal(t, 3, *fu.MaxValues)
}

// ── Media gallery ─────────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestMediaGalleryBuilder_AllSetters() {
	t := su.T()
	item := NewMediaGalleryItem("https://example.com/img.png").
		SetDescription("A photo").
		SetSpoiler(true).
		Build()

	require.NotNil(t, item.Media)
	assert.Equal(t, "https://example.com/img.png", item.Media.URL)
	assert.Equal(t, "A photo", *item.Description)
	assert.True(t, item.Spoiler)

	gallery := NewMediaGallery().
		AddItems(item).
		AddItems(NewMediaGalleryItem("https://example.com/b.png").Build()).
		Build()

	require.NotNil(t, gallery.Items)
	assert.Len(t, gallery.Items, 2)
}

// ── Radio group ───────────────────────────────────────────────────────────────

func (su *builderComponentsSuite) TestRadioGroupOptionBuilder_AllSetters() {
	t := su.T()
	opt := NewRadioGroupOption("Yes", "yes").
		SetDescription("Confirm").
		SetDefault(true).
		Build()

	assert.Equal(t, "Yes", opt.Label)
	assert.Equal(t, "yes", opt.Value)
	require.NotNil(t, opt.Description)
	assert.Equal(t, "Confirm", opt.Description)
	require.NotNil(t, opt.Default)
	assert.True(t, opt.Default)
}

func (su *builderComponentsSuite) TestRadioGroupBuilder_AllSetters() {
	t := su.T()
	rg := NewRadioGroup("size_pick").
		SetCustomID("size_pick2").
		AddOptions(NewRadioGroupOption("Small", "sm").Build()).
		AddOptions(NewRadioGroupOption("Large", "lg").SetDefault(true).Build()).
		SetRequired(true).
		Build()

	assert.Equal(t, "size_pick2", rg.CustomID)
	require.NotNil(t, rg.Options)
	assert.Len(t, rg.Options, 2)
	require.NotNil(t, rg.Required)
	assert.True(t, rg.Required)
}

// ── Checkbox group + checkbox ─────────────────────────────────────────────────

func (su *builderComponentsSuite) TestCheckboxGroupOptionBuilder_AllSetters() {
	t := su.T()
	opt := NewCheckboxGroupOption("Notifications", "notif").
		SetDescription("email me").
		SetDefault(true).
		Build()

	assert.Equal(t, "Notifications", opt.Label)
	assert.Equal(t, "notif", opt.Value)
	require.NotNil(t, opt.Description)
	assert.Equal(t, "email me", opt.Description)
	require.NotNil(t, opt.Default)
	assert.True(t, opt.Default)
}

func (su *builderComponentsSuite) TestCheckboxGroupBuilder_AllSetters() {
	t := su.T()
	cg := NewCheckboxGroup("prefs").
		SetCustomID("prefs2").
		AddOptions(NewCheckboxGroupOption("Dark mode", "dark").SetDefault(true).Build()).
		AddOptions(NewCheckboxGroupOption("Notifications", "notif").Build()).
		SetMinValues(1).
		SetMaxValues(2).
		SetRequired(true).
		Build()

	assert.Equal(t, "prefs2", cg.CustomID)
	require.NotNil(t, cg.Options)
	assert.Len(t, cg.Options, 2)
	require.NotNil(t, cg.MinValues)
	assert.Equal(t, 1, *cg.MinValues)
	require.NotNil(t, cg.MaxValues)
	assert.Equal(t, 2, *cg.MaxValues)
	require.NotNil(t, cg.Required)
	assert.True(t, cg.Required)
}

func (su *builderComponentsSuite) TestCheckboxBuilder_AllSetters() {
	t := su.T()
	cb := NewCheckbox("agree").
		SetCustomID("agree2").
		SetDefault(true).
		Build()

	assert.Equal(t, "agree2", cb.CustomID)
	require.NotNil(t, cb.Default)
	assert.True(t, cb.Default)
}

type builderComponentsSuite struct{ suite.Suite }

func TestBuilderComponentsSuite(t *testing.T) { suite.Run(t, new(builderComponentsSuite)) }
