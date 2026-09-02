package panorama

import "errors"

// --------------------------------- GET SSR AVAILABILITY REQUEST ---------------------------------

// GetSSRAvailabilityRequest is the request structure for GetSSRAvailability.
// The signature is the only parameter; it is sent in the Action/Signature headers
// and echoed in the request body (see the GetBalanceRequest pattern).
type GetSSRAvailabilityRequest struct {
	Signature string `json:"signature"`
}

func (r *GetSSRAvailabilityRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

// --------------------------------- GET SSR AVAILABILITY RESPONSE ---------------------------------

// GetSSRAvailabilityResponse is the response structure for GetSSRAvailability
type GetSSRAvailabilityResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SSRAvailabilityResponse `json:"response"`
	} `json:"result"`
}

// SSRAvailabilityResponse is the response payload for GetSSRAvailability
type SSRAvailabilityResponse struct {
	SSRAvailabilityProvider []SSRAvailabilityProvider `json:"ssr_availability_provider"`
}

// SSRAvailabilityProvider is a provider block in the GetSSRAvailability response
type SSRAvailabilityProvider struct {
	Provider        string          `json:"provider"`
	SSRAvailability SSRAvailability `json:"ssr_availability"`
	Status          string          `json:"status"`
	ScheduleID      int             `json:"schedule_id"`
	Sequence        int             `json:"sequence"`
}

// SSRAvailability is the SSR availability grouped by category (e.g. baggage, meal)
type SSRAvailability struct {
	Baggage []SSRAvailabilityItem `json:"baggage"`
	Meal    []SSRAvailabilityItem `json:"meal"`
}

// SSRAvailabilityItem is a journey-level SSR availability block for one category
type SSRAvailabilityItem struct {
	JourneyCode      string                 `json:"journey_code"`
	Segments         []SSRSegment           `json:"segments"`
	AvailabilityType string                 `json:"availability_type"`
	SSRs             []SSRAvailabilityEntry `json:"ssrs"`
}

// SSRSegment is a segment within an SSR availability item
type SSRSegment struct {
	CarrierCode   string `json:"carrier_code"`
	CarrierNumber string `json:"carrier_number"`
	Origin        string `json:"origin"`
	DepartureDate string `json:"departure_date"`
	Destination   string `json:"destination"`
	ArrivalDate   string `json:"arrival_date"`
	Provider      string `json:"provider"`
}

// SSRAvailabilityEntry is a single sellable SSR in the GetSSRAvailability response
type SSRAvailabilityEntry struct {
	SSRCode          string             `json:"ssr_code"`
	SSRType          string             `json:"ssr_type"`
	ImageURL         string             `json:"image_url"`
	AvailableCount   int                `json:"available_count"`
	Name             string             `json:"name"`
	Description      []any              `json:"description"`
	SSRRef           []any              `json:"ssr_ref"`
	ServiceCharges   []SSRServiceCharge `json:"service_charges"`
	BaseFare         float64            `json:"base_fare"`
	BaseTax          float64            `json:"base_tax"`
	BaseNta          float64            `json:"base_nta"`
	BaseCommission   float64            `json:"base_commission"`
	AvailabilityType string             `json:"availability_type"`
	JourneyCode      string             `json:"journey_code"`
}

// SSRServiceCharge is a service charge line for an SSR
type SSRServiceCharge struct {
	ChargeCode      string  `json:"charge_code"`
	ChargeType      string  `json:"charge_type"`
	Currency        string  `json:"currency"`
	PaxType         string  `json:"pax_type"`
	PaxCount        int     `json:"pax_count"`
	Amount          float64 `json:"amount"`
	ForeignCurrency string  `json:"foreign_currency"`
	ForeignAmount   float64 `json:"foreign_amount"`
	Total           float64 `json:"total"`
	Sequence        int     `json:"sequence"`
}
