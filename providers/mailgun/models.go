package mailgun

import "fmt"

// Attachment represents a file attached to an email.
type Attachment struct {
	Filename    string // Filename shown to the recipient
	Content     []byte // File content
	ContentType string // MIME type (optional, detected from content when empty)
}

// SendEmailRequest represents an email to be sent through Mailgun.
type SendEmailRequest struct {
	From        string            // Sender address (required)
	To          []string          // Recipient addresses (required)
	CC          []string          // CC addresses
	BCC         []string          // BCC addresses
	Subject     string            // Email subject (required)
	Text        string            // Plain-text body
	HTML        string            // HTML body
	ReplyTo     string            // Reply-To address
	Attachments []Attachment      // Files attached to the email
	Inline      []Attachment      // Files referenced from the HTML body
	Tags        []string          // Mailgun tags for analytics
	Headers     map[string]string // Additional message headers
	Variables   map[string]any    // Template variables
	TestMode    bool              // When true, Mailgun accepts the message but does not deliver it
}

// Validate checks that the request has the fields required to send an email.
func (r SendEmailRequest) Validate() error {
	if r.From == "" {
		return ErrMissingFrom
	}

	if len(r.To) == 0 {
		return ErrMissingRecipients
	}

	if r.Subject == "" {
		return ErrMissingSubject
	}

	if r.Text == "" && r.HTML == "" {
		return ErrMissingBody
	}

	for _, attachments := range [][]Attachment{r.Attachments, r.Inline} {
		for _, attachment := range attachments {
			if attachment.Filename == "" || len(attachment.Content) == 0 {
				return fmt.Errorf("%w: %q", ErrInvalidAttachment, attachment.Filename)
			}
		}
	}

	return nil
}

// SendEmailResponse is the response returned by Mailgun after accepting a message.
type SendEmailResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
