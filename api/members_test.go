package api

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func TestModifyGuildMemberParams_RolesJSON(t *testing.T) {
	roleID := discord.Snowflake(754246997266923571)

	cases := []struct {
		name        string
		params      ModifyGuildMemberParams
		wantKey     bool   // whether "roles" key should appear in JSON at all
		wantPayload string // expected value substring when key is present
	}{
		{
			name:    "nil roles omitted - do not modify",
			params:  ModifyGuildMemberParams{Roles: nil},
			wantKey: false,
		},
		{
			name:        "empty slice sends roles:[] - removes all roles",
			params:      ModifyGuildMemberParams{Roles: &[]discord.Snowflake{}},
			wantKey:     true,
			wantPayload: `"roles":[]`,
		},
		{
			name:        "non-empty slice sends role IDs",
			params:      ModifyGuildMemberParams{Roles: &[]discord.Snowflake{roleID}},
			wantKey:     true,
			wantPayload: `"roles":["754246997266923571"]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.params)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			got := string(b)
			hasKey := strings.Contains(got, `"roles"`)
			if hasKey != tc.wantKey {
				t.Errorf("json %s: wantKey=%v gotKey=%v", got, tc.wantKey, hasKey)
			}
			if tc.wantPayload != "" && !strings.Contains(got, tc.wantPayload) {
				t.Errorf("json %s: want substring %q", got, tc.wantPayload)
			}
		})
	}
}
