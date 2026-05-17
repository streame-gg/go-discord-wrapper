package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func applyReason(reason string) *http.Request {
	req, _ := http.NewRequest(http.MethodPatch, "https://discord.com/api/v10/test", nil)
	WithAuditLogReason(reason)(req)
	return req
}

func TestWithAuditLogReason(t *testing.T) {
	t.Run("short reason sets header", func(t *testing.T) {
		req := applyReason("ban evasion")
		assert.NotEmpty(t, req.Header.Get("X-Audit-Log-Reason"))
	})

	t.Run("empty reason sets no header", func(t *testing.T) {
		req := applyReason("")
		assert.Empty(t, req.Header.Get("X-Audit-Log-Reason"))
	})

	t.Run("600-char reason truncated to 512 chars", func(t *testing.T) {
		long := strings.Repeat("a", 600)
		req := applyReason(long)
		// URL-decoded length must be ≤ 512 runes.
		got := req.Header.Get("X-Audit-Log-Reason")
		assert.NotEmpty(t, got)
		// All ASCII — no percent-encoding, so byte length == rune count.
		assert.LessOrEqual(t, len([]rune(got)), 512)
	})

	t.Run("reason with Umlaut is percent-encoded", func(t *testing.T) {
		req := applyReason("Verstoß")
		got := req.Header.Get("X-Audit-Log-Reason")
		assert.Contains(t, got, "%C3%9F", "ß should be percent-encoded as %%C3%%9F")
	})

	t.Run("reason with emoji is percent-encoded", func(t *testing.T) {
		req := applyReason("spam 🚫")
		got := req.Header.Get("X-Audit-Log-Reason")
		// Any percent-encoding present means non-ASCII was encoded.
		assert.Contains(t, got, "%", "emoji should be percent-encoded")
	})
}
