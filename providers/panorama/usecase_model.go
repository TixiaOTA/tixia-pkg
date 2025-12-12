package panorama

import "errors"

// CancelBookingRequest represents the request for canceling bookings
type CancelBookingRequest struct {
	OrderNumbers []string `json:"order_numbers"`
}

func (r *CancelBookingRequest) Validate() error {
	if len(r.OrderNumbers) == 0 {
		return errors.New("order_numbers is required")
	}
	return nil
}

// CancelBookingResponse represents the response for canceling bookings
type CancelBookingResponse struct {
	Results []CancelBookingResult `json:"results"`
}

// CancelBookingResult represents the result for a single order number
type CancelBookingResult struct {
	OrderNumber string `json:"order_number"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	State       string `json:"state,omitempty"`
}
