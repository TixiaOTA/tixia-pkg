package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- CONSTANTS ---------------------------------

// PanoramaHeaderAction is the action header for the Panorama API
type PanoramaHeaderAction string

const (
	SignIn     PanoramaHeaderAction = "signin"
	Search     PanoramaHeaderAction = "search"
	GetBooking PanoramaHeaderAction = "get_booking"
	Cancel     PanoramaHeaderAction = "cancel"
)

// --------------------------------- REQUESTS ---------------------------------

// SearchRequest is the request structure for Search
type SearchRequest struct {
	Signature    string               `json:"signature,omitempty"`
	JourneyList  []JourneyListRequest `json:"journey_list"`
	Adult        int                  `json:"adult"`
	Child        int                  `json:"child"`
	Infant       int                  `json:"infant"`
	CabinClass   string               `json:"cabin_class"`
	CarrierCodes []string             `json:"carrier_codes"`
	IsComboPrice bool                 `json:"is_combo_price"`
	Provider     string               `json:"provider"`
	PromoCodes   []string             `json:"promo_codes"`
}

// JourneyListRequest represents the request structure for flight search
type JourneyListRequest struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	Pax           *[]int `json:"pax,omitempty"`
}

func (r *SearchRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if len(r.JourneyList) <= 0 {
		return errors.New("journey list is required")
	}
	for _, v := range r.JourneyList {
		if v.Origin == "" {
			return errors.New("origin is required")
		}
		if v.Destination == "" {
			return errors.New("destination is required")
		}
		if v.DepartureDate == "" {
			return errors.New("departure date is required")
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

// GetBookingRequest is the request structure for GetBooking
type GetBookingRequest struct {
	Signature   string `json:"signature"`
	OrderNumber string `json:"order_number"`
}

func (r *GetBookingRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	return nil
}

// CancelVoidRequest is the request structure for CancelVoid
type CancelVoidRequest struct {
	Signature   string `json:"signature"`
	OrderNumber string `json:"order_number"`
}

func (r *CancelVoidRequest) Validate() error {
	if r.Signature == "" {
		return errors.New("signature is required")
	}
	if r.OrderNumber == "" {
		return errors.New("order_number is required")
	}
	return nil
}

// --------------------------------- RESPONSES ---------------------------------

// SignInResponse is the response structure for SignIn
type SignInResponse struct {
	Result struct {
		BasicPanoramaResult
		Response struct {
			Signature string `json:"signature"`
		} `json:"response"`
	} `json:"result"`
}

type SearchResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SearchResult `json:"response"`
	} `json:"result"`
}

// GetBookingResponse is the response structure for GetBooking
type GetBookingResponse struct {
	Result struct {
		BasicPanoramaResult
		Response BookingResponse `json:"response"`
	} `json:"result"`
}

// --------------------------------- SHARED TYPES ---------------------------------

type Credentials struct {
	APIKey   string `json:"api_key" log:"redact"`
	User     string `json:"user" log:"redact"`
	Password string `json:"password" log:"redact"`
}

func (c *Credentials) Validate() error {
	if c.APIKey == "" {
		return errors.New("api_key is required")
	}
	if c.User == "" {
		return errors.New("user is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// APIError is the error structure for all API errors
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (code %d): %s", e.Code, e.Message)
}

// BasicPanoramaResult is the result structure for all Panorama results
type BasicPanoramaResult struct {
	ErrorCode PanoramaErrorCode `json:"error_code"`
	ErrorMsg  string            `json:"error_msg"`
}

// PanoramaError is the error structure for all Panorama errors
type PanoramaError struct {
	Result PanoramaErrorResult `json:"result"`
}

// PanoramaErrorResult is the result structure for all Panorama errors
type PanoramaErrorResult struct {
	BasicPanoramaResult
	Response           string `json:"response"`
	ErrorAdditionalMsg string `json:"error_additional_msg"`
}

func (e *PanoramaError) Error() string {
	return fmt.Sprintf("Panorama error (code %d): %s", e.Result.ErrorCode, e.Result.ErrorMsg)
}

type BookingResponse struct {
	OrderNumber string `json:"order_number"`
	State       string `json:"state"`
	// Add more fields as needed
}

type SearchResult struct {
	Schedules []struct {
		Journeys []struct {
			Origin             string   `json:"origin"`
			OriginCountry      string   `json:"origin_country"`
			DepartureDate      string   `json:"departure_date"`
			Destination        string   `json:"destination"`
			DestinationCountry string   `json:"destination_country"`
			ArrivalDate        string   `json:"arrival_date"`
			JourneyCode        string   `json:"journey_code"`
			RouteCode          string   `json:"route_code"`
			ElapsedTime        string   `json:"elapsed_time"`
			TransitCount       int      `json:"transit_count"`
			Provider           string   `json:"provider"`
			CarrierCodeList    []string `json:"carrier_code_list"`
			Segments           []struct {
				Origin             string `json:"origin"`
				OriginCountry      string `json:"origin_country"`
				DepartureDate      string `json:"departure_date"`
				Destination        string `json:"destination"`
				DestinationCountry string `json:"destination_country"`
				ArrivalDate        string `json:"arrival_date"`
				TransitCount       int    `json:"transit_count"`
				SegmentCode        string `json:"segment_code"`
				CarrierName        string `json:"carrier_name"`
				CarrierCode        string `json:"carrier_code"`
				CarrierTypeName    string `json:"carrier_type_name"`
				Provider           string `json:"provider"`
				Legs               []struct {
					Origin             string `json:"origin"`
					OriginCountry      string `json:"origin_country"`
					DepartureDate      string `json:"departure_date"`
					Destination        string `json:"destination"`
					DestinationCountry string `json:"destination_country"`
					ArrivalDate        string `json:"arrival_date"`
					ElapsedTime        string `json:"elapsed_time"`
					CarrierName        string `json:"carrier_name"`
					CarrierCode        string `json:"carrier_code"`
					CarrierTypeName    string `json:"carrier_type_name"`
					Provider           string `json:"provider"`
					Sequence           int    `json:"sequence"`
				} `json:"legs"`
				Fares []struct {
					CabinClass     string `json:"cabin_class"`
					ClassOfService string `json:"class_of_service"`
					FareCode       string `json:"fare_code"`
					AvailableCount int    `json:"available_count"`
					FareDetails    []struct {
						DetailCode  string        `json:"detail_code"`
						DetailType  string        `json:"detail_type"`
						Amount      float64       `json:"amount"`
						Unit        string        `json:"unit"`
						DetailName  string        `json:"detail_name"`
						Description []interface{} `json:"description"`
						Sequence    int           `json:"sequence"`
					} `json:"fare_details"`
					FareRefID            string `json:"fare_ref_id"`
					IsTravelFair         bool   `json:"is_travel_fair"`
					ServiceChargeSummary []struct {
						PaxType   string  `json:"pax_type"`
						BaseFare  float64 `json:"base_fare"`
						BaseTax   float64 `json:"base_tax"`
						BasePrice float64 `json:"base_price"`
						PaxCount  int     `json:"pax_count"`
					} `json:"service_charge_summary"`
					Sequence int `json:"sequence"`
				} `json:"fares"`
				Status     string        `json:"status"`
				SeatCabins []interface{} `json:"seat_cabins"`
				Sequence   int           `json:"sequence"`
			} `json:"segments"`
			JourneyRefID string `json:"journey_ref_id"`
			IsVtlFlight  bool   `json:"is_vtl_flight"`
			Sequence     int    `json:"sequence"`
		} `json:"journeys"`
		Provider      string `json:"provider"`
		Origin        string `json:"origin"`
		Destination   string `json:"destination"`
		DepartureDate string `json:"departure_date"`
	} `json:"schedules"`
	Recommendations []struct {
		JourneyFlightRefs []struct {
			JourneyRefID   string `json:"journey_ref_id"`
			FareFlightRefs []struct {
				FareRefID string `json:"fare_ref_id"`
			} `json:"fare_flight_refs"`
		} `json:"journey_flight_refs"`
	} `json:"recommendations"`
	JourneyKey string `json:"journey_key"`
	Provider   string `json:"provider"`
	// Add more fields as needed
}
