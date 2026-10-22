# Mailgun Package

A Go client library for sending email through the Mailgun API. This package supports plain-text and HTML messages, multiple recipients, attachments, inline files, tags, custom headers, and template variables.

## Overview

The Mailgun package enables you to:
- Send email through the Mailgun v3 API
- Send plain-text, HTML, or multipart (text + HTML) messages
- Address multiple To, CC, and BCC recipients
- Attach files and inline images
- Add tags, custom headers, and template variables
- Use the US or EU Mailgun API endpoint
- Retry transient failures with exponential backoff
- Trace requests with OpenTelemetry and structured logging

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/providers/mailgun
```

## Configuration

### Client Configuration

Create a Mailgun client with your API key and sending domain:

```go
import (
    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/mailgun"
)

config := mailgun.Config{
    APIKey: "your-private-api-key",
    Domain: "mg.example.com",
}

mailgunClient, err := mailgun.New(config, logger.Default())
if err != nil {
    log.Fatal(err)
}
```

For the EU region, set the base URL:

```go
config := mailgun.Config{
    APIKey:  "your-private-api-key",
    Domain:  "mg.example.com",
    BaseURL: "https://api.eu.mailgun.net",
}
```

### Configuration Fields

- **APIKey** (string, required): Mailgun private API key. Used as the Basic Auth password with the username `api`.
- **Domain** (string, required): Mailgun sending domain (for example, `mg.example.com`). Do not include a protocol or path.
- **BaseURL** (string, optional): Mailgun API base URL. Defaults to `https://api.mailgun.net`. Use `https://api.eu.mailgun.net` for the EU region.
- **Timeout** (*time.Duration, optional): Timeout for each request attempt. Defaults to `30s` in the HTTP client.
- **RetryConfig** (*client.RetryConfig, optional): Retry behavior. Defaults to 2 retries with a 1 second initial backoff, 8 second maximum backoff, and a 2x multiplier.

## Usage

### Send a Plain-Text Email

```go
import (
    "context"

    "github.com/TixiaOTA/tixia-pkg/providers/mailgun"
)

response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:    "Tixia <no-reply@mg.example.com>",
    To:      []string{"customer@example.com"},
    Subject: "Your booking is confirmed",
    Text:    "Thank you for your booking.",
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(response.ID)
fmt.Println(response.Message)
```

### Send an HTML Email

```go
response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:    "Tixia <no-reply@mg.example.com>",
    To:      []string{"customer@example.com"},
    CC:      []string{"operations@example.com"},
    Subject: "Your booking is confirmed",
    Text:    "Thank you for your booking.",
    HTML:    "<h1>Booking confirmed</h1><p>Thank you for your booking.</p>",
    ReplyTo: "support@example.com",
})
```

### Send an Email with Attachments

```go
file, err := os.ReadFile("invoice.pdf")
if err != nil {
    log.Fatal(err)
}

response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:    "Tixia <no-reply@mg.example.com>",
    To:      []string{"customer@example.com"},
    Subject: "Your invoice",
    Text:    "Your invoice is attached.",
    Attachments: []mailgun.Attachment{
        {
            Filename:    "invoice.pdf",
            Content:     file,
            ContentType: "application/pdf", // Optional; detected when omitted
        },
    },
})
```

### Send an Email with an Inline Image

Inline files are referenced from the HTML body with a `cid:` URL:

```go
logo, err := os.ReadFile("logo.png")
if err != nil {
    log.Fatal(err)
}

response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:    "Tixia <no-reply@mg.example.com>",
    To:      []string{"customer@example.com"},
    Subject: "Welcome",
    HTML:    `<p>Welcome to Tixia</p><img src="cid:logo.png">`,
    Inline: []mailgun.Attachment{
        {Filename: "logo.png", Content: logo},
    },
})
```

### Add Tags, Headers, and Variables

```go
response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:     "Tixia <no-reply@mg.example.com>",
    To:       []string{"customer@example.com"},
    Subject:  "Your booking is confirmed",
    Text:     "Hello {{name}}",
    Tags:     []string{"booking-confirmation"},
    Headers:  map[string]string{"X-Booking-ID": "BK-123"},
    Variables: map[string]any{
        "name": "Ada",
    },
})
```

Variables are JSON-encoded and sent as `v:<key>` fields, which Mailgun exposes to templates.

### Test Mode

Set `TestMode` to have Mailgun accept the message without delivering it:

```go
response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
    From:     "Tixia <no-reply@mg.example.com>",
    To:       []string{"customer@example.com"},
    Subject:  "Test message",
    Text:     "This message will not be delivered.",
    TestMode: true,
})
```

## API Reference

### Public Functions

#### New

Creates a new Mailgun client.

```go
func New(config Config, logger *logger.Logger) (*Client, error)
```

**Parameters:**

