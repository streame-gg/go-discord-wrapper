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
		{"ButtonComponent", func() error { var v ButtonComponent; return json.Unmarshal(null, &v) }},
		{"ChannelSelectMenuComponent", func() error {
			var v ChannelSelectMenuComponent
			return json.Unmarshal(null, &v)
		}},
		{"ChannelComponentInteractionResponse", func() error {
			var v ChannelComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"CheckboxComponent", func() error { var v CheckboxComponent; return json.Unmarshal(null, &v) }},
		{"CheckboxComponentInteractionResponse", func() error {
			var v CheckboxComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"CheckboxGroupComponent", func() error {
			var v CheckboxGroupComponent
			return json.Unmarshal(null, &v)
		}},
		{"CheckboxGroupComponentInteractionResponse", func() error {
			var v CheckboxGroupComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"Container", func() error { var v Container; return json.Unmarshal(null, &v) }},
		{"FileComponent", func() error { var v FileComponent; return json.Unmarshal(null, &v) }},
		{"FileUploadComponent", func() error {
			var v FileUploadComponent
			return json.Unmarshal(null, &v)
		}},
		{"FileUploadComponentInteractionResponse", func() error {
			var v FileUploadComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"LabelComponent", func() error { var v LabelComponent; return json.Unmarshal(null, &v) }},
		{"LabelComponentInteractionResponse", func() error {
			var v LabelComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"MediaGalleryComponent", func() error {
			var v MediaGalleryComponent
			return json.Unmarshal(null, &v)
		}},
		{"MentionableSelectMenuComponent", func() error {
			var v MentionableSelectMenuComponent
			return json.Unmarshal(null, &v)
		}},
		{"MentionableComponentInteractionResponse", func() error {
			var v MentionableComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"RadioGroupComponent", func() error {
			var v RadioGroupComponent
			return json.Unmarshal(null, &v)
		}},
		{"RadioGroupComponentInteractionResponse", func() error {
			var v RadioGroupComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"RoleSelectMenuComponent", func() error {
			var v RoleSelectMenuComponent
			return json.Unmarshal(null, &v)
		}},
		{"RoleComponentInteractionResponse", func() error {
			var v RoleComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"SeparatorComponent", func() error {
			var v SeparatorComponent
			return json.Unmarshal(null, &v)
		}},
		{"StringSelectMenuComponent", func() error {
			var v StringSelectMenuComponent
			return json.Unmarshal(null, &v)
		}},
		{"StringSelectComponentInteractionResponse", func() error {
			var v StringSelectComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"TextDisplayComponent", func() error {
			var v TextDisplayComponent
			return json.Unmarshal(null, &v)
		}},
		{"TextDisplayComponentInteractionResponse", func() error {
			var v TextDisplayComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"TextInputComponent", func() error {
			var v TextInputComponent
			return json.Unmarshal(null, &v)
		}},
		{"TextInputComponentInteractionResponse", func() error {
			var v TextInputComponentInteractionResponse
			return json.Unmarshal(null, &v)
		}},
		{"ThumbnailComponent", func() error {
			var v ThumbnailComponent
			return json.Unmarshal(null, &v)
		}},
		{"UserSelectMenuComponent", func() error {
			var v UserSelectMenuComponent
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
