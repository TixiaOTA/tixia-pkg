package panorama

import (
	"errors"
	"fmt"
)

// CabinClass identifies the cabin class requested from Panorama.
type CabinClass string

const (
	CabinClassEconomy        CabinClass = "Y"
	CabinClassPremiumEconomy CabinClass = "W"
	CabinClassBusiness       CabinClass = "C"
	CabinClassFirst          CabinClass = "F"
)

// SearchRequest is the request structure for Search.
type SearchRequest struct {
	Signature    string          `json:"signature,omitempty" log:"redact"`
	JourneyList  []SearchJourney `json:"journey_list"`
	Adult        int             `json:"adult"`
	Child        int             `json:"child"`
	Infant       int             `json:"infant"`
	CabinClass   CabinClass      `json:"cabin_class"`
	CarrierCodes []string        `json:"carrier_codes"`
	IsComboPrice bool            `json:"is_combo_price"`
	Provider     string          `json:"provider"`
	PromoCodes   []string        `json:"promo_codes"`
}

// SearchJourney is a route in a Search request.
type SearchJourney struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	Pax           *[]int `json:"pax,omitempty"`
}

func (r *SearchRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.Adult <= 0 {
		return errors.New("adult is required")
	}
	if r.Provider == "" {
		return errors.New("provider is required")
	}
	if r.CabinClass == "" {
		return errors.New("cabin_class is required")
	}
	switch r.CabinClass {
	case CabinClassEconomy, CabinClassPremiumEconomy, CabinClassBusiness, CabinClassFirst:
	default:
		return errors.New("cabin_class must be one of Y, W, C, F")
	}
	if len(r.JourneyList) == 0 {
		return errors.New("journey_list is required")
	}
	for i, journey := range r.JourneyList {
		if journey.Origin == "" {
			return fmt.Errorf("journey_list[%d].origin is required", i)
		}
		if journey.Destination == "" {
			return fmt.Errorf("journey_list[%d].destination is required", i)
		}
		if journey.DepartureDate == "" {
			return fmt.Errorf("journey_list[%d].departure_date is required", i)
		}
	}
	return nil
}

// SearchResponse is the response structure for Search.
type SearchResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  struct {
		BasicPanoramaResult
		Response SearchResponseData `json:"response"`
	} `json:"result"`
}

// SearchResponseData is the response payload for Search.
type SearchResponseData struct {
	Schedules       []SearchSchedule       `json:"schedules"`
	Recommendations []SearchRecommendation `json:"recommendations"`
	PromoCodes      []any                  `json:"promo_codes"`
	JourneyKey      string                 `json:"journey_key"`
	Provider        string                 `json:"provider"`
}

// SearchSchedule is a schedule returned by Search.
type SearchSchedule struct {
	Journeys      []SearchJourneyResult `json:"journeys"`
	Provider      string                `json:"provider"`
	Origin        string                `json:"origin"`
	Destination   string                `json:"destination"`
	DepartureDate string                `json:"departure_date"`
	ErrorCode     int                   `json:"error_code"`
	ErrorMsg      string                `json:"error_msg"`
	Sequence      int                   `json:"sequence"`
}

// SearchJourneyResult is a journey within a Search schedule.
type SearchJourneyResult struct {
	Origin                   string          `json:"origin"`
	OriginTerminal           string          `json:"origin_terminal"`
	OriginUtc                int             `json:"origin_utc"`
	OriginCountry            string          `json:"origin_country"`
	DepartureDate            string          `json:"departure_date"`
	Destination              string          `json:"destination"`
	DestinationTerminal      string          `json:"destination_terminal"`
	DestinationUtc           int             `json:"destination_utc"`
	DestinationCountry       string          `json:"destination_country"`
	ArrivalDate              string          `json:"arrival_date"`
	JourneyCode              string          `json:"journey_code"`
	RouteCode                string          `json:"route_code"`
	ElapsedTime              string          `json:"elapsed_time"`
	TransitCount             int             `json:"transit_count"`
	Provider                 string          `json:"provider"`
	CarrierCodeList          []string        `json:"carrier_code_list"`
	OperatingAirlineCodeList []string        `json:"operating_airline_code_list"`
	Segments                 []SearchSegment `json:"segments"`
	JourneyRefID             string          `json:"journey_ref_id"`
	Sequence                 int             `json:"sequence"`
	SearchBanner             []any           `json:"search_banner"`
	ComplianceText           string          `json:"compliance_text"`
	TransitDuration          string          `json:"transit_duration"`
}

