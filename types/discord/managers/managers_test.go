package managers_test

import (
	"context"
	"strconv"
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
func (s *stubEntityClient) EditMessage(_ context.Context, _, _ discord.Snowflake, _ discord.MessageEditOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteMessage(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) CreateMessage(_ context.Context, _ discord.Snowflake, _ discord.MessageCreateOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) PinMessage(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) UnpinMessage(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) CrosspostMessage(_ context.Context, _, _ discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) AddReaction(_ context.Context, _, _ discord.Snowflake, _ string) error {
	return nil
}
func (s *stubEntityClient) ModifyChannel(_ context.Context, _ discord.Snowflake, _ discord.ChannelEditOptions) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteChannel(_ context.Context, _ discord.Snowflake, _ *string) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) BulkDeleteMessages(_ context.Context, _ discord.Snowflake, _ []discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) GetChannelMessages(_ context.Context, _ discord.Snowflake, _ discord.FetchMessagesOptions) ([]*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) TriggerTypingIndicator(_ context.Context, _ discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) SetVoiceChannelStatus(_ context.Context, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) CreateChannelInvite(_ context.Context, _ discord.Snowflake, _ discord.InviteCreateOptions) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuild(_ context.Context, _ discord.Snowflake, _ discord.GuildEditOptions) (*discord.Guild, error) {
	return nil, nil
}
func (s *stubEntityClient) LeaveGuild(_ context.Context, _ discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) ListGuildMembers(_ context.Context, _ discord.Snowflake, _ discord.FetchMembersOptions) ([]*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildRole(_ context.Context, _ discord.Snowflake, _ discord.RoleCreateOptions) (*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildChannel(_ context.Context, _ discord.Snowflake, _ discord.ChannelCreateOptions) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildEmoji(_ context.Context, _ discord.Snowflake, _ discord.EmojiCreateOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildSticker(_ context.Context, _ discord.Snowflake, _ discord.StickerCreateOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildAuditLog(_ context.Context, _ discord.Snowflake, _ discord.AuditLogOptions) (*discord.AuditLog, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateGuildScheduledEvent(_ context.Context, _ discord.Snowflake, _ discord.ScheduledEventCreateOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildMember(_ context.Context, _, _ discord.Snowflake, _ discord.MemberEditOptions) (*discord.GuildMember, error) {
	return nil, nil
}
func (s *stubEntityClient) KickGuildMember(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) CreateGuildBan(_ context.Context, _, _ discord.Snowflake, _ discord.BanOptions) error {
	return nil
}
func (s *stubEntityClient) AddGuildMemberRole(_ context.Context, _, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) RemoveGuildMemberRole(_ context.Context, _, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildRole(_ context.Context, _, _ discord.Snowflake, _ discord.RoleEditOptions) (*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildRole(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildRolePositions(_ context.Context, _ discord.Snowflake, _ discord.RolePositionOptions) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateDM(_ context.Context, _ discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildEmoji(_ context.Context, _, _ discord.Snowflake, _ discord.EmojiEditOptions) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildEmoji(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ModifyWebhook(_ context.Context, _ discord.Snowflake, _ discord.WebhookEditOptions) (*discord.Webhook, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteWebhook(_ context.Context, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ExecuteWebhook(_ context.Context, _ discord.Snowflake, _ string, _ discord.WebhookExecuteOptions) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) GetWebhookMessage(_ context.Context, _ discord.Snowflake, _ string, _ discord.Snowflake) (*discord.Message, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteInvite(_ context.Context, _ string, _ *string) (*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyStageInstance(_ context.Context, _ discord.Snowflake, _ discord.StageEditOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteStageInstance(_ context.Context, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildScheduledEvent(_ context.Context, _, _ discord.Snowflake, _ discord.ScheduledEventEditOptions) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildScheduledEvent(_ context.Context, _, _ discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) GetGuildScheduledEventUsers(_ context.Context, _, _ discord.Snowflake, _ discord.FetchUsersOptions) ([]*discord.GuildScheduledEventUser, error) {
	return nil, nil
}
func (s *stubEntityClient) ModifyGuildSticker(_ context.Context, _, _ discord.Snowflake, _ discord.StickerEditOptions) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildSticker(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) ModifyAutoModerationRule(_ context.Context, _, _ discord.Snowflake, _ discord.RuleEditOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteAutoModerationRule(_ context.Context, _, _ discord.Snowflake) error {
	return nil
}
func (s *stubEntityClient) ModifyGuildSoundboardSound(_ context.Context, _, _ discord.Snowflake, _ discord.SoundEditOptions) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) DeleteGuildSoundboardSound(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) DeleteGuildIntegration(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) GetGuildIntegrations(_ context.Context, _ discord.Snowflake) ([]*discord.Integration, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildRoles(_ context.Context, _ discord.Snowflake) ([]*discord.Role, error) {
	return nil, nil
}
func (s *stubEntityClient) GetChannel(_ context.Context, _ discord.Snowflake) (*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildChannels(_ context.Context, _ discord.Snowflake) ([]*discord.Channel, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildEmoji(_ context.Context, _, _ discord.Snowflake) (*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildEmojis(_ context.Context, _ discord.Snowflake) ([]*discord.Emoji, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildSticker(_ context.Context, _, _ discord.Snowflake) (*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildStickers(_ context.Context, _ discord.Snowflake) ([]*discord.Sticker, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildBan(_ context.Context, _, _ discord.Snowflake) (*discord.Ban, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildBans(_ context.Context, _ discord.Snowflake, _ discord.FetchBansOptions) ([]*discord.Ban, error) {
	return nil, nil
}
func (s *stubEntityClient) RemoveGuildBan(_ context.Context, _, _ discord.Snowflake, _ *string) error {
	return nil
}
func (s *stubEntityClient) GetGuildScheduledEvent(_ context.Context, _, _ discord.Snowflake) (*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildScheduledEvents(_ context.Context, _ discord.Snowflake) ([]*discord.GuildScheduledEvent, error) {
	return nil, nil
}
func (s *stubEntityClient) GetStageInstance(_ context.Context, _ discord.Snowflake) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateStageInstance(_ context.Context, _ discord.StageCreateOptions) (*discord.StageInstance, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildSoundboardSound(_ context.Context, _, _ discord.Snowflake) (*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) ListGuildSoundboardSounds(_ context.Context, _ discord.Snowflake) ([]*discord.SoundboardSound, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildInvites(_ context.Context, _ discord.Snowflake) ([]*discord.Invite, error) {
	return nil, nil
}
func (s *stubEntityClient) GetGuildWebhooks(_ context.Context, _ discord.Snowflake) ([]*discord.Webhook, error) {
	return nil, nil
}
func (s *stubEntityClient) GetAutoModerationRule(_ context.Context, _, _ discord.Snowflake) (*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) ListAutoModerationRules(_ context.Context, _ discord.Snowflake) ([]*discord.AutoModerationRule, error) {
	return nil, nil
}
func (s *stubEntityClient) CreateAutoModerationRule(_ context.Context, _ discord.Snowflake, _ discord.RuleCreateOptions) (*discord.AutoModerationRule, error) {
	return nil, nil
}

var _ discord.EntityClient = (*stubEntityClient)(nil)

// ── MemberManager tests ───────────────────────────────────────────────────────

func TestMemberManager_Cache_ReturnsGuildMembers(t *testing.T) {
	const guildID discord.Snowflake = 1
	const userID discord.Snowflake = 2
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
	const guildID discord.Snowflake = 1
	const userID discord.Snowflake = 2
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{
		userID: {UserID: userID},
	}
	m := managers.NewMemberManager(guildID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(userID); !ok {
		t.Fatal("expected Get to return member")
	}
	if _, ok := m.Get(123); ok {
		t.Fatal("expected Get to return false for unknown ID")
	}
}

func TestMemberManager_Get_NoCache(t *testing.T) {
	m := managers.NewMemberManager(1, &stubEntityClient{cache: nil})
	_, ok := m.Get(2)
	if ok {
		t.Fatal("expected false when no cache")
	}
}

func TestMemberManager_Fetch_CallsClient(t *testing.T) {
	const guildID discord.Snowflake = 1
	const userID discord.Snowflake = 2
	want := &discord.GuildMember{UserID: userID}
	called := false
	client := &stubEntityClient{
		fetchMember: func(_ context.Context, gid, uid discord.Snowflake) (*discord.GuildMember, error) {
			called = true
			if gid != guildID || uid != userID {
				t.Errorf("unexpected IDs: guild=%s user=%s", &gid, &uid)
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
	const guildID discord.Snowflake = 213456789012345
	const userID discord.Snowflake = 123456789012345
	mem := &discord.GuildMember{UserID: userID}
	cache := newStubCache()
	cache.members.members[guildID] = map[discord.Snowflake]*discord.GuildMember{userID: mem}
	m := managers.NewMemberManager(guildID, &stubEntityClient{cache: cache})

	if got, _ := m.Resolve(mem); got != mem {
		t.Fatal("Resolve(*GuildMember) should return it directly")
	}
	if got, _ := m.Resolve(userID); got != mem {
		t.Fatalf("Resolve(Snowflake) should return cached member, got %v", got)
	}
	if got, _ := m.Resolve(strconv.FormatUint(uint64(userID), 10)); got != mem {
		t.Fatalf("Resolve(string) should return cached member, got %v", got)
	}
	if got, _ := m.Resolve(42); got != nil {
		t.Fatal("Resolve with unknown type should return nil")
	}
}

func TestMemberManager_Size_NoCache(t *testing.T) {
	m := managers.NewMemberManager(1, &stubEntityClient{})
	if s := m.Size(); s != 0 {
		t.Fatalf("expected size 0 with no cache, got %d", s)
	}
}

// ── RoleManager tests ─────────────────────────────────────────────────────────

func TestRoleManager_Cache_ReturnsGuildRoles(t *testing.T) {
	const guildID discord.Snowflake = 1
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{
		2: {ID: 2},
		3: {ID: 3},
	}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	coll := m.Cache()
	if coll.Len() != 2 {
		t.Fatalf("expected 2 roles, got %d", coll.Len())
	}
}

func TestRoleManager_Get(t *testing.T) {
	const guildID discord.Snowflake = 1
	const roleID discord.Snowflake = 2
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{roleID: {ID: roleID}}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(roleID); !ok {
		t.Fatal("expected Get to find role")
	}
}

func TestRoleManager_Resolve(t *testing.T) {
	const guildID discord.Snowflake = 10
	const roleID discord.Snowflake = 11
	role := &discord.Role{ID: roleID}
	cache := newStubCache()
	cache.roles.roles[guildID] = map[discord.Snowflake]*discord.Role{roleID: role}
	m := managers.NewRoleManager(guildID, &stubEntityClient{cache: cache})

	if got, _ := m.Resolve(role); got != role {
		t.Fatal("Resolve(*Role) should return it directly")
	}
	if got, _ := m.Resolve(roleID); got != role {
		t.Fatalf("Resolve(Snowflake) returned %v", got)
	}
	if got, _ := m.Resolve(nil); got != nil {
		t.Fatal("Resolve(nil) should return nil")
	}
}

// ── MessageManager tests ──────────────────────────────────────────────────────

func TestMessageManager_Cache_ReturnsChannelMessages(t *testing.T) {
	const channelID discord.Snowflake = 12
	const msgID discord.Snowflake = 13
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
	const channelID discord.Snowflake = 12
	const msgID discord.Snowflake = 13
	cache := newStubCache()
	cache.messages.messages[channelID] = map[discord.Snowflake]*discord.Message{
		msgID: {ID: msgID},
	}
	m := managers.NewMessageManager(channelID, &stubEntityClient{cache: cache})

	if _, ok := m.Get(msgID); !ok {
		t.Fatal("expected Get to find message")
	}
	if _, ok := m.Get(123); ok {
		t.Fatal("expected Get to miss unknown message")
	}
}

func TestMessageManager_Fetch_CallsClient(t *testing.T) {
	const channelID discord.Snowflake = 12
	const msgID discord.Snowflake = 13
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
	const guildID discord.Snowflake = 10
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
	if got, _ := m.Resolve(10); got != nil {
		t.Fatal("Resolve with no cache should return nil")
	}
}

// ── ClientUserManager tests ───────────────────────────────────────────────────

func TestClientUserManager_Fetch(t *testing.T) {
	const userID discord.Snowflake = 13123
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
	m := managers.NewMemberManager(10, &stubEntityClient{cache: nil})
	coll := m.Cache()
	if coll == nil {
		t.Fatal("Cache() should return non-nil empty collection when cache is nil")
	}
	if coll.Len() != 0 {
		t.Fatal("Cache() should return empty collection when cache is nil")
	}
}
