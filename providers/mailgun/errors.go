package mailgun

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAPIKey     = errors.New("mailgun: api key is required")
	ErrMissingDomain     = errors.New("mailgun: domain is required")
	ErrNilLogger         = errors.New("mailgun: logger is required")
	ErrMissingFrom       = errors.New("mailgun: from address is required")
	ErrMissingRecipients = errors.New("mailgun: at least one recipient is required")
	ErrMissingSubject    = errors.New("mailgun: subject is required")
	ErrMissingBody       = errors.New("mailgun: text or html body is required")
	ErrInvalidAttachment = errors.New("mailgun: attachment filename and content are required")
)

// APIError represents a non-2xx response from the Mailgun API.
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	RawBody    []byte `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("mailgun: API error (status %d): %s", e.StatusCode, e.Message)
}