// SearchSegment is a segment within a Search journey.
type SearchSegment struct {
	Origin               string       `json:"origin"`
	OriginTerminal       string       `json:"origin_terminal"`
	OriginUtc            int          `json:"origin_utc"`
	OriginCountry        string       `json:"origin_country"`
	DepartureDate        string       `json:"departure_date"`
	Destination          string       `json:"destination"`
	DestinationTerminal  string       `json:"destination_terminal"`
	DestinationUtc       int          `json:"destination_utc"`
	DestinationCountry   string       `json:"destination_country"`
	ArrivalDate          string       `json:"arrival_date"`
	ElapsedTime          string       `json:"elapsed_time"`
	TransitCount         int          `json:"transit_count"`
	SegmentCode          string       `json:"segment_code"`
	SegmentCode2         string       `json:"segment_code2"`
	RouteCode            string       `json:"route_code"`
	TransitDuration      string       `json:"transit_duration"`
	Provider             string       `json:"provider"`
	Sequence             int          `json:"sequence"`
	CarrierName          string       `json:"carrier_name"`
	CarrierCode          string       `json:"carrier_code"`
	CarrierTypeCode      string       `json:"carrier_type_code"`
	CarrierTypeName      string       `json:"carrier_type_name"`
	CarrierNumber        string       `json:"carrier_number"`
	OperatingAirlineCode string       `json:"operating_airline_code"`
	Legs                 []SearchLeg  `json:"legs"`
	Fares                []SearchFare `json:"fares"`
	EquipmentType        string       `json:"equipment_type"`
	Status               string       `json:"status"`
	SeatCabins           []any        `json:"seat_cabins"`
}

// SearchLeg is a leg within a Search segment.
type SearchLeg struct {
	Origin               string         `json:"origin"`
	OriginTerminal       string         `json:"origin_terminal"`
	OriginUtc            int            `json:"origin_utc"`
	OriginCountry        string         `json:"origin_country"`
	DepartureDate        string         `json:"departure_date"`
	Destination          string         `json:"destination"`
	DestinationTerminal  string         `json:"destination_terminal"`
	DestinationUtc       int            `json:"destination_utc"`
	DestinationCountry   string         `json:"destination_country"`
	ArrivalDate          string         `json:"arrival_date"`
	LegCode              FlexibleString `json:"leg_code"`
	LegCode2             FlexibleString `json:"leg_code2"`
	RouteCode            string         `json:"route_code"`
	ElapsedTime          string         `json:"elapsed_time"`
	TransitDuration      string         `json:"transit_duration"`
	Provider             string         `json:"provider"`
	Sequence             int            `json:"sequence"`
	CarrierName          string         `json:"carrier_name"`
	CarrierCode          string         `json:"carrier_code"`
	CarrierTypeCode      string         `json:"carrier_type_code"`
	CarrierTypeName      string         `json:"carrier_type_name"`
	CarrierNumber        string         `json:"carrier_number"`
	OperatingAirlineCode string         `json:"operating_airline_code"`
	EquipmentType        string         `json:"equipment_type"`
}

// SearchFare is a fare within a Search segment.
type SearchFare struct {
	CabinClass           string                       `json:"cabin_class"`
	ClassOfService       string                       `json:"class_of_service"`
	TourCode             string                       `json:"tour_code"`
	FareCode             string                       `json:"fare_code"`
	FareBasisCode        string                       `json:"fare_basis_code"`
	FareClass            string                       `json:"fare_class"`
	FareName             string                       `json:"fare_name"`
	FareRules            []any                        `json:"fare_rules"`
	RuleCode             string                       `json:"rule_code"`
	Description          []any                        `json:"description"`
	AvailableCount       FlexibleNumber               `json:"available_count"`
	ServiceCharges       []SearchServiceCharge        `json:"service_charges"`
	FareDetails          []SearchFareDetail           `json:"fare_details"`
	GroupFareID          string                       `json:"group_fare_id"`
	FareRefID            string                       `json:"fare_ref_id"`
	PackageList          []any                        `json:"package_list"`
	EquipmentType        string                       `json:"equipment_type"`
	ValidatingCarrier    string                       `json:"validating_carrier"`
	IsTravelFair         bool                         `json:"is_travel_fair"`
	TravelFairText       string                       `json:"travel_fair_text"`
	FareRemarks          []any                        `json:"fare_remarks"`
	ServiceChargeSummary []SearchServiceChargeSummary `json:"service_charge_summary"`
	Sequence             int                          `json:"sequence"`
}

