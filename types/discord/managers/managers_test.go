package managers_test

import (
	"context"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/collection"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
	"github.com/streame-gg/go-discord-wrapper/types/discord/managers"
)

// ── Stub cache stores ─────────────────────────────────────────────────────────

type stubMemberStore struct {
	members map[discord.Snowflake]map[discord.Snowflake]*discord.GuildMember
}

func (s *stubMemberStore) Get(guildID, userID discord.Snowflake) (*discord.GuildMember, bool) {
	gm, ok := s.members[guildID]
	if !ok {
		return nil, false
	}
	m, ok := gm[userID]
	return m, ok
}
func (s *stubMemberStore) AllInGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.GuildMember] {
	c := collection.New[discord.Snowflake, *discord.GuildMember]()
	for uid, m := range s.members[guildID] {
		c.Set(uid, m)
	}
	return c
}
func (s *stubMemberStore) Size() int { return len(s.members) }

type stubRoleStore struct {
	roles map[discord.Snowflake]map[discord.Snowflake]*discord.Role
}

func (s *stubRoleStore) Get(roleID discord.Snowflake) (*discord.Role, bool) {
	for _, gm := range s.roles {
		if r, ok := gm[roleID]; ok {
			return r, true
		}
	}
	return nil, false
}
func (s *stubRoleStore) GetByGuild(guildID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Role] {
	c := collection.New[discord.Snowflake, *discord.Role]()
	for rid, r := range s.roles[guildID] {
		c.Set(rid, r)
	}
	return c
}
func (s *stubRoleStore) All() *collection.Collection[discord.Snowflake, *discord.Role] {
	c := collection.New[discord.Snowflake, *discord.Role]()
	for _, gm := range s.roles {
		for rid, r := range gm {
			c.Set(rid, r)
		}
	}
	return c
}
func (s *stubRoleStore) Size() int { return len(s.roles) }

type stubChannelStore struct {
	channels map[discord.Snowflake]*discord.Channel
}

func (s *stubChannelStore) Get(id discord.Snowflake) (*discord.Channel, bool) {
	ch, ok := s.channels[id]
	return ch, ok
}
func (s *stubChannelStore) All() *collection.Collection[discord.Snowflake, *discord.Channel] {
	c := collection.New[discord.Snowflake, *discord.Channel]()
	for id, ch := range s.channels {
		c.Set(id, ch)
	}
	return c
}
func (s *stubChannelStore) Size() int { return len(s.channels) }

type stubMessageStore struct {
	messages map[discord.Snowflake]map[discord.Snowflake]*discord.Message
}

func (s *stubMessageStore) Get(channelID, messageID discord.Snowflake) (*discord.Message, bool) {
	cm, ok := s.messages[channelID]
	if !ok {
		return nil, false
	}
	m, ok := cm[messageID]
	return m, ok
}
func (s *stubMessageStore) Channel(channelID discord.Snowflake) *collection.Collection[discord.Snowflake, *discord.Message] {
	c := collection.New[discord.Snowflake, *discord.Message]()
	for mid, m := range s.messages[channelID] {
		c.Set(mid, m)
	}
	return c
}
func (s *stubMessageStore) Size() int { return len(s.messages) }

type stubGuildStore struct {
	guilds map[discord.Snowflake]*discord.Guild
}

func (s *stubGuildStore) Get(id discord.Snowflake) (*discord.Guild, bool) {
	g, ok := s.guilds[id]
	return g, ok
}
func (s *stubGuildStore) All() *collection.Collection[discord.Snowflake, *discord.Guild] {
	c := collection.New[discord.Snowflake, *discord.Guild]()
	for id, g := range s.guilds {
		c.Set(id, g)
	}
	return c
}
func (s *stubGuildStore) Size() int { return len(s.guilds) }

type stubUserStore struct {
	users map[discord.Snowflake]*discord.User
}

func (s *stubUserStore) Get(id discord.Snowflake) (*discord.User, bool) {
	u, ok := s.users[id]
	return u, ok
}
func (s *stubUserStore) All() *collection.Collection[discord.Snowflake, *discord.User] {
	c := collection.New[discord.Snowflake, *discord.User]()
	for id, u := range s.users {
		c.Set(id, u)
	}
	return c
}
func (s *stubUserStore) Size() int { return len(s.users) }

// ── Stub cache ────────────────────────────────────────────────────────────────

