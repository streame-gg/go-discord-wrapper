package api

import (
	"mime"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func TestBasenameFilename(t *testing.T) {
	cases := []struct{ in, want string }{
		{"report.pdf", "report.pdf"},
		{"/tmp/upload/report.pdf", "report.pdf"},
		{`C:\Users\foo\report.pdf`, "report.pdf"},
		{"/path/with/trailing/", "file"},
		{"file.txt", "file.txt"},
		{"", "file"},
		{"./relative.pdf", "relative.pdf"},
		{"../escape.pdf", "escape.pdf"},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := basenameFilename(tc.in)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuildMultipartMessage_BasenameInHeader(t *testing.T) {
	payload := []byte(`{"content":"hello"}`)
	files := []discord.MessageFile{
		{Name: "/home/user/Downloads/secret_report.pdf", ContentType: "application/pdf", Data: []byte("pdfdata")},
	}

	buf, ct, err := buildMultipartMessage(payload, files)
	if err != nil {
		t.Fatalf("buildMultipartMessage error: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil || mediaType != "multipart/form-data" {
		t.Fatalf("unexpected Content-Type %q: %v", ct, err)
	}

	mr := multipart.NewReader(buf, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err != nil {
			break
		}
		cd := p.Header.Get("Content-Disposition")
		if strings.Contains(cd, "files[") {
			if strings.Contains(cd, "secret_report.pdf") && !strings.Contains(cd, "/home/") {
				// basename present, path stripped — pass
				return
			}
			if strings.Contains(cd, "/home/") {
				t.Errorf("filesystem path leaked into Content-Disposition: %q", cd)
				return
			}
			if !strings.Contains(cd, "secret_report.pdf") {
				t.Errorf("expected basename %q in Content-Disposition, got %q", "secret_report.pdf", cd)
			}
			return
		}
	}
	t.Error("no file part found in multipart body")
}
