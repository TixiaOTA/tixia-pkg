package mg

import (
	"context"
	"encoding/json"
	"io"
)

func (c *Client) SearchHotel(ctx context.Context, req SearchHotelRequest) (*SearchHotelResponse, *MGError) {
	path := "/SearchHotel"

	httpResp, err := c.sendRequest(ctx, req, path, nil)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}
	defer httpResp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	// Try to unmarshal as MGError first
	var mgError MGError
	if err := json.Unmarshal(bodyBytes, &mgError); err == nil {
		// Check if it's a valid MGError
		// Status = false, then it's an error
		if !mgError.Status || mgError.ErrorCode != "" || mgError.ErrorMessage != "" {
			return nil, &MGError{
				Status:       false,
				ErrorCode:    mgError.ErrorCode,
				ErrorMessage: mgError.ErrorMessage,
			}
		}
	}

	// If not an MGError, unmarshal as GetHotelsResponse
	var resp SearchHotelResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, &MGError{
			Status:       false,
			ErrorCode:    InternalServerError,
			ErrorMessage: err.Error(),
		}
	}

	return &resp, nil
}
