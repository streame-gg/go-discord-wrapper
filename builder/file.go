package builder

import "github.com/streame-gg/go-discord-wrapper/types/components"

// ── File display ──────────────────────────────────────────────────────────────

// FileDisplayBuilder builds a FileComponent (Components v2 file preview) using a fluent API.
//
//	f := builder.NewFileDisplay("https://cdn.example.com/report.pdf").
//	    SetSpoiler(true).
//	    Build()
//
// https://docs.discord.com/developers/components/reference#file
type FileDisplayBuilder struct {
	f components.FileComponent
}

func NewFileDisplay(url string) *FileDisplayBuilder {
	return &FileDisplayBuilder{f: components.FileComponent{File: &components.UnfurledMediaItem{URL: url}}}
}

func (b *FileDisplayBuilder) SetName(name string) *FileDisplayBuilder {
	b.f.Name = name
	return b
}

func (b *FileDisplayBuilder) SetSpoiler(spoiler bool) *FileDisplayBuilder {
	b.f.Spoiler = spoiler
	return b
}

func (b *FileDisplayBuilder) Build() *components.FileComponent {
	return &b.f
}

// ── File upload ───────────────────────────────────────────────────────────────

// FileUploadBuilder builds a FileUploadComponent using a fluent API.
//
//	fu := builder.NewFileUpload("attachment").
//	    SetRequired(true).
//	    SetMaxValues(3).
//	    Build()
//
// https://docs.discord.com/developers/components/reference#file-upload
type FileUploadBuilder struct {
	fu components.FileUploadComponent
}

func NewFileUpload(customID string) *FileUploadBuilder {
	return &FileUploadBuilder{fu: components.FileUploadComponent{CustomID: customID}}
}

func (b *FileUploadBuilder) SetCustomID(id string) *FileUploadBuilder {
	b.fu.CustomID = id
	return b
}

func (b *FileUploadBuilder) SetRequired(required bool) *FileUploadBuilder {
	b.fu.Required = required
	return b
}

func (b *FileUploadBuilder) SetMinValues(min int) *FileUploadBuilder {
	b.fu.MinValues = &min
	return b
}

func (b *FileUploadBuilder) SetMaxValues(max int) *FileUploadBuilder {
	b.fu.MaxValues = &max
	return b
}

func (b *FileUploadBuilder) Build() *components.FileUploadComponent {
	return &b.fu
}
