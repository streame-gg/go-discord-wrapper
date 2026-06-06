package stores

import (
	"github.com/stretchr/testify/suite"
	"sync/atomic"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

// TestBug92_RoleSet_IndexAfterBase verifies that a concurrent GetByGuild call
// never observes a role ID in the index that does not yet exist in the base store.
//
// Root cause: index.add was called BEFORE base.Set, creating a brief window where
// a reader could find an ID via s.index.forGuild but then miss it in s.base.Get.
// Fix: base.Set first, then index.add — any ID visible in the index is guaranteed
// to be present in base.
func (su *storesRoleSuite) TestBug92_RoleSet_IndexAfterBase() {
	t := su.T()
	s := NewRoleStore(StoreOptions{})
	defer s.Close()

	guildID := discord.Snowflake(1)
	const n = 5000

	var violations atomic.Int64
	done := make(chan struct{})

	// Reader: take a raw snapshot of the index and then probe base for each ID.
	// With the old ordering (index first), violations > 0 would be observed.
	var readerDone = make(chan struct{})
	go func() {
		defer close(readerDone)
		for {
			select {
			case <-done:
				return
			default:
			}
			for _, id := range s.index.forGuild(guildID) {
				if _, ok := s.base.Get(id); !ok {
					violations.Add(1)
				}
			}
		}
	}()

	for i := 0; i < n; i++ {
		s.Set(guildID, &discord.Role{ID: discord.Snowflake(i + 1)})
	}
	close(done)
	<-readerDone

	if v := violations.Load(); v > 0 {
		t.Errorf("Bug92: %d observation(s) where role ID was visible in index before base.Set completed", v)
	}
}

type storesRoleSuite struct{ suite.Suite }

func TestStoresRoleSuite(t *testing.T) { suite.Run(t, new(storesRoleSuite)) }
