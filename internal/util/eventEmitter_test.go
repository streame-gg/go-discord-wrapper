package util_test

import (
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/internal/util"
)

func (su *eventEmitterSuite) TestEventEmitter_OnReturnsIncrementingIDs() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()

	id1 := e.On("evt", func() {})
	id2 := e.On("evt", func() {})
	id3 := e.On("other", func() {})

	if id1 == id2 || id2 == id3 || id1 == id3 {
		t.Fatalf("expected unique ids, got %d, %d, %d", id1, id2, id3)
	}
}

func (su *eventEmitterSuite) TestEventEmitter_HandlersEmptyEventReturnsNil() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()

	if got := e.Handlers("missing"); got != nil {
		t.Fatalf("expected nil for unknown event, got %#v", got)
	}
}

func (su *eventEmitterSuite) TestEventEmitter_HandlersFastPathNoOnce() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()
	e.On("evt", func() {})
	e.On("evt", func() {})

	// First call returns both handlers.
	if got := e.Handlers("evt"); len(got) != 2 {
		t.Fatalf("expected 2 handlers, got %d", len(got))
	}
	// Persistent handlers remain after retrieval.
	if got := e.Handlers("evt"); len(got) != 2 {
		t.Fatalf("expected handlers to persist, got %d", len(got))
	}
}

func (su *eventEmitterSuite) TestEventEmitter_OnceFiresExactlyOnce() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()
	e.Once("evt", func() {})

	if got := e.Handlers("evt"); len(got) != 1 {
		t.Fatalf("expected once-handler returned, got %d", len(got))
	}
	// Once-handler is removed after first retrieval (slow path).
	if got := e.Handlers("evt"); got != nil {
		t.Fatalf("expected once-handler removed, got %#v", got)
	}
}

func (su *eventEmitterSuite) TestEventEmitter_OnceMixedWithPersistent() {
	t := su.T()
	e := util.NewEventEmitter[string, int]()
	e.On("evt", 1)
	e.Once("evt", 2)
	e.On("evt", 3)

	first := e.Handlers("evt")
	if len(first) != 3 {
		t.Fatalf("expected 3 handlers on first call, got %d", len(first))
	}

	// After the once-handler fires, only the two persistent handlers remain.
	second := e.Handlers("evt")
	if len(second) != 2 {
		t.Fatalf("expected 2 handlers after once removed, got %d", len(second))
	}
	for _, v := range second {
		if v == 2 {
			t.Fatalf("once-handler (2) should have been removed, got %v", second)
		}
	}
}

func (su *eventEmitterSuite) TestEventEmitter_OffRemovesHandler() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()
	id := e.On("evt", func() {})

	if !e.Off("evt", id) {
		t.Fatal("expected Off to report removal")
	}
	if got := e.Handlers("evt"); got != nil {
		t.Fatalf("expected no handlers after Off, got %#v", got)
	}
}

func (su *eventEmitterSuite) TestEventEmitter_OffUnknownIDReturnsFalse() {
	t := su.T()
	e := util.NewEventEmitter[string, func()]()
	e.On("evt", func() {})

	if e.Off("evt", 999) {
		t.Fatal("expected Off to return false for unknown id")
	}
	if e.Off("missing", 0) {
		t.Fatal("expected Off to return false for unknown event")
	}
}

func (su *eventEmitterSuite) TestEventEmitter_OffRemovesCorrectHandlerAmongMany() {
	t := su.T()
	e := util.NewEventEmitter[string, int]()
	id1 := e.On("evt", 10)
	id2 := e.On("evt", 20)
	id3 := e.On("evt", 30)

	if !e.Off("evt", id2) {
		t.Fatal("expected removal of middle handler")
	}

	got := e.Handlers("evt")
	if len(got) != 2 {
		t.Fatalf("expected 2 remaining, got %d", len(got))
	}
	for _, v := range got {
		if v == 20 {
			t.Fatal("handler 20 should have been removed")
		}
	}
	_ = id1
	_ = id3
}

func (su *eventEmitterSuite) TestEventEmitter_ConcurrentAccessIsSafe() {
	e := util.NewEventEmitter[string, func()]()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := e.On("evt", func() {})
			e.Once("evt", func() {})
			_ = e.Handlers("evt")
			e.Off("evt", id)
		}()
	}
	wg.Wait()
}

type eventEmitterSuite struct{ suite.Suite }

func TestEventEmitterSuite(t *testing.T) { suite.Run(t, new(eventEmitterSuite)) }
