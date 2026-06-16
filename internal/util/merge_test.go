package util_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// --- helpers -----------------------------------------------------------------

type inner struct {
	X int
	Y *string
}

type flat struct {
	Name    string
	Count   int
	Flag    bool
	Ptr     *string
	Nested  inner
	Slice   []int
	created time.Time // unexported — must never be touched
}

func strPtr(s string) *string { return &s }

// --- scalar fields -----------------------------------------------------------

func (su *mergeSuite) TestMergePartial_EmptyPartialKeepsOld() {
	t := su.T()
	old := flat{Name: "hello", Count: 5, Flag: true, Ptr: strPtr("x")}
	result := util.MergePartial(old, flat{})

	if result.Name != "hello" {
		t.Errorf("Name: got %q, want %q", result.Name, "hello")
	}
	if result.Count != 5 {
		t.Errorf("Count: got %d, want %d", result.Count, 5)
	}
	if !result.Flag {
		t.Error("Flag: got false, want true")
	}
	if result.Ptr == nil || *result.Ptr != "x" {
		t.Errorf("Ptr: got %v, want \"x\"", result.Ptr)
	}
}

func (su *mergeSuite) TestMergePartial_NonZeroScalarApplied() {
	t := su.T()
	old := flat{Name: "old", Count: 1}
	result := util.MergePartial(old, flat{Name: "new", Count: 99})

	if result.Name != "new" {
		t.Errorf("Name: got %q, want %q", result.Name, "new")
	}
	if result.Count != 99 {
		t.Errorf("Count: got %d, want %d", result.Count, 99)
	}
}

func (su *mergeSuite) TestMergePartial_ZeroScalarKeepsOld() {
	t := su.T()
	old := flat{Count: 42}
	result := util.MergePartial(old, flat{Name: "only-name"})

	if result.Count != 42 {
		t.Errorf("Count: got %d, want %d", result.Count, 42)
	}
}

// --- pointer fields ----------------------------------------------------------

func (su *mergeSuite) TestMergePartial_NonNilPtrApplied() {
	t := su.T()
	old := flat{Ptr: strPtr("old")}
	result := util.MergePartial(old, flat{Ptr: strPtr("new")})

	if result.Ptr == nil || *result.Ptr != "new" {
		t.Errorf("Ptr: got %v, want \"new\"", result.Ptr)
	}
}

func (su *mergeSuite) TestMergePartial_NilPtrKeepsOld() {
	t := su.T()
	old := flat{Ptr: strPtr("keep")}
	result := util.MergePartial(old, flat{})

	if result.Ptr == nil || *result.Ptr != "keep" {
		t.Errorf("Ptr: got %v, want \"keep\"", result.Ptr)
	}
}

// --- slice fields ------------------------------------------------------------

func (su *mergeSuite) TestMergePartial_NonNilSliceApplied() {
	t := su.T()
	old := flat{Slice: []int{1, 2, 3}}
	result := util.MergePartial(old, flat{Slice: []int{9}})

	if len(result.Slice) != 1 || result.Slice[0] != 9 {
		t.Errorf("Slice: got %v, want [9]", result.Slice)
	}
}

func (su *mergeSuite) TestMergePartial_NilSliceKeepsOld() {
	t := su.T()
	old := flat{Slice: []int{1, 2, 3}}
	result := util.MergePartial(old, flat{})

	if len(result.Slice) != 3 {
		t.Errorf("Slice: got %v, want [1 2 3]", result.Slice)
	}
}

func (su *mergeSuite) TestMergePartial_EmptySliceApplied() {
	t := su.T()
	// An explicit empty (non-nil) slice overwrites the old slice.
	old := flat{Slice: []int{1, 2, 3}}
	result := util.MergePartial(old, flat{Slice: []int{}})

	if result.Slice == nil || len(result.Slice) != 0 {
		t.Errorf("Slice: got %v, want []", result.Slice)
	}
}

// --- nested struct -----------------------------------------------------------

func (su *mergeSuite) TestMergePartial_NestedStructMergedRecursively() {
	t := su.T()
	old := flat{Nested: inner{X: 10, Y: strPtr("deep")}}
	result := util.MergePartial(old, flat{Nested: inner{X: 20}})

	if result.Nested.X != 20 {
		t.Errorf("Nested.X: got %d, want 20", result.Nested.X)
	}
	if result.Nested.Y == nil || *result.Nested.Y != "deep" {
		t.Errorf("Nested.Y: got %v, want \"deep\"", result.Nested.Y)
	}
}

// --- does not mutate inputs --------------------------------------------------

func (su *mergeSuite) TestMergePartial_DoesNotMutateOld() {
	t := su.T()
	old := flat{Name: "original"}
	_ = util.MergePartial(old, flat{Name: "changed"})

	if old.Name != "original" {
		t.Errorf("old was mutated: got %q, want \"original\"", old.Name)
	}
}

// --- discord.Message ---------------------------------------------------------

