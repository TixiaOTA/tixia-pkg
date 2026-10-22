package mailgun

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

func basicAuthHeader(apiKey string) string {
	credentials := base64.StdEncoding.EncodeToString([]byte("api:" + apiKey))
	return "Basic " + credentials
}

func buildForm(req SendEmailRequest) url.Values {
	values := url.Values{}

	values.Set("from", req.From)

	for _, recipient := range req.To {
		values.Add("to", recipient)
	}

	for _, recipient := range req.CC {
		values.Add("cc", recipient)
	}

	for _, recipient := range req.BCC {
		values.Add("bcc", recipient)
	}

	values.Set("subject", req.Subject)

	if req.Text != "" {
		values.Set("text", req.Text)
	}

	if req.HTML != "" {
		values.Set("html", req.HTML)
	}

	if req.ReplyTo != "" {
		values.Set("h:Reply-To", req.ReplyTo)
	}

	for _, tag := range req.Tags {
		values.Add("o:tag", tag)
	}

	for key, value := range req.Headers {
		values.Set("h:"+key, value)
	}

	for key, value := range req.Variables {
		encoded, err := json.Marshal(value)
		if err != nil {
			continue
		}

		values.Set("v:"+key, string(encoded))
	}

	if req.TestMode {
		values.Set("o:testmode", "yes")
	}

	return values
}

// buildMultipartBody encodes the request as multipart/form-data, which Mailgun
// requires when attachments or inline files are included.
func buildMultipartBody(req SendEmailRequest) ([]byte, string, error) {
	var buffer bytes.Buffer

	writer := multipart.NewWriter(&buffer)

	for key, values := range buildForm(req) {
		for _, value := range values {
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", err
			}
		}
	}

	for _, attachment := range req.Attachments {
		if err := writeAttachment(writer, "attachment", attachment); err != nil {
			return nil, "", err
		}
	}

	for _, inline := range req.Inline {
		if err := writeAttachment(writer, "inline", inline); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return buffer.Bytes(), writer.FormDataContentType(), nil
}

func writeAttachment(writer *multipart.Writer, field string, attachment Attachment) error {
	contentType := attachment.ContentType
	if contentType == "" {
		contentType = http.DetectContentType(attachment.Content)
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(
		`form-data; name="%s"; filename="%s"`,
		field,
		escapeQuotes(attachment.Filename),
	))
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = part.Write(attachment.Content)
	return err
}

func escapeQuotes(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
	return replacer.Replace(value)
}

func parseAPIError(statusCode int, body []byte) *APIError {
	apiError := &APIError{
		StatusCode: statusCode,
		RawBody:    body,
	}

	var payload struct {
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &payload); err == nil && payload.Message != "" {
		apiError.Message = payload.Message
		return apiError
	}

	if message := strings.TrimSpace(string(body)); message != "" {
		apiError.Message = message
		return apiError
	}

	apiError.Message = http.StatusText(statusCode)
	return apiError
}
