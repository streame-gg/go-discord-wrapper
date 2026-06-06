package connection

// Tests for the new AutoMod rule and Invite gateway event cache handlers.

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/events"
)

// ── AutoMod gateway events ────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestAutoModRuleCreate_CachesRule() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	ruleID := discord.Snowflake(500000000000000005)

	payload := map[string]any{
		"id":               ruleID.String(),
		"guild_id":         guildID.String(),
		"name":             "no-spam",
		"creator_id":       "111",
		"event_type":       1,
		"trigger_type":     1,
		"trigger_metadata": map[string]any{},
		"actions":          []any{},
		"enabled":          true,
		"exempt_roles":     []any{},
		"exempt_channels":  []any{},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventAutoModerationRuleCreate, nil)

	rule, ok := c.Cache.AutoModRules().Get(ruleID)
	require.True(t, ok, "AUTO_MODERATION_RULE_CREATE must add rule to AutoModRuleStore")
	assert.Equal(t, "no-spam", rule.Name)
}

func (cs *ConnectionSuite) TestAutoModRuleUpdate_UpdatesCache() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	ruleID := discord.Snowflake(500000000000000005)

	// Seed with original rule.
	c.Cache.AutoModRules().Set(guildID, &discord.AutoModerationRule{
		ID:      ruleID,
		GuildID: guildID,
		Name:    "old-name",
	})

	payload := map[string]any{
		"id":               ruleID.String(),
		"guild_id":         guildID.String(),
		"name":             "new-name",
		"creator_id":       "111",
		"event_type":       1,
		"trigger_type":     1,
		"trigger_metadata": map[string]any{},
		"actions":          []any{},
		"enabled":          true,
		"exempt_roles":     []any{},
		"exempt_channels":  []any{},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	// AUTO_MODERATION_RULE_UPDATE uses a typed cast on the pre-parsed event, so
	// we must unmarshal into the right type before passing it to internalEventHandler.
	var ev events.AutoModerationRuleUpdateEvent
	require.NoError(t, json.Unmarshal(raw, &ev))

	c.internalEventHandler(json.RawMessage(raw), events.EventAutoModerationRuleUpdate, &ev)

	rule, ok := c.Cache.AutoModRules().Get(ruleID)
	require.True(t, ok, "AUTO_MODERATION_RULE_UPDATE must keep rule in AutoModRuleStore")
	assert.Equal(t, "new-name", rule.Name)
}

func (cs *ConnectionSuite) TestAutoModRuleDelete_EvictsRule() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	ruleID := discord.Snowflake(500000000000000005)

	c.Cache.AutoModRules().Set(guildID, &discord.AutoModerationRule{
		ID:      ruleID,
		GuildID: guildID,
		Name:    "to-delete",
	})
	_, ok := c.Cache.AutoModRules().Get(ruleID)
	require.True(t, ok)

	payload := map[string]any{
		"id":               ruleID.String(),
		"guild_id":         guildID.String(),
		"name":             "to-delete",
		"creator_id":       "111",
		"event_type":       1,
		"trigger_type":     1,
		"trigger_metadata": map[string]any{},
		"actions":          []any{},
		"enabled":          true,
		"exempt_roles":     []any{},
		"exempt_channels":  []any{},
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventAutoModerationRuleDelete, nil)

	_, ok = c.Cache.AutoModRules().Get(ruleID)
	assert.False(t, ok, "AUTO_MODERATION_RULE_DELETE must evict rule from AutoModRuleStore")
}

func (cs *ConnectionSuite) TestGuildDelete_ClearsAutoModRules() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	ruleID := discord.Snowflake(500000000000000005)

	c.Cache.AutoModRules().Set(guildID, &discord.AutoModerationRule{
		ID:      ruleID,
		GuildID: guildID,
	})
	require.Equal(t, 1, c.Cache.AutoModRules().GetByGuild(guildID).Len())

	dispatchGuildDelete(t, c, guildID)

	assert.Equal(t, 0, c.Cache.AutoModRules().GetByGuild(guildID).Len(),
		"GUILD_DELETE must evict all automod rules for the guild")
}

// ── Invite gateway events ─────────────────────────────────────────────────────

func (cs *ConnectionSuite) TestInviteCreate_CachesInvite() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	channelID := discord.Snowflake(200000000000000002)

	payload := map[string]any{
		"code":       "abc123",
		"channel_id": channelID.String(),
		"guild_id":   guildID.String(),
		"created_at": "2024-01-01T00:00:00Z",
		"temporary":  false,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventInviteCreate, nil)

	inv, ok := c.Cache.Invites().Get("abc123")
	require.True(t, ok, "INVITE_CREATE must add invite to InviteStore")
	assert.Equal(t, "abc123", inv.Code)
}

func (cs *ConnectionSuite) TestInviteCreate_GuildScoped_GetByGuild() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)
	channelID := discord.Snowflake(200000000000000002)

	payload := map[string]any{
		"code":       "abc123",
		"channel_id": channelID.String(),
		"guild_id":   guildID.String(),
		"created_at": "2024-01-01T00:00:00Z",
		"temporary":  false,
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventInviteCreate, nil)

	all := c.Cache.Invites().GetByGuild(guildID)
	assert.Equal(t, 1, all.Len(), "invite must be retrievable via GetByGuild after INVITE_CREATE")
}

func (cs *ConnectionSuite) TestInviteDelete_EvictsInvite() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)

	// Pre-seed an invite.
	c.Cache.Invites().SetWithGuild(guildID, &discord.Invite{Code: "abc123"})
	_, ok := c.Cache.Invites().Get("abc123")
	require.True(t, ok)

	channelID := discord.Snowflake(200000000000000002)
	payload := map[string]any{
		"code":       "abc123",
		"channel_id": channelID.String(),
		"guild_id":   guildID.String(),
	}
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	c.internalEventHandler(json.RawMessage(raw), events.EventInviteDelete, nil)

	_, ok = c.Cache.Invites().Get("abc123")
	assert.False(t, ok, "INVITE_DELETE must evict invite from InviteStore")
}

func (cs *ConnectionSuite) TestGuildDelete_ClearsInvites() {
	t := cs.T()
	c := newClientWithCache(t)

	guildID := discord.Snowflake(100000000000000001)

	c.Cache.Invites().SetWithGuild(guildID, &discord.Invite{Code: "code1"})
	c.Cache.Invites().SetWithGuild(guildID, &discord.Invite{Code: "code2"})
	require.Equal(t, 2, c.Cache.Invites().GetByGuild(guildID).Len())

	dispatchGuildDelete(t, c, guildID)

	assert.Equal(t, 0, c.Cache.Invites().GetByGuild(guildID).Len(),
		"GUILD_DELETE must evict all invites for the guild")
}
