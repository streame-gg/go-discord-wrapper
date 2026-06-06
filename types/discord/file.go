package discord

import "io"

// MessageFile is a binary file attachment included in a message or webhook.
// https://docs.discord.com/developers/reference#uploading-files
type MessageFile struct {
	// Name is the filename Discord will display (e.g. "image.png").
	Name string
	// ContentType is the MIME type (e.g. "image/png"). Defaults to "application/octet-stream".
	ContentType string
	// Data is the raw file content. Used when Reader is nil.
	Data []byte
	// Reader is an optional streaming source for large files.
	// When set, the file part is written via io.Copy instead of buffering all
	// of Data into memory. Data is ignored when Reader is non-nil.
	Reader io.Reader
}
