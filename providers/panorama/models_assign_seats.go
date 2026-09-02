package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- ASSIGN SEATS REQUEST ---------------------------------

// AssignSeatsRequest is the request structure for AssignSeats
type AssignSeatsRequest struct {
	Signature          string               `json:"signature"`
	AssignSeatsRequest []AssignSeatProvider `json:"assign_seats_request"`
}

func (r *AssignSeatsRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if len(r.AssignSeatsRequest) == 0 {
		return errors.New("assign_seats_request is required")
	}
	for i, provider := range r.AssignSeatsRequest {
		if provider.Provider == "" {
			return fmt.Errorf("assign_seats_request[%d].provider is required", i)
		}
		if len(provider.AssignSeats) == 0 {
			return fmt.Errorf("assign_seats_request[%d].assign_seats is required", i)
		}
		for j, assignSeat := range provider.AssignSeats {
			if assignSeat.SegmentCode == "" {
				return fmt.Errorf("assign_seats_request[%d].assign_seats[%d].segment_code is required", i, j)
			}
			if len(assignSeat.Passengers) == 0 {
				return fmt.Errorf("assign_seats_request[%d].assign_seats[%d].passengers is required", i, j)
			}
			for k, passenger := range assignSeat.Passengers {
				if passenger.SeatCode == "" {
					return fmt.Errorf("assign_seats_request[%d].assign_seats[%d].passengers[%d].seat_code is required", i, j, k)
				}
			}
		}
	}
	return nil
}

// AssignSeatProvider is a provider block in the AssignSeats request
type AssignSeatProvider struct {
	AssignSeats []AssignSeat `json:"assign_seats"`
	Provider    string       `json:"provider"`
}

// AssignSeat is a seat assignment for a segment
type AssignSeat struct {
	SegmentCode string                `json:"segment_code"`
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
