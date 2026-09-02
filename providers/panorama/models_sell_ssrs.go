package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- SELL SSRS REQUEST ---------------------------------

// SellSSRsRequest is the request structure for SellSSRs
type SellSSRsRequest struct {
	Signature       string            `json:"signature"`
	SellSSRsRequest []SellSSRProvider `json:"sell_ssrs_request"`
}

func (r *SellSSRsRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if len(r.SellSSRsRequest) == 0 {
		return errors.New("sell_ssrs_request is required")
	}
	for i, provider := range r.SellSSRsRequest {
		if provider.Provider == "" {
			return fmt.Errorf("sell_ssrs_request[%d].provider is required", i)
		}
		if len(provider.SellSSRs) == 0 {
			return fmt.Errorf("sell_ssrs_request[%d].sell_ssrs is required", i)
		}
		for j, sellSSR := range provider.SellSSRs {
			if sellSSR.JourneyCode == "" {
				return fmt.Errorf("sell_ssrs_request[%d].sell_ssrs[%d].journey_code is required", i, j)
			}
			if sellSSR.AvailabilityType == "" {
				return fmt.Errorf("sell_ssrs_request[%d].sell_ssrs[%d].availability_type is required", i, j)
			}
			if len(sellSSR.Passengers) == 0 {
				return fmt.Errorf("sell_ssrs_request[%d].sell_ssrs[%d].passengers is required", i, j)
			}
			for k, passenger := range sellSSR.Passengers {
				if passenger.SSRCode == "" {
					return fmt.Errorf("sell_ssrs_request[%d].sell_ssrs[%d].passengers[%d].ssr_code is required", i, j, k)
				}
			}
		}
	}
	return nil
}

// SellSSRProvider is a provider block in the SellSSRs request
type SellSSRProvider struct {
	SellSSRs []SellSSR `json:"sell_ssrs"`
	Provider string    `json:"provider"`
}

// SellSSR is a single SSR to sell for a journey
type SellSSR struct {
	JourneyCode      string             `json:"journey_code"`
	AvailabilityType string             `json:"availability_type"`
	Passengers       []SellSSRPassenger `json:"passengers"`
}

// SellSSRPassenger is a passenger selection for a sold SSR
type SellSSRPassenger struct {
	PassengerNumber int    `json:"passenger_number"`
	SSRCode         string `json:"ssr_code"`
}

// --------------------------------- SELL SSRS RESPONSE ---------------------------------

// SellSSRsResponse is the response structure for SellSSRs
type SellSSRsResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SellSSRsResult `json:"response"`
	} `json:"result"`
}

// SellSSRsResult is the response payload for SellSSRs
type SellSSRsResult struct {
	TotalPrice      float64                 `json:"total_price"`
	TotalNta        float64                 `json:"total_nta"`
	SellSSRProvider []SellSSRProviderResult `json:"sell_ssr_provider"`
}

// SellSSRProviderResult is a provider block in the SellSSRs response.
// The response is large and provider-specific, so the journeys/passengers are
// modeled permissively with only the fields callers need.
type SellSSRProviderResult struct {
	Provider   string                   `json:"provider"`
	Journeys   []SellSSRJourney         `json:"journeys"`
	Passengers []SellSSRPassengerResult `json:"passengers"`
}

// SellSSRJourney is a journey within a SellSSRs provider response
type SellSSRJourney struct {
	JourneyCode string           `json:"journey_code"`
	Segments    []SellSSRSegment `json:"segments"`
}

// SellSSRSegment is a segment within a SellSSRs journey
type SellSSRSegment struct {
	SegmentCode   string `json:"segment_code"`
	FareCode      string `json:"fare_code"`
	CarrierCode   string `json:"carrier_code"`
	CarrierNumber string `json:"carrier_number"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	ArrivalDate   string `json:"arrival_date"`
}

// SellSSRPassengerResult is a passenger within a SellSSRs provider response
type SellSSRPassengerResult struct {
	PassengerNumber int    `json:"passenger_number"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PaxType         string `json:"pax_type"`
}
