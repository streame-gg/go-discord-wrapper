package events

import (
	"github.com/streame-gg/go-discord-wrapper/internal/testutil"
	"github.com/streame-gg/go-discord-wrapper/testdata"
)

func (s *eventSuite) TestGuildCreate() {
	s.T().Log("Testing Guild Create Unmarshal Logic")

	sub := testutil.InitSub[GuildCreateEvent](s)

	sub.RunCommonEdgeCases()

	guild := testdata.NewAvailableGuild()
	unavailableGuild := testdata.NewUnavailableGuild()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildCreateEvent]{
		{
			Name:  "valid full payload",
			Input: sub.MustMarshal(guild),
			Validate: func(got GuildCreateEvent) {
				s.Require().NotNil(got.Guild)
				s.EqualValues(guild["id"], got.Guild.ID)
				s.EqualValues(guild["name"], got.Guild.Name)
				s.EqualValues(guild["icon"], *got.Guild.Icon)
				s.EqualValues(guild["icon_hash"], *got.Guild.IconHash)
				s.EqualValues(guild["splash"], *got.Guild.Splash)
				s.EqualValues(guild["discovery_splash"], *got.Guild.DiscoverySplash)
				s.EqualValues(guild["owner"], got.Guild.Owner)
				s.EqualValues(guild["owner_id"], got.Guild.OwnerID)
				s.EqualValues(guild["permissions"], got.Guild.Permissions)
				s.EqualValues(guild["region"], *got.Guild.Region)
				s.EqualValues(guild["afk_channel_id"], *got.Guild.AfkChannelID)
				s.EqualValues(guild["afk_timeout"], got.Guild.AfkTimeout)
				s.EqualValues(guild["widget_enabled"], got.Guild.WidgetEnabled)
				s.EqualValues(guild["widget_channel_id"], *got.Guild.WidgetChannelID)
				s.EqualValues(guild["verification_level"], got.Guild.VerificationLevel)
				s.EqualValues(guild["default_message_notifications"], got.Guild.DefaultMessageNotifications)
				s.EqualValues(guild["explicit_content_filter"], got.Guild.ExplicitContentFilter)
				s.EqualValues(guild["mfa_level"], got.Guild.MfaLevel)
				s.EqualValues(guild["application_id"], *got.Guild.ApplicationID)
				s.EqualValues(guild["system_channel_id"], *got.Guild.SystemChannelID)
				s.EqualValues(guild["system_channel_flags"], got.Guild.SystemChannelFlags)
				s.EqualValues(guild["rules_channel_id"], *got.Guild.RulesChannelID)
				s.EqualValues(guild["max_presences"], *got.Guild.MaxPresences)
				s.EqualValues(guild["max_members"], got.Guild.MaxMembers)
				s.EqualValues(guild["vanity_url_code"], *got.Guild.VanityUrlCode)
				s.EqualValues(guild["description"], *got.Guild.Description)
				s.EqualValues(guild["banner"], *got.Guild.Banner)
				s.EqualValues(guild["premium_tier"], got.Guild.PremiumTier)
				s.EqualValues(guild["premium_subscription_count"], got.Guild.PremiumSubscriptionCount)
				s.EqualValues(guild["preferred_locale"], got.Guild.PreferredLocale)
				s.EqualValues(guild["public_updates_channel_id"], *got.Guild.PublicUpdatesChannelID)
				s.EqualValues(guild["max_video_channel_users"], got.Guild.MaxVideoChannelUsers)
				s.EqualValues(guild["max_stage_video_channel_users"], got.Guild.MaxStageVideoChannelUsers)
				s.EqualValues(guild["approximate_member_count"], got.Guild.ApproximateMemberCount)
				s.EqualValues(guild["approximate_presence_count"], got.Guild.ApproximatePresenceCount)
				s.EqualValues(guild["nsfw_level"], got.Guild.NSFWLevel)
				s.EqualValues(guild["premium_progress_bar_enabled"], got.Guild.PremiumProgressBarEnabled)
				s.EqualValues(guild["safety_alerts_channel_id"], *got.Guild.SafetyAlertsChannelID)

				welcomeScreen := guild["welcome_screen"].(map[string]interface{})
				s.EqualValues(welcomeScreen["description"], *got.Guild.WelcomeScreen.Description)
				welcomeChannels := welcomeScreen["welcome_channels"].([]map[string]interface{})
				for i, wc := range welcomeChannels {
					s.EqualValues(wc["channel_id"], got.Guild.WelcomeScreen.WelcomeChannels[i].ChannelID)
					s.EqualValues(wc["description"], got.Guild.WelcomeScreen.WelcomeChannels[i].Description)
					s.EqualValues(wc["emoji_id"], *got.Guild.WelcomeScreen.WelcomeChannels[i].EmojiID)
					s.EqualValues(wc["emoji_name"], *got.Guild.WelcomeScreen.WelcomeChannels[i].EmojiName)
				}

				incidentsData := guild["incidents_data"].(map[string]interface{})
				s.EqualValues(incidentsData["invites_disabled_until"], *got.Guild.IncidentsData.InvitesDisabledUntil)
				s.EqualValues(incidentsData["dms_disabled_until"], *got.Guild.IncidentsData.DmsDisabledUntil)
				s.EqualValues(incidentsData["dm_spam_detected_at"], *got.Guild.IncidentsData.DmSpanDetectedAt)
				s.EqualValues(incidentsData["raid_detected_at"], *got.Guild.IncidentsData.RaidDetectedAt)

				roles := guild["roles"].([]map[string]interface{})
				s.Len(got.Guild.RawRoles, len(roles))

				for i, r := range roles {
					s.compareRole(r, got.Guild.RawRoles[i])
				}

				stickers := guild["stickers"].([]map[string]interface{})
				s.Len(got.Guild.RawStickers, len(stickers))

				for i, st := range stickers {
					s.compareSticker(st, got.Guild.RawStickers[i])
				}

				emojis := guild["emojis"].([]map[string]interface{})
				s.Len(got.Guild.RawEmojis, len(emojis))

				for i, em := range emojis {
					s.compareEmoji(em, got.Guild.RawEmojis[i])
				}

				s.False(got.Unavailable)
			},
		},
		{
			Name:  "unavailable guild",
			Input: sub.MustMarshal(unavailableGuild),
			Validate: func(got GuildCreateEvent) {
				s.True(got.Unavailable)
				s.EqualValues(unavailableGuild["id"], got.ID)
				s.Nil(got.Guild)
			},
		},
	})
}

func (s *eventSuite) TestGuildDelete() {
	s.T().Log("Testing Guild Delete Unmarshal Logic")

	sub := testutil.InitSub[GuildDeleteEvent](s)

	sub.RunCommonEdgeCases()

	unavailableGuild := testdata.NewUnavailableGuild()

	sub.RunCases([]testutil.UnmarshalTestCase[GuildDeleteEvent]{
		{
			Name:  "unavailable guild",
			Input: sub.MustMarshal(unavailableGuild),
			Validate: func(got GuildDeleteEvent) {
				s.True(got.Unavailable)
				s.EqualValues(unavailableGuild["id"], got.ID)
				s.Nil(got.Guild)
			},
		},
	})
}