func (su *mergeSuite) TestMergePartial_EmptyMessageKeepsOld() {
	t := su.T()
	ts := time.Now()
	old := discord.Message{
		ID:        123,
		ChannelID: 456,
		Content:   "hello world",
		Timestamp: &ts,
		Author:    &discord.User{ID: 789},
		Embeds:    []discord.Embed{{Title: "embed"}},
	}

	result := util.MergePartial(old, discord.Message{})

	if result.ID != 123 {
		t.Errorf("ID: got %d, want 123", result.ID)
	}
	if result.ChannelID != 456 {
		t.Errorf("ChannelID: got %d, want 456", result.ChannelID)
	}
	if result.Content != "hello world" {
		t.Errorf("Content: got %q, want %q", result.Content, "hello world")
	}
	if result.Timestamp == nil {
		t.Error("Timestamp: got nil, want non-nil")
	}
	if result.Author == nil || result.Author.ID != 789 {
		t.Errorf("Author: got %v, want ID=789", result.Author)
	}
	if len(result.Embeds) != 1 {
		t.Errorf("Embeds: got %v, want 1 embed", result.Embeds)
	}
}

func (su *mergeSuite) TestMergePartial_OldMessageIsEmpty() {
	t := su.T()
	ts := time.Now()
	newMsg := discord.Message{
		ID:        123,
		ChannelID: 456,
		Content:   "hello world",
		Timestamp: &ts,
		Author:    &discord.User{ID: 789},
		Embeds:    []discord.Embed{{Title: "embed"}},
	}

	result := util.MergePartial(discord.Message{}, newMsg)

	if result.ID != 123 {
		t.Errorf("ID: got %d, want 123", result.ID)
	}
	if result.ChannelID != 456 {
		t.Errorf("ChannelID: got %d, want 456", result.ChannelID)
	}
	if result.Content != "hello world" {
		t.Errorf("Content: got %q, want %q", result.Content, "hello world")
	}
	if result.Timestamp == nil {
		t.Error("Timestamp: got nil, want non-nil")
	}
	if result.Author == nil || result.Author.ID != 789 {
		t.Errorf("Author: got %v, want ID=789", result.Author)
	}
	if len(result.Embeds) != 1 {
		t.Errorf("Embeds: got %v, want 1 embed", result.Embeds)
	}
}

func (su *mergeSuite) TestMergePartial_PartialMessageUpdatesContent() {
	t := su.T()
	ts := time.Now()
	old := discord.Message{
		ID:        123,
		ChannelID: 456,
		Content:   "original",
		Timestamp: &ts,
		Author:    &discord.User{ID: 789},
	}
	// Discord partial: only ID, ChannelID, and Content are present.
	partial := discord.Message{
		ID:        123,
		ChannelID: 456,
		Content:   "edited",
	}

	result := util.MergePartial(old, partial)

	if result.Content != "edited" {
		t.Errorf("Content: got %q, want %q", result.Content, "edited")
	}
	if result.Author == nil || result.Author.ID != 789 {
		t.Errorf("Author: got %v, want ID=789 (preserved from old)", result.Author)
	}
	if result.Timestamp == nil {
		t.Error("Timestamp: got nil, want preserved from old")
	}
}

type mergeSuite struct{ suite.Suite }

func TestMergeSuite(t *testing.T) { suite.Run(t, new(mergeSuite)) }

// withTime exercises the atomic value-struct path: time.Time has only
// unexported fields, so a field-by-field recursive merge could never set it.
type withTime struct {
	Name string
	When time.Time
}

// TestMergePartial_TimeTimeReplacedWhenNonZero is a regression test: a non-zero
// time.Time in the partial must replace a zero/older destination value. Before
// the fix the recursive merge silently dropped it (unexported fields).
func (s *mergeSuite) TestMergePartial_TimeTimeReplacedWhenNonZero() {
	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	// Zero destination, non-zero source → source wins.
	got := util.MergePartial(withTime{Name: "old"}, withTime{When: when})
	s.True(got.When.Equal(when), "non-zero time.Time must replace a zero destination")
	s.Equal("old", got.Name, "zero-valued source scalar must not overwrite destination")
}

// TestMergePartial_TimeTimeZeroSourceKeepsDestination verifies the dual: a zero
// time.Time in the partial leaves the destination untouched.
func (s *mergeSuite) TestMergePartial_TimeTimeZeroSourceKeepsDestination() {
	when := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	got := util.MergePartial(withTime{Name: "keep", When: when}, withTime{Name: "new"})
	s.True(got.When.Equal(when), "zero-valued source time.Time must not clear the destination")
	s.Equal("new", got.Name)
}

// TestMergePartial_GuildMemberJoinedAt covers the concrete payload type the fix
// targets: GuildMember.JoinedAt is a value time.Time.
func (s *mergeSuite) TestMergePartial_GuildMemberJoinedAt() {
	joined := time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC)
	old := discord.GuildMember{Nick: ptr("old-nick")}
	partial := discord.GuildMember{JoinedAt: joined}

	got := util.MergePartial(old, partial)
	s.True(got.JoinedAt.Equal(joined), "JoinedAt must be applied from the partial")
	s.Require().NotNil(got.Nick)
	s.Equal("old-nick", *got.Nick, "fields absent from the partial must be preserved")
}

