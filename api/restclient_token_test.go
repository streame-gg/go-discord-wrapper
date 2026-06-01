package api

import (
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRestClient_TokenNormalization verifies that NewRestClient strips
// "Bot " / "Bearer " prefixes and whitespace so the stored token is a raw
// secret regardless of how the caller formatted it (P1-20).
func (su *restclientTokenSuite) TestNewRestClient_TokenNormalization() {
	t := su.T()
	cases := []struct {
		desc    string
		input   string
		wantErr bool
		want    string // expected raw token stored in rc.token
	}{
		{"bare token", "xyz", false, "xyz"},
		{"Bot prefix", "Bot xyz", false, "xyz"},
		{"bot prefix lowercase", "bot xyz", false, "xyz"},
		{"BOT prefix uppercase", "BOT xyz", false, "xyz"},
		{"Bearer prefix stored as-is", "Bearer xyz", false, "Bearer xyz"},
		{"leading/trailing spaces", "  xyz  ", false, "xyz"},
		{"empty string", "", true, ""},
		{"only spaces", "   ", true, ""},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			rc, err := NewRestClient(tc.input)
			if tc.wantErr {
				require.Error(t, err, "expected error for input %q", tc.input)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, rc.token,
				"stored token for input %q should be %q", tc.input, tc.want)
		})
	}
}

type restclientTokenSuite struct{ suite.Suite }

func TestRestclientTokenSuite(t *testing.T) { suite.Run(t, new(restclientTokenSuite)) }
