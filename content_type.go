package email

// ContentType is a string type alias for the different supported content types.
type ContentType string

const (
	// ContentTypeTextPlain is used for text/plain emails.
	ContentTypeTextPlain ContentType = "text/plain"
	// ContentTypeTextHTML is used for text/html emails.
	ContentTypeTextHTML ContentType = "text/html"
)
