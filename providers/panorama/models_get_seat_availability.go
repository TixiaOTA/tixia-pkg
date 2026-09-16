package panorama

import "errors"

// --------------------------------- GET SEAT AVAILABILITY REQUEST ---------------------------------

// GetSeatAvailabilityRequest is the request structure for GetSeatAvailability.
// The signature is the only parameter; it is sent in the Action/Signature headers
// and echoed in the request body (see the GetBalanceRequest pattern).
type GetSeatAvailabilityRequest struct {
	Signature string `json:"signature"`
}

func (r *GetSeatAvailabilityRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	return nil
}

// --------------------------------- GET SEAT AVAILABILITY RESPONSE ---------------------------------

// GetSeatAvailabilityResponse is the response structure for GetSeatAvailability
type GetSeatAvailabilityResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SeatAvailabilityResponse `json:"response"`
	} `json:"result"`
}

// SeatAvailabilityResponse is the response payload for GetSeatAvailability
type SeatAvailabilityResponse struct {
	SeatAvailabilityProvider []SeatAvailabilityProvider `json:"seat_availability_provider"`
}

// SeatAvailabilityProvider is a provider block in the GetSeatAvailability response
type SeatAvailabilityProvider struct {
	Provider   string                    `json:"provider"`
	Status     string                    `json:"status"`
	ScheduleID int                       `json:"schedule_id"`
	Sequence   int                       `json:"sequence"`
	Segments   []SeatAvailabilitySegment `json:"segments"`
}

// SeatAvailabilitySegment is a flight segment in the GetSeatAvailability response
type SeatAvailabilitySegment struct {
	Origin               string      `json:"origin"`
	OriginTerminal       string      `json:"origin_terminal"`
	OriginUtc            int         `json:"origin_utc"`
	DepartureDate        string      `json:"departure_date"`
	Destination          string      `json:"destination"`
	DestinationTerminal  string      `json:"destination_terminal"`
	DestinationUtc       int         `json:"destination_utc"`
	ArrivalDate          string      `json:"arrival_date"`
	SegmentCode          string      `json:"segment_code"`
	SegmentCode2         string      `json:"segment_code2"`
	RouteCode            string      `json:"route_code"`
	ElapsedTime          string      `json:"elapsed_time"`
	TransitCount         int         `json:"transit_count"`
	TransitDuration      string      `json:"transit_duration"`
	Provider             string      `json:"provider"`
	Sequence             int         `json:"sequence"`
	CarrierName          string      `json:"carrier_name"`
	CarrierCode          string      `json:"carrier_code"`
	CarrierTypeCode      string      `json:"carrier_type_code"`
	CarrierNumber        string      `json:"carrier_number"`
	OperatingAirlineCode string      `json:"operating_airline_code"`
	SeatCabins           []SeatCabin `json:"seat_cabins"`
}

// SeatCabin is a cabin (travel class) with its seat rows in the seat map
type SeatCabin struct {
	CabinName      string         `json:"cabin_name"`
	CabinCode      string         `json:"cabin_code"`
	CabinTypeCode  string         `json:"cabin_type_code"`
	FirstRowNumber FlexibleString `json:"first_row_number"`
	LastRowNumber  FlexibleString `json:"last_row_number"`
	SeatRows       []SeatRow      `json:"seat_rows"`
	Sequence       int            `json:"sequence"`
}

// SeatRow is a row of seats in the seat map
type SeatRow struct {
	RowNumber FlexibleString `json:"row_number"`
	Seats     []Seat         `json:"seats"`
	Sequence  int            `json:"sequence"`
}

// Seat is a single seat in the seat map
type Seat struct {
	SeatName       string              `json:"seat_name"`
	SeatCode       string              `json:"seat_code"`
	Column         string              `json:"column"`
	Availability   int                 `json:"availability"`
	SeatType       FlexibleString      `json:"seat_type"`
	SeatRef        string              `json:"seat_ref"`
	ServiceCharges []SeatServiceCharge `json:"service_charges"`
	BaseFare       float64             `json:"base_fare"`
	BaseTax        float64             `json:"base_tax"`
	BaseNta        float64             `json:"base_nta"`
	BaseCommission float64             `json:"base_commission"`
	Description    []any               `json:"description"`
	Sequence       int                 `json:"sequence"`
}

// SeatServiceCharge is a service charge line for a seat
type SeatServiceCharge struct {
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
