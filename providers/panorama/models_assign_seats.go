package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- ASSIGN SEATS REQUEST ---------------------------------

// AssignSeatsRequest is the request structure for AssignSeats
type AssignSeatsRequest struct {
	Signature          string       `json:"signature"`
	SegmentSeatRequest []AssignSeat `json:"segment_seat_request"`
}

func (r *AssignSeatsRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if len(r.SegmentSeatRequest) == 0 {
		return errors.New("segment_seat_request is required")
	}
	for i, segment := range r.SegmentSeatRequest {
		if segment.SegmentCode == "" {
			return fmt.Errorf("segment_seat_request[%d].segment_code is required", i)
		}
		if segment.Provider == "" {
			return fmt.Errorf("segment_seat_request[%d].provider is required", i)
		}
		if len(segment.Passengers) == 0 {
			return fmt.Errorf("segment_seat_request[%d].passengers is required", i)
		}
		for j, passenger := range segment.Passengers {
			if passenger.SeatCode == "" {
				return fmt.Errorf("segment_seat_request[%d].passengers[%d].seat_code is required", i, j)
			}
		}
	}
	return nil
}

// AssignSeat is a seat assignment for a segment.
type AssignSeat struct {
	SegmentCode string                `json:"segment_code"`
	Provider    string                `json:"provider"`
	Passengers  []AssignSeatPassenger `json:"passengers"`
}

// AssignSeatPassenger is a passenger selection for an assigned seat
type AssignSeatPassenger struct {
	PassengerNumber int    `json:"passenger_number"`
	SeatCode        string `json:"seat_code"`
}

// --------------------------------- ASSIGN SEATS RESPONSE ---------------------------------

// AssignSeatsResponse is the response structure for AssignSeats
type AssignSeatsResponse struct {
	Result struct {
		BasicPanoramaResult
		Response AssignSeatsResult `json:"response"`
	} `json:"result"`
}

// AssignSeatsResult is the response payload for AssignSeats.
// The response carries booking/price-itinerary-like data and is provider
// specific, so only the fields callers need are modeled.
type AssignSeatsResult struct {
	TotalPrice         float64                    `json:"total_price"`
	TotalNta           float64                    `json:"total_nta"`
	OrderNumber        string                     `json:"order_number"`
	AssignSeatProvider []AssignSeatProviderResult `json:"assign_seat_provider"`
}

// AssignSeatProviderResult is a provider block in the AssignSeats response
type AssignSeatProviderResult struct {
	Provider   string                      `json:"provider"`
	Journeys   []AssignSeatJourney         `json:"journeys"`
	Passengers []AssignSeatPassengerResult `json:"passengers"`
}

// AssignSeatJourney is a journey within an AssignSeats provider response
type AssignSeatJourney struct {
	JourneyCode string              `json:"journey_code"`
	Segments    []AssignSeatSegment `json:"segments"`
}

// AssignSeatSegment is a segment within an AssignSeats journey
type AssignSeatSegment struct {
	SegmentCode   string `json:"segment_code"`
	CarrierCode   string `json:"carrier_code"`
	CarrierNumber string `json:"carrier_number"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	ArrivalDate   string `json:"arrival_date"`
}

// AssignSeatPassengerResult is a passenger within an AssignSeats provider response
type AssignSeatPassengerResult struct {
	PassengerNumber int    `json:"passenger_number"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	PaxType         string `json:"pax_type"`
	Seats           []any  `json:"seats"`
}
