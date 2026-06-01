package api

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"io"
	"mime"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/streame-gg/go-discord-wrapper/types/discord"
)

func (su *apiMessagesSuite) TestBasenameFilename() {
	t := su.T()
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

// TestBug123_MessageFileReader verifies that MessageFile.Reader is streamed via
// io.Copy when set, and that the resulting multipart body contains the correct
// file content (Bug 123).
func (su *apiMessagesSuite) TestBug123_MessageFileReader() {
	t := su.T()
	content := "streaming file content"
	payload := []byte(`{"content":"hello"}`)
	files := []discord.MessageFile{
		{Name: "stream.txt", ContentType: "text/plain", Reader: strings.NewReader(content)},
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
		if !strings.Contains(cd, "files[") {
			continue
		}
		body, readErr := io.ReadAll(p)
		if readErr != nil {
			t.Fatalf("reading file part: %v", readErr)
		}
		if !bytes.Equal(body, []byte(content)) {
			t.Errorf("file part body = %q, want %q", body, content)
		}
		return
	}
	t.Error("no file part found in multipart body (Bug 123)")
}

// TestBug123_MessageFileData verifies that the existing Data path still works
// when Reader is nil (regression guard for Bug 123 change).
func (su *apiMessagesSuite) TestBug123_MessageFileData() {
	t := su.T()
	content := []byte("raw bytes content")
	payload := []byte(`{}`)
	files := []discord.MessageFile{
		{Name: "data.bin", ContentType: "application/octet-stream", Data: content},
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
		if !strings.Contains(p.Header.Get("Content-Disposition"), "files[") {
			continue
		}
		body, readErr := io.ReadAll(p)
		if readErr != nil {
			t.Fatalf("reading file part: %v", readErr)
		}
		if !bytes.Equal(body, content) {
			t.Errorf("file part body = %q, want %q", body, content)
		}
		return
	}
	t.Error("no file part found in multipart body (Bug 123 regression)")
}

func (su *apiMessagesSuite) TestBuildMultipartMessage_BasenameInHeader() {
	t := su.T()
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

type apiMessagesSuite struct{ suite.Suite }

func TestApiMessagesSuite(t *testing.T) { suite.Run(t, new(apiMessagesSuite)) }
