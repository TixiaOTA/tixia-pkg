package train

import (
	"fmt"
	"strings"
)

// AssignSeatsRequest follows the Orbis train assign_seats contract.
// Request passenger_sequence is zero-based; provider response is one-based.
type AssignSeatsRequest struct {
	OrderNumber      string                      `json:"order_number"`
	ProviderBookings []AssignSeatProviderBooking `json:"provider_bookings"`
}
type AssignSeatProviderBooking struct {
	Provider string              `json:"provider"`
	Sequence int                 `json:"sequence"`
	Journeys []AssignSeatJourney `json:"journeys"`
}
type AssignSeatJourney struct {
	Sequence int          `json:"sequence"`
	Seats    []AssignSeat `json:"seats"`
}
type AssignSeat struct {
	CabinCode         string `json:"cabin_code"`
	SeatRow           int    `json:"seat_row"`
	SeatColumn        string `json:"seat_column"`
	PassengerSequence int    `json:"passenger_sequence"`
}

func (r AssignSeatsRequest) Validate() error {
	if strings.TrimSpace(r.OrderNumber) == "" {
		return fmt.Errorf("order_number is required")
	}
	if len(r.ProviderBookings) == 0 || len(r.ProviderBookings) > 2 {
		return fmt.Errorf("provider_bookings must contain one or two entries")
	}
	seenProviders := map[int]bool{}
	for i, b := range r.ProviderBookings {
		if b.Provider != ProviderKAI || b.Sequence < 1 || seenProviders[b.Sequence] {
			return fmt.Errorf("provider_bookings[%d] invalid provider or sequence", i)
		}
		seenProviders[b.Sequence] = true
		if len(b.Journeys) != 1 {
			return fmt.Errorf("provider_bookings[%d] must contain one journey", i)
		}
		seenJourney := map[int]bool{}
		for j, journey := range b.Journeys {
			if journey.Sequence < 1 || seenJourney[journey.Sequence] || len(journey.Seats) == 0 {
				return fmt.Errorf("provider_bookings[%d].journeys[%d] invalid sequence or seats", i, j)
			}
			seenJourney[journey.Sequence] = true
			passengers := map[int]bool{}
			seats := map[string]bool{}
			for k, s := range journey.Seats {
				if strings.TrimSpace(s.CabinCode) == "" || s.SeatRow < 1 || strings.TrimSpace(s.SeatColumn) == "" || s.PassengerSequence < 0 {
					return fmt.Errorf("invalid seat at provider_bookings[%d].journeys[%d].seats[%d]", i, j, k)
				}
				key := fmt.Sprintf("%s:%d:%s", strings.ToUpper(s.CabinCode), s.SeatRow, strings.ToUpper(s.SeatColumn))
				if passengers[s.PassengerSequence] || seats[key] {
					return fmt.Errorf("duplicate passenger or seat in provider_bookings[%d].journeys[%d]", i, j)
				}
				passengers[s.PassengerSequence] = true
				seats[key] = true
			}
		}
	}
	return nil
}

type AssignSeatsResponse struct {
	JSONRPC string            `json:"jsonrpc"`
	Result  AssignSeatsResult `json:"result"`
}
type AssignSeatsResult struct {
	BasicResult
	Response []AssignedProviderBooking `json:"response"`
}
type AssignedProviderBooking struct {
	Sequence int               `json:"sequence"`
	Journeys []AssignedJourney `json:"journeys"`
	Status   string            `json:"status"`
}
type AssignedJourney struct {
	Sequence int            `json:"sequence"`
	Seats    []AssignedSeat `json:"seats"`
}
type AssignedSeat struct {
	Seat              string         `json:"seat"`
	SeatCode          FlexibleString `json:"seat_code"`
	PassengerSequence int            `json:"passenger_sequence"`
}

// GetBooking always uses a fresh session; order number survives a prior signin session.
type GetBookingRequest struct {
	OrderNumber string `json:"order_number"`
}

func (r GetBookingRequest) Validate() error {
	if strings.TrimSpace(r.OrderNumber) == "" {
		return fmt.Errorf("order_number is required")
	}
	return nil
}

type GetBookingResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  GetBookingResult `json:"result"`
}
type GetBookingResult struct {
	BasicResult
	Response CommitBookingData `json:"response"`
}
