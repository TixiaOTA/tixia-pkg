package panorama

import "errors"

// --------------------------------- SHARED REFUND TYPES ---------------------------------

// RefundPassengerNumber identifies a passenger on the order by their passenger number.
type RefundPassengerNumber struct {
	PassengerNumber int `json:"passenger_number"`
}

// RefundProviderBooking identifies one provider booking (PNR) to refund,
// with the journeys of that booking that take part in the refund.
type RefundProviderBooking struct {
	Pnr      string          `json:"pnr"`
	Journeys []RefundJourney `json:"journeys"`
}

// RefundJourney is one journey inside a RefundProviderBooking
type RefundJourney struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"` // Format: YYYY-MM-DD
	Pax           []int  `json:"pax"`
}

// --------------------------------- PRE REFUND LOGIN LIST REQUEST ---------------------------------

// PreRefundLoginListRequest is the request structure for PreRefundLoginList.
// It asks Panorama for the refund captcha challenge(s) of the given provider
// bookings before requesting the refund quote.
type PreRefundLoginListRequest struct {
	Signature   string   `json:"signature"`
	OrderNumber string   `json:"order_number"`
	Provider    []string `json:"provider,omitempty"`
	Pnr         []string `json:"pnr,omitempty"`
}

func (r *PreRefundLoginListRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	return nil
}

// --------------------------------- PRE REFUND LOGIN LIST RESPONSE ---------------------------------

// PreRefundLoginListResponse is the response structure for PreRefundLoginList
type PreRefundLoginListResponse struct {
	Result struct {
		BasicPanoramaResult
		Response PreRefundLoginListResponseData `json:"response"`
	} `json:"result"`
}

// PreRefundLoginListResponseData is the response payload for PreRefundLoginList
type PreRefundLoginListResponseData struct {
	PreRefundLogin []PreRefundLogin `json:"pre_refund_login"`
}

// PreRefundLogin is one refund captcha challenge, one per provider booking PNR.
// Img carries the captcha image (base64) to answer in GetRefundBookingRequest.Captcha.
type PreRefundLogin struct {
	Pnr string `json:"pnr"`
	Img string `json:"img"`
}

// --------------------------------- GET REFUND BOOKING REQUEST ---------------------------------

// GetRefundBookingRequest is the request structure for GetRefundBooking, the refund quote.
// Captcha must be provided by the caller, keyed by provider booking PNR with the
// answer to the challenge returned by PreRefundLoginList; the client does not fill it in.
type GetRefundBookingRequest struct {
	Signature        string                  `json:"signature"`
	OrderNumber      string                  `json:"order_number"`
	Passengers       []RefundPassengerNumber `json:"passengers"`
	ProviderBookings []RefundProviderBooking `json:"provider_bookings"`
	Captcha          map[string]string       `json:"captcha"`
}

func (r *GetRefundBookingRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	if len(r.Passengers) == 0 {
		return errors.New("passengers is required")
	}
	if len(r.ProviderBookings) == 0 {
		return errors.New("provider_bookings is required")
	}
	if len(r.Captcha) == 0 {
		return errors.New("captcha is required")
	}
	return nil
}

// --------------------------------- GET REFUND BOOKING RESPONSE ---------------------------------

// GetRefundBookingResponse is the response structure for GetRefundBooking
type GetRefundBookingResponse struct {
	Result struct {
		BasicPanoramaResult
		Response RefundBookingResponse `json:"response"`
	} `json:"result"`
}

// RefundBookingResponse is the refund quote payload for GetRefundBooking
type RefundBookingResponse struct {
	OrderNumber      string                       `json:"order_number"`
	Pnr              string                       `json:"pnr"`
	State            string                       `json:"state"`
	DepartureDate    string                       `json:"departure_date"`
	ArrivalDate      string                       `json:"arrival_date"`
	ProviderType     string                       `json:"provider_type"`
	ProviderBookings []RefundProviderBookingQuote `json:"provider_bookings"`
}

// RefundProviderBookingQuote is the refund quote returned for one provider booking.
// Amounts are FlexibleNumber because Panorama may return them as numbers or strings.
type RefundProviderBookingQuote struct {
	Pnr              string         `json:"pnr"`
	Provider         string         `json:"provider"`
	Status           string         `json:"status"`
	TotalPriceCost   FlexibleNumber `json:"total_price_cost"`
	PenaltyAmount    FlexibleNumber `json:"penalty_amount"`
	RefundableAmount FlexibleNumber `json:"refundable_amount"`
	AdminFee         FlexibleNumber `json:"admin_fee"`
	ResvTotalPrice   FlexibleNumber `json:"resv_total_price"`
}

// --------------------------------- CANCEL REFUND BOOKING REQUEST ---------------------------------

// CancelRefundBookingRequest is the request structure for CancelRefundBooking,
// which executes the refund for an issued order.
type CancelRefundBookingRequest struct {
	Signature        string                  `json:"signature"`
	OrderNumber      string                  `json:"order_number"`
	ProviderBookings []RefundProviderBooking `json:"provider_bookings"`
}

func (r *CancelRefundBookingRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	if len(r.ProviderBookings) == 0 {
		return errors.New("provider_bookings is required")
	}
	return nil
}

// --------------------------------- CANCEL REFUND BOOKING RESPONSE ---------------------------------

// CancelRefundBookingResponse is the response structure for CancelRefundBooking
type CancelRefundBookingResponse struct {
	Result struct {
		BasicPanoramaResult
		Response CancelRefundBookingResponseData `json:"response"`
	} `json:"result"`
}

// CancelRefundBookingResponseData is the response payload for CancelRefundBooking
type CancelRefundBookingResponseData struct {
	OrderNumber string `json:"order_number"`
}
