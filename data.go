package email

// Data is the data required to build and send an email.
type Data struct {
	To          []string     `json:"to"`
	CC          *[]string    `json:"cc"`
	BCC         *[]string    `json:"bcc"`
	ReplyTo     *[]string    `json:"replyTo"`
	From        string       `json:"from"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	ContentType ContentType  `json:"contentType"`
	Attachments []Attachment `json:"attachments"`
}
