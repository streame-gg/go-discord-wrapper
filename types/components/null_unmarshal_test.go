package components

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestNullUnmarshalNoPanic verifies that every component type with a custom
// UnmarshalJSON safely handles a JSON null input without panicking (P1-13).
func (su *nullUnmarshalSuite) TestNullUnmarshalNoPanic() {
	t := su.T()
	null := []byte("null")

	cases := []struct {
		name string
		fn   func() error
	}{
		{"Button", func() error { var v Button; return json.Unmarshal(null, &v) }},
		{"ChannelSelectMenu", func() error {
			var v ChannelSelectMenu
			return json.Unmarshal(null, &v)
		}},
		{"ChannelComponentInteractionResponse", func() error {
			var v ChannelComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Checkbox", func() error { var v Checkbox; return json.Unmarshal(null, &v) }},
		{"CheckboxComponentInteractionResponse", func() error {
			var v CheckboxComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"CheckboxGroup", func() error {
			var v CheckboxGroup
			return json.Unmarshal(null, &v)
		}},
		{"CheckboxGroupComponentInteractionResponse", func() error {
			var v CheckboxGroupComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Container", func() error { var v Container; return json.Unmarshal(null, &v) }},
		{"File", func() error { var v File; return json.Unmarshal(null, &v) }},
		{"FileUpload", func() error {
			var v FileUpload
			return json.Unmarshal(null, &v)
		}},
		{"FileUploadComponentInteractionResponse", func() error {
			var v FileUploadComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Label", func() error { var v Label; return json.Unmarshal(null, &v) }},
		{"LabelComponentInteractionResponse", func() error {
			var v LabelComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"MediaGallery", func() error {
			var v MediaGallery
			return json.Unmarshal(null, &v)
		}},
		{"MentionableSelectMenu", func() error {
			var v MentionableSelectMenu
			return json.Unmarshal(null, &v)
		}},
		{"MentionableComponentInteractionResponse", func() error {
			var v MentionableComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"RadioGroup", func() error {
			var v RadioGroup
			return json.Unmarshal(null, &v)
		}},
		{"RadioGroupComponentInteractionResponse", func() error {
			var v RadioGroupComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"RoleSelectMenu", func() error {
			var v RoleSelectMenu
			return json.Unmarshal(null, &v)
		}},
		{"RoleComponentInteractionResponse", func() error {
			var v RoleComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Separator", func() error {
			var v Separator
			return json.Unmarshal(null, &v)
		}},
		{"StringSelectMenu", func() error {
			var v StringSelectMenu
			return json.Unmarshal(null, &v)
		}},
		{"StringSelectComponentInteractionResponse", func() error {
			var v StringSelectComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"TextDisplay", func() error {
			var v TextDisplay
			return json.Unmarshal(null, &v)
		}},
		{"TextDisplayComponentInteractionResponse", func() error {
			var v TextDisplayComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"TextInput", func() error {
			var v TextInput
			return json.Unmarshal(null, &v)
		}},
		{"TextInputComponentInteractionResponse", func() error {
			var v TextInputComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Thumbnail", func() error {
			var v Thumbnail
			return json.Unmarshal(null, &v)
		}},
		{"UserSelectMenu", func() error {
			var v UserSelectMenu
			return json.Unmarshal(null, &v)
		}},
		{"UserSelectComponentInteractionResponse", func() error {
			var v UserSelectComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.fn(); err != nil {
				t.Errorf("unexpected error for %s: %v", tc.name, err)
			}
		})
	}
}

type nullUnmarshalSuite struct{ suite.Suite }

func TestNullUnmarshalSuite(t *testing.T) { suite.Run(t, new(nullUnmarshalSuite)) }