func ptr[T any](v T) *T { return &v }

// --- MergePartialJSON: presence-aware merging --------------------------------

// TestMergePartialJSON_PresentZeroOverwrites is the core fix: a field that is
// present in the raw payload with an intentional zero value ("" / false / 0)
// must overwrite the old value — something value-based MergePartial cannot do.
func (s *mergeSuite) TestMergePartialJSON_PresentZeroOverwrites() {
	old := discord.Message{ID: 1, ChannelID: 2, Content: "original", Pinned: true}
	partial := discord.Message{ID: 1, ChannelID: 2, Content: "", Pinned: false}
	raw := json.RawMessage(`{"id":"1","channel_id":"2","content":"","pinned":false}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Equal("", got.Content, "content present as empty string must clear the old value")
	s.False(got.Pinned, "pinned present as false must clear the old value")
}

// TestMergePartialJSON_AbsentFieldKeepsOld verifies that a zero-valued field
// which is NOT present in the payload is preserved from old.
func (s *mergeSuite) TestMergePartialJSON_AbsentFieldKeepsOld() {
	ts := time.Now()
	old := discord.Message{ID: 1, ChannelID: 2, Content: "keep me", Timestamp: &ts}
	// Only content is present; it changed. Timestamp/Author etc. are absent.
	partial := discord.Message{ID: 1, ChannelID: 2, Content: "edited"}
	raw := json.RawMessage(`{"id":"1","channel_id":"2","content":"edited"}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Equal("edited", got.Content)
	s.Require().NotNil(got.Timestamp, "timestamp absent from payload must be preserved")
	s.True(got.Timestamp.Equal(ts))
}

// TestMergePartialJSON_PresentNullClearsPointer verifies an explicit JSON null
// on a pointer field clears it (the partial decodes that field to nil).
func (s *mergeSuite) TestMergePartialJSON_PresentNullClearsPointer() {
	old := discord.GuildMember{Nick: ptr("old-nick")}
	partial := discord.GuildMember{Nick: nil}
	raw := json.RawMessage(`{"nick":null}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Nil(got.Nick, "nick present as null must clear the old pointer")
}

// TestMergePartialJSON_DashTaggedFieldNeverSourced confirms json:"-" fields
// (e.g. GuildMember.GuildID) are never taken from the payload, even if a key of
// the same name appears.
func (s *mergeSuite) TestMergePartialJSON_DashTaggedFieldNeverSourced() {
	old := discord.GuildMember{GuildID: 555}
	partial := discord.GuildMember{GuildID: 999, Nick: ptr("n")}
	raw := json.RawMessage(`{"guild_id":"999","nick":"n"}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Equal(discord.Snowflake(555), got.GuildID, "json:\"-\" field must retain old value")
	s.Require().NotNil(got.Nick)
	s.Equal("n", *got.Nick)
}

// TestMergePartialJSON_FallsBackOnNonObject verifies that a non-object raw
// payload falls back to value-based MergePartial (zero values do not overwrite).
func (s *mergeSuite) TestMergePartialJSON_FallsBackOnNonObject() {
	old := flat{Name: "keep", Count: 7}
	partial := flat{Name: "", Count: 0}

	// Empty, null, and array payloads all trigger the fallback.
	for _, raw := range []json.RawMessage{nil, json.RawMessage(`null`), json.RawMessage(`[1,2]`), json.RawMessage(`oops`)} {
		got := util.MergePartialJSON(old, partial, raw)
		s.Equal("keep", got.Name, "fallback must preserve old on zero-valued partial")
		s.Equal(7, got.Count)
	}
}

// TestMergePartialJSON_NoJSONTagUsesFieldName covers the json-tag-less branch:
// encoding/json matches such fields by (case-insensitive) field name.
func (s *mergeSuite) TestMergePartialJSON_NoJSONTagUsesFieldName() {
	old := flat{Name: "old", Count: 5}
	partial := flat{Name: "new", Count: 0}
	// "count" present as 0 must overwrite; "Name" matched case-insensitively.
	raw := json.RawMessage(`{"name":"new","count":0}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Equal("new", got.Name)
	s.Equal(0, got.Count, "count present as zero must overwrite via field-name match")
}

// TestMergePartialJSON_UnexportedFieldUntouched guards the CanSet check.
func (s *mergeSuite) TestMergePartialJSON_UnexportedFieldUntouched() {
	when := time.Now()
	old := flat{Name: "x", created: when}
	partial := flat{Name: "y"}
	raw := json.RawMessage(`{"name":"y","created":"whatever"}`)

	got := util.MergePartialJSON(old, partial, raw)
	s.Equal("y", got.Name)
	s.True(got.created.Equal(when), "unexported field must never be set from the payload")
}
