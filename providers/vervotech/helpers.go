package vervotech

import (
	"net/url"
	"time"
)

// buildQueryString builds a query string from parameters
func buildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Add(key, value)
		}
	}

	query := values.Encode()
	if query != "" {
		return "?" + query
	}
	return ""
}

// buildHeaders builds common headers for Vervotech API requests
func (c *Client) buildHeaders(includeSyncID bool) map[string]string {
	headers := map[string]string{
		"accountId": c.config.AccountID,
		"apikey":    c.config.APIKey,
	}

	if includeSyncID && c.config.SyncID != "" {
		headers["syncId"] = c.config.SyncID
	}

	return headers
}

// formatDateTime formats time.Time to RFC3339 format as expected by the API
func formatDateTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