// SearchFareDetail is a fare detail within a Search fare.
type SearchFareDetail struct {
	DetailCode  string         `json:"detail_code"`
	DetailType  string         `json:"detail_type"`
	Amount      FlexibleString `json:"amount"`
	Unit        string         `json:"unit"`
	DetailName  string         `json:"detail_name"`
	Description []any          `json:"description"`
	Sequence    int            `json:"sequence"`
}

// SearchServiceChargeSummary is the passenger-level fare summary.
type SearchServiceChargeSummary struct {
	ServiceCharges  []SearchServiceCharge `json:"service_charges"`
	PaxType         string                `json:"pax_type"`
	PaxCount        int                   `json:"pax_count"`
	BaseFare        float64               `json:"base_fare"`
	BaseTax         float64               `json:"base_tax"`
	BaseAmount      float64               `json:"base_amount"`
	BaseNta         float64               `json:"base_nta"`
	BasePrice       float64               `json:"base_price"`
	BaseCommission  float64               `json:"base_commission"`
	TotalFare       float64               `json:"total_fare"`
	TotalTax        float64               `json:"total_tax"`
	TotalNta        float64               `json:"total_nta"`
	TotalPrice      float64               `json:"total_price"`
	TotalCommission float64               `json:"total_commission"`
}

// SearchServiceCharge is a passenger-level service charge.
type SearchServiceCharge struct {
	PaxType           string  `json:"pax_type"`
	PaxCount          int     `json:"pax_count"`
	ChargeType        string  `json:"charge_type"`
	ChargeCode        string  `json:"charge_code"`
	Currency          string  `json:"currency"`
	Amount            float64 `json:"amount"`
	ForeignCurrency   string  `json:"foreign_currency"`
	ForeignAmount     float64 `json:"foreign_amount"`
	Total             float64 `json:"total"`
	Sequence          int     `json:"sequence"`
	IsData            bool    `json:"is_data,omitempty"`
	CommissionAgentID any     `json:"commission_agent_id,omitempty"`
}

// SearchRecommendation is a bundled fare recommendation.
type SearchRecommendation struct {
	JourneyFlightRefs    []SearchJourneyFlightRef     `json:"journey_flight_refs"`
	ServiceCharges       []SearchServiceCharge        `json:"service_charges"`
	ServiceChargeSummary []SearchServiceChargeSummary `json:"service_charge_summary"`
	Sequence             int                          `json:"sequence"`
}

// SearchJourneyFlightRef identifies a journey used by a recommendation.
type SearchJourneyFlightRef struct {
	JourneyRefID         string                       `json:"journey_ref_id"`
	FareFlightRefs       []SearchFareFlightRef        `json:"fare_flight_refs"`
	Origin               string                       `json:"origin"`
	Destination          string                       `json:"destination"`
	DepartureDate        string                       `json:"departure_date"`
	ArrivalDate          string                       `json:"arrival_date"`
	CarrierCode          string                       `json:"carrier_code"`
	CarrierNumber        string                       `json:"carrier_number"`
	CarrierName          string                       `json:"carrier_name"`
	ServiceChargeSummary []SearchServiceChargeSummary `json:"service_charge_summary,omitempty"`
}

// SearchFareFlightRef identifies a fare used by a recommendation.
type SearchFareFlightRef struct {
	FareRefID      string `json:"fare_ref_id"`
	ClassOfService string `json:"class_of_service"`
	FareBasisCode  string `json:"fare_basis_code"`
	FareName       string `json:"fare_name"`
	CarrierName    string `json:"carrier_name"`
}
