package train

import (
	"fmt"
	"strings"
	"time"
)

const (
	ProviderKAI = "kai"
	DirectionOW = "OW"
	DirectionRT = "RT"
)

// Credentials are used only for the Orbis /session request.
type Credentials struct {
	User     string `json:"user" log:"redact"`
	Password string `json:"password" log:"redact"`
	APIKey   string `json:"api_key" log:"redact"`
}

func (c Credentials) Validate() error {
	if strings.TrimSpace(c.User) == "" || strings.TrimSpace(c.Password) == "" || strings.TrimSpace(c.APIKey) == "" {
		return fmt.Errorf("orbis user, password, and api_key are required")
	}
	return nil
}

type ErrorCode int

const (
	ErrorSuccess      ErrorCode = 0
	ErrorInvalidInput ErrorCode = 400
	ErrorInternal     ErrorCode = 500
)

type ProviderError struct {
	Code    ErrorCode
	Message string
}

func (e *ProviderError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("Orbis train error (code %d): %s", e.Code, e.Message)
}

func newProviderError(code ErrorCode, message string) *ProviderError {
	return &ProviderError{Code: code, Message: message}
}

type BasicResult struct {
	ErrorCode ErrorCode `json:"error_code"`
	ErrorMsg  string    `json:"error_msg"`
}

type SignInResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  SignInResult `json:"result"`
}

type SignInResult struct {
	BasicResult
	Response struct {
		Signature string `json:"signature"`
	} `json:"response"`
}

// SearchRequest contains provider search parameters. The client obtains
// a fresh signature internally by calling SignIn before Search.
type SearchRequest struct {
	JourneyList []SearchJourney `json:"journey_list"`
	Direction   string          `json:"direction"`
	Adult       int             `json:"adult"`
	Infant      int             `json:"infant"`
	Provider    string          `json:"provider"`
}

type SearchJourney struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
}

func (r *SearchRequest) Validate() error {
	if r.Provider != ProviderKAI {
		return fmt.Errorf("provider must be %q", ProviderKAI)
	}
	if r.Adult < 1 {
		return fmt.Errorf("adult must be at least 1")
	}
	if r.Infant < 0 {
		return fmt.Errorf("infant cannot be negative")
	}

	var expectedJourneys int
	switch r.Direction {
	case DirectionOW:
		expectedJourneys = 1
	case DirectionRT:
		expectedJourneys = 2
	default:
		return fmt.Errorf("direction must be OW or RT")
	}
	if len(r.JourneyList) != expectedJourneys {
		return fmt.Errorf("%s search requires %d journey(s)", r.Direction, expectedJourneys)
	}

	for i, journey := range r.JourneyList {
		if strings.TrimSpace(journey.Origin) == "" || strings.TrimSpace(journey.Destination) == "" {
			return fmt.Errorf("journey %d requires origin and destination", i+1)
		}
		if strings.EqualFold(journey.Origin, journey.Destination) {
			return fmt.Errorf("journey %d origin and destination must differ", i+1)
		}
		if _, err := time.Parse("2006-01-02", journey.DepartureDate); err != nil {
			return fmt.Errorf("journey %d departure_date must be YYYY-MM-DD: %w", i+1, err)
		}
	}
	if r.Direction == DirectionRT {
		outbound, inbound := r.JourneyList[0], r.JourneyList[1]
		if !strings.EqualFold(outbound.Origin, inbound.Destination) ||
			!strings.EqualFold(outbound.Destination, inbound.Origin) {
			return fmt.Errorf("round-trip inbound journey must reverse outbound stations")
		}
		if inbound.DepartureDate < outbound.DepartureDate {
			return fmt.Errorf("return date cannot be before departure date")
		}
	}
	return nil
}

type SearchResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  SearchResult `json:"result"`
	// Signature is for backend session storage and future KAI booking operations.
	// It must not be serialized into frontend responses.
	Signature string `json:"-"`
}

type SearchResult struct {
	BasicResult
	Response SearchData `json:"response"`
}

type SearchData struct {
	Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
	JourneyList SearchJourney `json:"journey_list"`
	Journeys    []Journey     `json:"journeys"`
}

type Journey struct {
	JourneyCode    string `json:"journey_code"`
	Origin         string `json:"origin"`
	DepartureDate  string `json:"departure_date"`
	Destination    string `json:"destination"`
	ArrivalDate    string `json:"arrival_date"`
	ElapsedTime    string `json:"elapsed_time"`
	Provider       string `json:"provider"`
	Sequence       int    `json:"sequence"`
	CabinClass     string `json:"cabin_class"`
	ClassOfService string `json:"class_of_service"`
	AvailableCount int    `json:"available_count"`
	CarrierCode    string `json:"carrier_code"`
	CarrierName    string `json:"carrier_name"`
	CarrierNumber  string `json:"carrier_number"`
	Fares          []Fare `json:"fares"`
}

type Fare struct {
	FareCode             string                 `json:"fare_code"`
	ServiceChargeSummary []ServiceChargeSummary `json:"service_charge_summary"`
}

type ServiceChargeSummary struct {
	ServiceCharges  []ServiceCharge `json:"service_charges"`
	PaxType         string          `json:"pax_type"`
	PaxCount        int             `json:"pax_count"`
	BaseFare        float64         `json:"base_fare"`
	BaseTax         float64         `json:"base_tax"`
	BaseNTA         float64         `json:"base_nta"`
	BasePrice       float64         `json:"base_price"`
	BaseCommission  float64         `json:"base_commission"`
	TotalFare       float64         `json:"total_fare"`
	TotalTax        float64         `json:"total_tax"`
	TotalNTA        float64         `json:"total_nta"`
	TotalPrice      float64         `json:"total_price"`
	TotalCommission float64         `json:"total_commission"`
	TotalRAC        float64         `json:"total_rac"`
}

type ServiceCharge struct {
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