type stubCache struct {
	members  *stubMemberStore
	roles    *stubRoleStore
	channels *stubChannelStore
	messages *stubMessageStore
	guilds   *stubGuildStore
	users    *stubUserStore
}

func newStubCache() *stubCache {
	return &stubCache{
		members:  &stubMemberStore{members: map[discord.Snowflake]map[discord.Snowflake]*discord.GuildMember{}},
		roles:    &stubRoleStore{roles: map[discord.Snowflake]map[discord.Snowflake]*discord.Role{}},
		channels: &stubChannelStore{channels: map[discord.Snowflake]*discord.Channel{}},
		messages: &stubMessageStore{messages: map[discord.Snowflake]map[discord.Snowflake]*discord.Message{}},
		guilds:   &stubGuildStore{guilds: map[discord.Snowflake]*discord.Guild{}},
		users:    &stubUserStore{users: map[discord.Snowflake]*discord.User{}},
	}
}

func (c *stubCache) Guilds() discord.GuildStore                   { return c.guilds }
func (c *stubCache) Channels() discord.ChannelStore               { return c.channels }
func (c *stubCache) Users() discord.UserStore                     { return c.users }
func (c *stubCache) Members() discord.MemberStore                 { return c.members }
func (c *stubCache) Messages() discord.MessageStore               { return c.messages }
func (c *stubCache) Roles() discord.RoleStore                     { return c.roles }
func (c *stubCache) VoiceStates() discord.VoiceStateStore         { return nil }
func (c *stubCache) Soundboard() discord.SoundboardStore          { return nil }
func (c *stubCache) ScheduledEvents() discord.ScheduledEventStore { return nil }
func (c *stubCache) StageInstances() discord.StageInstanceStore   { return nil }
func (c *stubCache) Emojis() discord.EmojiStore                   { return nil }
func (c *stubCache) Stickers() discord.StickerStore               { return nil }

var _ discord.Cache = (*stubCache)(nil)

// ── Stub EntityClient ─────────────────────────────────────────────────────────

type stubEntityClient struct {
	cache        discord.Cache
	fetchMember  func(ctx context.Context, guildID, userID discord.Snowflake) (*discord.GuildMember, error)
	fetchRole    func(ctx context.Context, guildID, roleID discord.Snowflake) (*discord.Role, error)
	fetchMessage func(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error)
	fetchGuild   func(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error)
	fetchUser    func(ctx context.Context, userID discord.Snowflake) (*discord.User, error)
}

func (s *stubEntityClient) ClientCache() discord.Cache { return s.cache }
func (s *stubEntityClient) GetGuildMember(ctx context.Context, guildID, userID discord.Snowflake) (*discord.GuildMember, error) {
	if s.fetchMember != nil {
		return s.fetchMember(ctx, guildID, userID)
	}
	return nil, nil
}
func (s *stubEntityClient) GetGuildRole(ctx context.Context, guildID, roleID discord.Snowflake) (*discord.Role, error) {
	if s.fetchRole != nil {
		return s.fetchRole(ctx, guildID, roleID)
	}
	return nil, nil
}
func (s *stubEntityClient) GetMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	if s.fetchMessage != nil {
		return s.fetchMessage(ctx, channelID, messageID)
	}
	return nil, nil
}
func (s *stubEntityClient) GetGuild(ctx context.Context, guildID discord.Snowflake) (*discord.Guild, error) {
	if s.fetchGuild != nil {
		return s.fetchGuild(ctx, guildID)
	}
	return nil, nil
}
func (s *stubEntityClient) GetUser(ctx context.Context, userID discord.Snowflake) (*discord.User, error) {
	if s.fetchUser != nil {
		return s.fetchUser(ctx, userID)
	}
	return nil, nil
}

