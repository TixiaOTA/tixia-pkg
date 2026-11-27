package train

import (
	"fmt"
	"strings"
)

// IssueBookingRequest issues an existing held Orbis train booking.
type IssueBookingRequest struct {
	OrderNumber string `json:"order_number"`
}

func (r IssueBookingRequest) Validate() error {
	if strings.TrimSpace(r.OrderNumber) == "" {
		return fmt.Errorf("order_number is required")
	}
	return nil
}

// IssueBookingResponse retains the issued booking, including provider tickets.
type IssueBookingResponse struct {
	JSONRPC string             `json:"jsonrpc"`
	Result  IssueBookingResult `json:"result"`
}

type IssueBookingResult struct {
	BasicResult
	Response CommitBookingData `json:"response"`
}

// CancelBookingRequest cancels an unissued Orbis train booking.
type CancelBookingRequest struct {
	OrderNumber string `json:"order_number"`
}

func (r CancelBookingRequest) Validate() error {
	if strings.TrimSpace(r.OrderNumber) == "" {
		return fmt.Errorf("order_number is required")
	}
	return nil
}

// CancelBookingResponse contains one provider result per cancelled booking.
// A successful outer result does not necessarily mean all provider bookings succeeded.
type CancelBookingResponse struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  CancelBookingResult `json:"result"`
}

type CancelBookingResult struct {
	BasicResult
	Response []CancelBookingProviderResult `json:"response"`
}

type CancelBookingProviderResult struct {
	BasicResult
	Response string `json:"response"`
	Status   string `json:"status"`
}