- `config` (Config, required): Mailgun client configuration
- `logger` (*logger.Logger, required): Logger instance

**Return Values:**

- `*Client`: Mailgun client instance
- `error`: `ErrMissingAPIKey`, `ErrMissingDomain`, or `ErrNilLogger` when configuration is invalid

**Special Notes:** Applies default base URL and retry configuration when they are not provided. Trailing slashes are removed from `BaseURL`.

### Client Methods

#### SendEmail

Sends an email through the Mailgun `POST /v3/{domain}/messages` endpoint.

```go
func (c *Client) SendEmail(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error)
```

**Request Parameters:**

- `From` (string, required): Sender address. Display names are allowed.
- `To` ([]string, required): Recipient addresses. At least one is required.
- `CC` ([]string, optional): CC recipients.
- `BCC` ([]string, optional): BCC recipients.
- `Subject` (string, required): Email subject.
- `Text` (string, optional): Plain-text body.
- `HTML` (string, optional): HTML body. At least one of `Text` or `HTML` is required.
- `ReplyTo` (string, optional): Reply-To address.
- `Attachments` ([]Attachment, optional): Files attached to the message.
- `Inline` ([]Attachment, optional): Files referenced from the HTML body.
- `Tags` ([]string, optional): Mailgun tags.
- `Headers` (map[string]string, optional): Additional message headers.
- `Variables` (map[string]any, optional): Template variables.
- `TestMode` (bool, optional): When true, Mailgun accepts the message but does not deliver it.

**Response:**

- `ID` (string): Mailgun message ID.
- `Message` (string): Mailgun status message, for example `Queued. Thank you.`

**Timeout:** 30 seconds per attempt by default.

**Retries:** 2 retries by default on network errors and HTTP 5xx responses.

**Special Notes:** Requests with attachments or inline files are sent as `multipart/form-data`; other requests use `application/x-www-form-urlencoded`. Retrying a send can deliver a duplicate message if Mailgun accepted the original request but the response was lost.

## Error Handling

Configuration and validation errors are sentinel values:

- **ErrMissingAPIKey**: `APIKey` was not provided
- **ErrMissingDomain**: `Domain` was not provided
- **ErrNilLogger**: A nil logger was passed to `New`
- **ErrMissingFrom**: `From` was not provided
- **ErrMissingRecipients**: `To` is empty
- **ErrMissingSubject**: `Subject` was not provided
- **ErrMissingBody**: Neither `Text` nor `HTML` was provided
- **ErrInvalidAttachment**: An attachment is missing a filename or content

Non-2xx responses return `*APIError`:

```go
type APIError struct {
    StatusCode int
    Message    string
    RawBody    []byte
}
```

**Example:**

```go
_, err := mailgunClient.SendEmail(ctx, request)
if err != nil {
    var apiError *mailgun.APIError
    if errors.As(err, &apiError) {
        log.Printf("Mailgun rejected the message (%d): %s", apiError.StatusCode, apiError.Message)
    } else if errors.Is(err, mailgun.ErrMissingRecipients) {
        log.Println("No recipients provided")
    } else {
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Data Models

### SendEmailRequest

```go
type SendEmailRequest struct {
    From        string
    To          []string
    CC          []string
    BCC         []string
    Subject     string
    Text        string
    HTML        string
    ReplyTo     string
    Attachments []Attachment
    Inline      []Attachment
    Tags        []string
    Headers     map[string]string
    Variables   map[string]any
    TestMode    bool
}
```

### Attachment

```go
type Attachment struct {
    Filename    string
    Content     []byte
    ContentType string
}
```

### SendEmailResponse

```go
type SendEmailResponse struct {
    ID      string
    Message string
}
```

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/mailgun"
)

func main() {
    mailgunClient, err := mailgun.New(mailgun.Config{
        APIKey: os.Getenv("MAILGUN_API_KEY"),
        Domain: os.Getenv("MAILGUN_DOMAIN"),
    }, logger.Default())
    if err != nil {
        log.Fatal(err)
    }

    invoice, err := os.ReadFile("invoice.pdf")
    if err != nil {
        log.Fatal(err)
    }

    response, err := mailgunClient.SendEmail(context.Background(), mailgun.SendEmailRequest{
        From:    "Tixia <no-reply@mg.example.com>",
        To:      []string{"customer@example.com"},
        Subject: "Your booking is confirmed",
        Text:    "Thank you for your booking. Your invoice is attached.",
        HTML:    "<h1>Booking confirmed</h1><p>Your invoice is attached.</p>",
        Tags:    []string{"booking-confirmation"},
        Attachments: []mailgun.Attachment{
            {Filename: "invoice.pdf", Content: invoice, ContentType: "application/pdf"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Message queued: %s\n", response.ID)
}
```

## License

This package is part of the TixiaOTA project.