// Unused but required by EntityClient interface.
func (s *stubEntityClient) EditMessage(ctx context.Context, channelID, messageID discord.Snowflake, opts discord.MessageEditOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) CreateMessage(ctx context.Context, channelID discord.Snowflake, opts discord.MessageCreateOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) PinMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) UnpinMessage(ctx context.Context, channelID, messageID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) CrosspostMessage(ctx context.Context, channelID, messageID discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) AddReaction(ctx context.Context, channelID, messageID discord.Snowflake, emoji string) error {
	return nil
}
func (s *stubEntityClient) ModifyChannel(ctx context.Context, channelID discord.Snowflake, opts discord.ChannelEditOptions) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteChannel(ctx context.Context, channelID discord.Snowflake, reason *string) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) BulkDeleteMessages(ctx context.Context, channelID discord.Snowflake, messageIDs []discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) GetChannelMessages(ctx context.Context, channelID discord.Snowflake, opts discord.FetchMessagesOptions) ([]*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) TriggerTypingIndicator(ctx context.Context, channelID discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) SetVoiceChannelStatus(ctx context.Context, channelID discord.Snowflake, status *string) error {
	return nil
}
func (s *stubEntityClient) CreateChannelInvite(ctx context.Context, channelID discord.Snowflake, opts discord.InviteCreateOptions) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuild(ctx context.Context, guildID discord.Snowflake, opts discord.GuildEditOptions) (*discord.Guild, error) {
	return nil, nil
}
func (s *stubEntityClient) LeaveGuild(ctx context.Context, guildID discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) ListGuildMembers(ctx context.Context, guildID discord.Snowflake, opts discord.FetchMembersOptions) ([]*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildRole(ctx context.Context, guildID discord.Snowflake, opts discord.RoleCreateOptions) (*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildChannel(ctx context.Context, guildID discord.Snowflake, opts discord.ChannelCreateOptions) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildEmoji(ctx context.Context, guildID discord.Snowflake, opts discord.EmojiCreateOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildSticker(ctx context.Context, guildID discord.Snowflake, opts discord.StickerCreateOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildAuditLog(ctx context.Context, guildID discord.Snowflake, opts discord.AuditLogOptions) (*discord.AuditLog, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildScheduledEvent(ctx context.Context, guildID discord.Snowflake, opts discord.ScheduledEventCreateOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildMember(ctx context.Context, guildID, userID discord.Snowflake, opts discord.MemberEditOptions) (*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubEntityClient) KickGuildMember(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) CreateGuildBan(ctx context.Context, guildID, userID discord.Snowflake, opts discord.BanOptions) error {
	return nil
}
func (s *stubEntityClient) AddGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) RemoveGuildMemberRole(ctx context.Context, guildID, userID, roleID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, opts discord.RoleEditOptions) (*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildRole(ctx context.Context, guildID, roleID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildRolePositions(ctx context.Context, guildID discord.Snowflake, opts discord.RolePositionOptions) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateDM(ctx context.Context, recipientID discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, opts discord.EmojiEditOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ModifyWebhook(ctx context.Context, webhookID discord.Snowflake, opts discord.WebhookEditOptions) (*discord.Webhook, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteWebhook(ctx context.Context, webhookID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ExecuteWebhook(ctx context.Context, webhookID discord.Snowflake, token string, opts discord.WebhookExecuteOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) GetWebhookMessage(ctx context.Context, webhookID discord.Snowflake, token string, messageID discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteInvite(ctx context.Context, code string, reason *string) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyStageInstance(ctx context.Context, channelID discord.Snowflake, opts discord.StageEditOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteStageInstance(ctx context.Context, channelID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.ScheduledEventEditOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) GetGuildScheduledEventUsers(ctx context.Context, guildID, eventID discord.Snowflake, opts discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, opts discord.StickerEditOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) ModifyAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake, opts discord.RuleEditOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, opts discord.SoundEditOptions) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) DeleteGuildIntegration(ctx context.Context, guildID, integrationID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) GetGuildIntegrations(ctx context.Context, guildID discord.Snowflake) ([]*discord.Integration, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildRoles(ctx context.Context, guildID discord.Snowflake) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) GetChannel(ctx context.Context, channelID discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildChannels(ctx context.Context, guildID discord.Snowflake) ([]*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildEmoji(ctx context.Context, guildID, emojiID discord.Snowflake) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildEmojis(ctx context.Context, guildID discord.Snowflake) ([]*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildSticker(ctx context.Context, guildID, stickerID discord.Snowflake) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildStickers(ctx context.Context, guildID discord.Snowflake) ([]*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildBan(ctx context.Context, guildID, userID discord.Snowflake) (*discord.Ban, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildBans(ctx context.Context, guildID discord.Snowflake, opts discord.FetchBansOptions) ([]*discord.Ban, error) {
	return nil, nil
}
func (s *stubEntityClient) RemoveGuildBan(ctx context.Context, guildID, userID discord.Snowflake, reason *string) error {
	return nil
}
func (s *stubEntityClient) GetGuildScheduledEvent(ctx context.Context, guildID, eventID discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildScheduledEvents(ctx context.Context, guildID discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) GetStageInstance(ctx context.Context, channelID discord.Snowflake) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateStageInstance(ctx context.Context, opts discord.StageCreateOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildSoundboardSound(ctx context.Context, guildID, soundID discord.Snowflake) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildSoundboardSounds(ctx context.Context, guildID discord.Snowflake) ([]*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildInvites(ctx context.Context, guildID discord.Snowflake) ([]*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildWebhooks(ctx context.Context, guildID discord.Snowflake) ([]*discord.Webhook, error) {
	return nil, nil
}
func (s *stubEntityClient) GetAutoModerationRule(ctx context.Context, guildID, ruleID discord.Snowflake) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) ListAutoModerationRules(ctx context.Context, guildID discord.Snowflake) ([]*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateAutoModerationRule(ctx context.Context, guildID discord.Snowflake, opts discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}

var _ discord.EntityClient = (*stubEntityClient)(nil)

// ── MemberManager tests ───────────────────────────────────────────────────────

func TestMemberManager_Cache_ReturnsGuildMembers(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const userID discord.Snowflake = "u1"
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{
		userID: {UserID: userID, GuildID: guildID},
	}
	m := managers.NewMemberManager(guildID, &stubEntityClient{cache: cache})

	coll := m.Cache()
	if coll.Len() != 1 {
		t.Fatalf("expected 1 member in cache, got %d", coll.Len())
	}
}

func TestMemberManager_Get_HitAndMiss(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const userID discord.Snowflake = "u1"
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{
		userID: {UserID: userID},
	}
	m := managers.NewMemberManager(guildID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(userID); !ok {
		t.Fatal("expected Get to return member")
	}
	if _, ok := m.Get("unknown"); ok {
		t.Fatal("expected Get to return false for unknown ID")
	}
}

func TestMemberManager_Get_NoCache(t *testing.T) {
	m := managers.NewMemberManager("g1", &stubEntityClient{cache: nil})
	_, ok := m.Get("u1")
	if ok {
		t.Fatal("expected false when no cache")
	}
}

func TestMemberManager_Fetch_CallsClient(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const userID discord.Snowflake = "u1"
	want := &discord.GuildMember{UserID: userID}
	called := false
	client := &stubEntityClient{
		fetchMember: func(_ context.Context, gid, uid discord.Snowflake) (*discord.GuildMember, error) {
			called = true
			if gid != guildID || uid != userID {
				t.Errorf("unexpected IDs: guild=%s user=%s", gid, uid)
			}
			return want, nil
		},
	}
	m := managers.NewMemberManager(guildID, client)
	got, err := m.Fetch(context.Background(), userID)
	if err != nil || got != want {
		t.Fatalf("Fetch: err=%v got=%v", err, got)
	}
	if !called {
		t.Fatal("GetGuildMember was not called")
	}
}

func TestMemberManager_Resolve(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const userID discord.Snowflake = "u1"
	mem := &discord.GuildMember{UserID: userID}
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{userID: mem}
	m := managers.NewMemberManager(guildID, &stubEntityClient{cache: cache})

	if got := m.Resolve(mem); got != mem {
		t.Fatal("Resolve(*GuildMember) should return it directly")
	}
	if got := m.Resolve(userID); got != mem {
		t.Fatalf("Resolve(Snowflake) should return cached member, got %v", got)
	}
	if got := m.Resolve(string(userID)); got != mem {
		t.Fatalf("Resolve(string) should return cached member, got %v", got)
	}
	if got := m.Resolve(42); got != nil {
		t.Fatal("Resolve with unknown type should return nil")
	}
}

func TestMemberManager_Size_NoCache(t *testing.T) {
	m := managers.NewMemberManager("g1", &stubEntityClient{})
	if s := m.Size(); s != 0 {
		t.Fatalf("expected size 0 with no cache, got %d", s)
	}
}

// ── RoleManager tests ─────────────────────────────────────────────────────────

func TestRoleManager_Cache_ReturnsGuildRoles(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{
		"r1": {ID: "r1"},
		"r2": {ID: "r2"},
	}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	coll := m.Cache()
	if coll.Len() != 2 {
		t.Fatalf("expected 2 roles, got %d", coll.Len())
	}
}

func TestRoleManager_Get(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const roleID discord.Snowflake = "r1"
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{roleID: {ID: roleID}}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(roleID); !ok {
		t.Fatal("expected Get to find role")
	}
}

func TestRoleManager_Resolve(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	const roleID discord.Snowflake = "r1"
	role := &discord.Role{ID: roleID}
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{roleID: role}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	if got := m.Resolve(role); got != role {
		t.Fatal("Resolve(*Role) should return it directly")
	}
	if got := m.Resolve(roleID); got != role {
		t.Fatalf("Resolve(Snowflake) returned %v", got)
	}
	if got := m.Resolve(nil); got != nil {
		t.Fatal("Resolve(nil) should return nil")
	}
}

// ── MessageManager tests ──────────────────────────────────────────────────────

func TestMessageManager_Cache_ReturnsChannelMessages(t *testing.T) {
	const channelID discord.Snowflake = "ch1"
	const msgID discord.Snowflake = "m1"
	cache := newStubCache()
	cache.messages.messages[channelID] = map[discord.Snowflake]*discord.Message{
		msgID: {ID: msgID},
	}
	m := managers.NewMessageManager(channelID, &stubEntityClient{cache: cache})

	coll := m.Cache()
	if coll.Len() != 1 {
		t.Fatalf("expected 1 message, got %d", coll.Len())
	}
}

func TestMessageManager_Get(t *testing.T) {
	const channelID discord.Snowflake = "ch1"
	const msgID discord.Snowflake = "m1"
	cache := newStubCache()
	cache.messages.messages[channelID] = map[discord.Snowflake]*discord.Message{
		msgID: {ID: msgID},
	}
	m := managers.NewMessageManager(channelID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(msgID); !ok {
		t.Fatal("expected Get to find message")
	}
	if _, ok := m.Get("unknown"); ok {
		t.Fatal("expected Get to miss unknown message")
	}
}

func TestMessageManager_Fetch_CallsClient(t *testing.T) {
	const channelID discord.Snowflake = "ch1"
	const msgID discord.Snowflake = "m1"
	want := &discord.Message{ID: msgID}
	client := &stubEntityClient{
		fetchMessage: func(_ context.Context, chID, mID discord.Snowflake) (*discord.Message, error) {
			return want, nil
		},
	}
	m := managers.NewMessageManager(channelID, client)
	got, err := m.Fetch(context.Background(), msgID)
	if err != nil || got != want {
		t.Fatalf("Fetch: err=%v got=%v", err, got)
	}
}

// ── ClientGuildManager tests ──────────────────────────────────────────────────

func TestClientGuildManager_Cache_Empty(t *testing.T) {
	m := managers.NewClientGuildManager(&stubEntityClient{cache: newStubCache()})
	if m.Cache().Len() != 0 {
		t.Fatal("expected empty cache")
	}
}

func TestClientGuildManager_Fetch(t *testing.T) {
	const guildID discord.Snowflake = "g1"
	want := &discord.Guild{ID: guildID}
	client := &stubEntityClient{
		fetchGuild: func(_ context.Context, gid discord.Snowflake) (*discord.Guild, error) {
			return want, nil
		},
	}
	m := managers.NewClientGuildManager(client)
	got, err := m.Fetch(context.Background(), guildID)
	if err != nil || got != want {
		t.Fatalf("Fetch: err=%v got=%v", err, got)
	}
}

func TestClientGuildManager_Resolve_NoCache(t *testing.T) {
	m := managers.NewClientGuildManager(&stubEntityClient{})
	if got := m.Resolve("g1"); got != nil {
		t.Fatal("Resolve with no cache should return nil")
	}
}

// ── ClientUserManager tests ───────────────────────────────────────────────────

func TestClientUserManager_Fetch(t *testing.T) {
	const userID discord.Snowflake = "u1"
	want := &discord.User{ID: userID}
	client := &stubEntityClient{
		fetchUser: func(_ context.Context, uid discord.Snowflake) (*discord.User, error) {
			return want, nil
		},
	}
	m := managers.NewClientUserManager(client)
	got, err := m.Fetch(context.Background(), userID)
	if err != nil || got != want {
		t.Fatalf("Fetch: err=%v got=%v", err, got)
	}
}

// ── Nil cache guard ───────────────────────────────────────────────────────────

func TestMemberManager_CacheReturnsEmpty_WhenNilCache(t *testing.T) {
	m := managers.NewMemberManager("g1", &stubEntityClient{cache: nil})
	coll := m.Cache()
	if coll == nil {
		t.Fatal("Cache() should return non-nil empty collection when cache is nil")
	}
	if coll.Len() != 0 {
		t.Fatal("Cache() should return empty collection when cache is nil")
	}
}
