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

// SearchBookingRequest represents the request for searching bookings
type SearchBookingRequest struct {
	JourneyList []struct {
		Origin        string `json:"origin"`
		Destination   string `json:"destination"`
		DepartureDate string `json:"departure_date"`
	} `json:"journey_list"`
	Adult        int      `json:"adult"`
	Child        int      `json:"child"`
	Infant       int      `json:"infant"`
	CabinClass   string   `json:"cabin_class"`
	CarrierCodes []string `json:"carrier_codes"`
	IsComboPrice bool     `json:"is_combo_price"`
	Provider     string   `json:"provider"`
	PromoCodes   []string `json:"promo_codes"`
}

func (r *SearchBookingRequest) Validate() error {
	if len(r.JourneyList) == 0 {
		return errors.New("journey_list is required")
	}
	for _, v := range r.JourneyList {
		if v.Origin == "" {
			return errors.New("origin is required")
		}
		if v.Destination == "" {
			return errors.New("destination is required")
		}
		if v.DepartureDate == "" {
			return errors.New("departure_date is required")
		}
	}
	if r.Adult <= 0 {
		return errors.New("adult is required")
	}
	if r.CabinClass == "" {
		return errors.New("cabin class is required")
	}
	if r.Provider == "" {
		return errors.New("provider is required")
	}
	return nil
}

// SearchBookingResponse represents the response for searching bookings
type SearchBookingResponse struct {
	Signature string       `json:"signature"`
	Results   SearchResult `json:"results"`
}
