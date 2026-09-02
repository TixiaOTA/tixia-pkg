package panorama

import (
	"errors"
	"fmt"
)

// --------------------------------- CONSTANTS ---------------------------------

// PanoramaHeaderAction is the action header for the Panorama API
type PanoramaHeaderAction string

const (
	SignIn              PanoramaHeaderAction = "signin"
	GetBooking          PanoramaHeaderAction = "get_booking"
	Cancel              PanoramaHeaderAction = "cancel"
	GetPriceItinerary   PanoramaHeaderAction = "get_price_itinerary"
	SellJourneys        PanoramaHeaderAction = "sell_journeys"
	GetBalance          PanoramaHeaderAction = "get_balance"
	UpdateContacts      PanoramaHeaderAction = "update_contacts"
	UpdatePassengers    PanoramaHeaderAction = "update_passengers"
	CommitBooking       PanoramaHeaderAction = "commit_booking"
	Issued              PanoramaHeaderAction = "issued"
	GetFFAvailability   PanoramaHeaderAction = "get_ff_availability"
	GetSSRAvailability  PanoramaHeaderAction = "get_ssr_availability"
	GetSeatAvailability PanoramaHeaderAction = "get_seat_availability"
	SellSSRs            PanoramaHeaderAction = "sell_ssrs"
	AssignSeats         PanoramaHeaderAction = "assign_seats"
)

// --------------------------------- REQUESTS ---------------------------------

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
	Code    PanoramaErrorCode `json:"code"`
	Message string            `json:"message"`
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
