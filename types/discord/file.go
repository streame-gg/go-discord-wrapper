package discord

// MessageFile is a binary file attachment included in a message or webhook.
type MessageFile struct {
	// Name is the filename Discord will display (e.g. "image.png").
	Name string
	// ContentType is the MIME type (e.g. "image/png"). Defaults to "application/octet-stream".
	ContentType string
	// Data is the raw file content.
	Data []byte
}
