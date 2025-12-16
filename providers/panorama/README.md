# Panorama Provider

A Go client library for interacting with the Panorama API, providing airline booking management functionality including authentication, booking retrieval, and booking cancellation.

## Features

- **Authentication**: Sign in to the Panorama API using credentials
- **Search**: Retrieve flight to the Panorama API
- **Booking Management**: Retrieve booking details by order number
- **Booking Cancellation**: Cancel bookings with "booked" status to "void" status
- **High-Level Usecase**: Convenient wrapper for complex operations like bulk booking cancellation
- **Error Handling**: Comprehensive error handling with custom error types
- **Observability**: Built-in OpenTelemetry tracing and structured logging
- **Retry Logic**: Configurable HTTP client with automatic retries

## Package Structure

```
providers/panorama/
├── README.md
├── client.go
├── endpoint.go
├── endpoint_test.go
├── errors.go
├── helpers.go
├── models.go
├── usecase.go
├── usecase_model.go
└── usecase_test.go
```

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/providers/panorama
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/panorama"
)

func main() {
    // Initialize logger
    logger := logger.New(logger.Config{
        Level:      "info",
        JSONOutput: false,
    })

    // Create client configuration
    config := panorama.Config{
        BaseURL: "https://api.panorama.example.com",
        Credentials: panorama.Credentials{
            APIKey:   "your-api-key",
            User:     "your-username",
            Password: "your-password",
        },
    }

    // Create client
    client, err := panorama.New(config, logger)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Sign in to get authentication signature
    signInResp, err := client.SignIn(ctx)
    if err != nil {
        log.Fatal(err)
    }

    signature := signInResp.Result.Response.Signature

    // Get booking details
    bookingResp, err := client.GetBooking(ctx, panorama.GetBookingRequest{
        Signature:   signature,
        OrderNumber: "ORDER123",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Booking state: %s", bookingResp.Result.Response.State)
}
```

## Configuration

### Client Configuration

```go
type Config struct {
    BaseURL     string                    // Base URL for the API (required)
    Credentials Credentials               // Authentication credentials (required)
    RetryConfig *client.RetryConfig       // HTTP retry configuration (optional)
}
```

### Credentials

```go
type Credentials struct {
    APIKey   string `json:"api_key"`
    User     string `json:"user"`
    Password string `json:"password"`
}
```

### Retry Configuration

The client uses a default retry configuration if none is provided:

```go
defaultRetryConfig := &client.RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    5 * time.Second,
    MaxBackoff:        15 * time.Second,
    BackoffMultiplier: 2.0,
}
```

## API Methods

### SignIn

Authenticates with the Panorama API and returns a signature for subsequent requests.

```go
func (c *Client) SignIn(ctx context.Context) (*SignInResponse, *PanoramaError)
```

**Response:**

```go
type SignInResponse struct {
    Result struct {
        BasicPanoramaResult
        Response struct {
            Signature string `json:"signature"`
        } `json:"response"`
    } `json:"result"`
}
```

### Search

Retrieve flight to the Panorama API

```go
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, *PanoramaError)
```

**Response:**

```go
type SearchResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SearchResult `json:"response"`
	} `json:"result"`
}
```

### GetBooking

Retrieves booking details by order number.

```go
func (c *Client) GetBooking(ctx context.Context, req GetBookingRequest) (*GetBookingResponse, *PanoramaError)
```

**Request:**

```go
type GetBookingRequest struct {
    Signature   string `json:"signature"`    // From SignIn response
    OrderNumber string `json:"order_number"` // Booking order number
}
```

**Response:**

```go
type GetBookingResponse struct {
    Result struct {
        BasicPanoramaResult
        Response BookingResponse `json:"response"`
    } `json:"result"`
}
```

### CancelVoid

Cancels a booking by changing its state from "booked" to "void".

```go
func (c *Client) CancelVoid(ctx context.Context, req CancelVoidRequest) (*GetBookingResponse, *PanoramaError)
```

**Request:**

```go
type CancelVoidRequest struct {
    Signature   string `json:"signature"`    // From SignIn response
    OrderNumber string `json:"order_number"` // Booking order number
}
```

## High-Level Usecase

The package provides a `Usecase` wrapper for complex operations that combine multiple API calls:

### Creating a Usecase Instance

```go
usecase := panorama.NewUsecase(client)
```

### SearchBooking

Search bookings, handling the complete workflow including authentication, and state validation.

**Workflow:**

1. Authenticates with the Panorama API (SignIn)
2. Search to get flight

```go
func (u *Usecase) Search(ctx context.Context, req SearchBookingRequest) (*SearchBookingResponse, *PanoramaError)
```

**Request:**

```go
type SearchBookingRequest struct {
	JourneyList []struct {
		Origin        string `json:"origin"`
		Destination   string `json:"destination"`
		DepartureDate string `json:"departure_date"`
	} `json:"journey_list"`
	Adult        int      `json:"adult"`
	Child        int      `json:"child"`
	Infant       int      `json:"infant"`
	CabinClass   string   `json:"cabin_class"`
	CarrierCodes []string `json:"carrier_codes"`
	IsComboPrice bool     `json:"is_combo_price"`
	Provider     string   `json:"provider"`
	PromoCodes   []string `json:"promo_codes"`
}
```

**Response:**

```go
type SearchBookingResponse struct {
	Signature string       `json:"signature"`
	Results   SearchResult `json:"results"`
}
```

**Example:**

```go
usecase := panorama.NewUsecase(client)

req := SearchBookingRequest{
    JourneyList: []struct {
        Origin        string "json:\"origin\""
        Destination   string "json:\"destination\""
        DepartureDate string "json:\"departure_date\""
    }{
        {
            Origin:        "CGK",
            Destination:   "SIN",
            DepartureDate: "2025-12-19",
        },
    },
    Adult:        1,
    Child:        0,
    Infant:       0,
    CabinClass:   "Y",
    CarrierCodes: []string{"SQ"},
    IsComboPrice: true,
    Provider:     "sia",
    PromoCodes:   []string{},
}

resp, err := usecase.Search(ctx, req)
if err != nil {
    log.Printf("Failed to process bulk cancellation: %v", err)
    return
}

log.Printf("✓ Successfully Search booking %s (final results: %v)", resp.Signature, resp.Results)
```

### CancelBooking

Cancels multiple bookings in bulk, handling the complete workflow including authentication, state validation, and cancellation.

**Workflow:**

1. Authenticates with the Panorama API (SignIn)
2. For each order number:
   - Retrieves booking details (GetBooking)
   - Validates booking state (must be "booked")
   - Cancels the booking if valid (CancelVoid)
   - Returns the result for each booking

```go
func (u *Usecase) CancelBooking(ctx context.Context, req CancelBookingRequest) (*CancelBookingResponse, error)
```

**Request:**

```go
type CancelBookingRequest struct {
    OrderNumbers []string `json:"order_numbers"` // List of order numbers to cancel
}
```

**Response:**

```go
type CancelBookingResponse struct {
    Results []CancelBookingResult `json:"results"`
}

type CancelBookingResult struct {
    OrderNumber string `json:"order_number"`
    Success     bool   `json:"success"`
    Message     string `json:"message"`
    State       string `json:"state,omitempty"` // Current/final state of the booking
}
```

**Example:**

```go
usecase := panorama.NewUsecase(client)

resp, err := usecase.CancelBooking(ctx, panorama.CancelBookingRequest{
    OrderNumbers: []string{"ORDER123", "ORDER456", "ORDER789"},
})
if err != nil {
    log.Printf("Failed to process bulk cancellation: %v", err)
    return
}

for _, result := range resp.Results {
    if result.Success {
        log.Printf("✓ Successfully canceled booking %s (final state: %s)", result.OrderNumber, result.State)
    } else {
        log.Printf("✗ Failed to cancel booking %s: %s (current state: %s)", result.OrderNumber, result.Message, result.State)
    }
}
```

**Notes:**

- Only bookings with state "booked" can be canceled
- Authentication is handled automatically for each batch
- Results are returned for all bookings, regardless of success/failure
- The operation is atomic per booking but not across multiple bookings

## API Reference

### Client Methods

#### SignIn

Authenticates with the Panorama API and returns a signature for subsequent requests.

```go
func (c *Client) SignIn(ctx context.Context) (*SignInResponse, *PanoramaError)
```

**Response:**

```go
type SignInResponse struct {
    Result struct {
        BasicPanoramaResult
        Response struct {
            Signature string `json:"signature"`
        } `json:"response"`
    } `json:"result"`
}
```

**Timeout:** 10 seconds

#### SearchBooking

Retrieve flight to the Panorama API.

```go
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, *PanoramaError)
```

**Request Parameters:**

- `Signature` (string, required): Authentication signature from SignIn
- `SearchRequest` (struct, required): Booking search to retrieve

**Response:**

```go
type SearchResponse struct {
	Result struct {
		BasicPanoramaResult
		Response SearchResult `json:"response"`
	} `json:"result"`
}
```

**Timeout:** 59 seconds

#### GetBooking

Retrieves booking details by order number.

```go
func (c *Client) GetBooking(ctx context.Context, req GetBookingRequest) (*GetBookingResponse, *PanoramaError)
```

**Request Parameters:**

- `Signature` (string, required): Authentication signature from SignIn
- `OrderNumber` (string, required): Booking order number to retrieve

**Response:**

```go
type GetBookingResponse struct {
    Result struct {
        BasicPanoramaResult
        Response BookingResponse `json:"response"`
    } `json:"result"`
}
```

**Timeout:** 10 seconds

#### CancelVoid

Cancels a booking by changing its state from "booked" to "void".

```go
func (c *Client) CancelVoid(ctx context.Context, req CancelVoidRequest) (*GetBookingResponse, *PanoramaError)
```

**Request Parameters:**

- `Signature` (string, required): Authentication signature from SignIn
- `OrderNumber` (string, required): Booking order number to cancel

**Response:** Returns the updated booking details after cancellation.

**Timeout:** 10 seconds

### Usecase Methods

#### CancelBooking

Cancels multiple bookings in bulk with automatic authentication and state validation.

```go
func (u *Usecase) CancelBooking(ctx context.Context, req CancelBookingRequest) (*CancelBookingResponse, error)
```

**Request Parameters:**

- `OrderNumbers` ([]string, required): List of booking order numbers to cancel

**Response:**

- `Results` ([]CancelBookingResult): Individual results for each booking cancellation attempt
  - `OrderNumber` (string): The booking order number
  - `Success` (bool): Whether the cancellation was successful
  - `Message` (string): Success message or error description
  - `State` (string): Final state of the booking after the operation

## Error Handling

The package defines custom error types for different scenarios:

### PanoramaError

```go
type PanoramaError struct {
    Result PanoramaErrorResult `json:"result"`
}

type PanoramaErrorResult struct {
    BasicPanoramaResult
    Response           string `json:"response"`
    ErrorAdditionalMsg string `json:"error_additional_msg"`
}
```

### Error Codes

```go
const (
    InvalidRequest        PanoramaErrorCode = 400   // Invalid request data
    InternalServerError   PanoramaErrorCode = 500   // Server or network error
    Success               PanoramaErrorCode = 0     // Operation successful
    InvalidActionType     PanoramaErrorCode = 4001  // Invalid action type
    AccessDenied          PanoramaErrorCode = 4002  // Authentication failed
    SessionExpired        PanoramaErrorCode = 4003  // Your session has been expired
    PermissionDenied      PanoramaErrorCode = 4004  // Permission denied
	ProviderUnavailable   PanoramaErrorCode = 4005  // Provider is unavailable
	PriceChanged          PanoramaErrorCode = 4006  // Price changed
	SsrUnavailable        PanoramaErrorCode = 4007  // SSR not available
	SeatUnavailable       PanoramaErrorCode = 4008  // Seat not available
	FailedSsr             PanoramaErrorCode = 4009  // Failed to assign SSR
	FailedSeat            PanoramaErrorCode = 4010  // Failed to assign Seat
	MissingRequirement    PanoramaErrorCode = 4011  // Missing Requirements
	DuplicateSession      PanoramaErrorCode = 4012  // Duplicate Session, Please create new session
	ValidationError       PanoramaErrorCode = 4013  // Data Validation Error
	AlreadyCall           PanoramaErrorCode = 4014  // You had already call this function
	NotAllowed            PanoramaErrorCode = 4015  // Destinations are not allowed
    FailedCommitBooking   PanoramaErrorCode = 4016  // Failed to commit booking
	FailedGetBooking      PanoramaErrorCode = 4017  // Failed to get booking
	AccessRestriction     PanoramaErrorCode = 4018  // Access restriction
	FailedReissue         PanoramaErrorCode = 4019  // Failed to Reissue
	NoSchedules           PanoramaErrorCode = 4020  // No Schedules
	PriceUnavailable      PanoramaErrorCode = 4021  // Price is not available. Please choose another fare class combination
	FailedCancelBooking   PanoramaErrorCode = 4022  // Cancel booking failed
	ReservationNotFound   PanoramaErrorCode = 4023  // Reservation not found
	VendorSessionExpired  PanoramaErrorCode = 4024  // Vendor session expired
	FlightUnavailable     PanoramaErrorCode = 4025  // Flight is not available
	HoldBookingRestricted PanoramaErrorCode = 4026  // Hold booking restricted, You have to Direct Issued this reservation
	NotBalance            PanoramaErrorCode = 4027  // Not enough vendor balance
	CannotInfant          PanoramaErrorCode = 4028  // Restricted booking, infant cannot be paired with senior citizen (YCD)
	FailedIssued          PanoramaErrorCode = 4029  // Failed issued
	UserBanned            PanoramaErrorCode = 4030  // User is banned
	ResourceNotFound      PanoramaErrorCode = 4031  // Resource not found
	FailedTransaction     PanoramaErrorCode = 4032  // Failed Transaction
	FailedPayment         PanoramaErrorCode = 4033  // Failed to execute payment
	TransactionPaid       PanoramaErrorCode = 4034  // Transaction has been paid on vendor
	TransactionManual     PanoramaErrorCode = 4035  // Failed Transaction, manual resync needed
	TransactionProgress   PanoramaErrorCode = 4036  // Transaction is still in progress, please try again later
	FfUnavailable         PanoramaErrorCode = 4037  // Frequent Flyer is not available
	PnrNotFound           PanoramaErrorCode = 4038  // PNR is not found
	PnrNotAllowed         PanoramaErrorCode = 4039  // PNR not allowed to process
	FailedSplitBooking    PanoramaErrorCode = 4040  // Failed Split Booking
	FailedVendorBalance   PanoramaErrorCode = 4041  // Failed to retrieve vendor balance
	NoUpdateDetected      PanoramaErrorCode = 4042  // No update change detected
	CancelUnavailable     PanoramaErrorCode = 4043  // Cancel is not available
	SplitUnavailable      PanoramaErrorCode = 4044  // Split is not available
)
```

**Example Error Handling:**

```go
resp, err := client.GetBooking(ctx, request)
if err != nil {
    if panoramaErr, ok := err.(*panorama.PanoramaError); ok {
        log.Printf("Panorama API error (code %d): %s", panoramaErr.Result.ErrorCode, panoramaErr.Result.ErrorMsg)
    } else {
        log.Printf("Other error: %v", err)
    }
    return
}
```

## Observability

The client includes built-in observability features:

- **OpenTelemetry Tracing**: All API calls are traced with spans
- **Structured Logging**: Requests and responses are logged with context
- **Request/Response Logging**: Sensitive data is automatically redacted

## Dependencies

This package depends on other packages in the tixia-pkg monorepo:

- `github.com/TixiaOTA/tixia-pkg/http/client` - HTTP client with retry logic
- `github.com/TixiaOTA/tixia-pkg/logger` - Structured logging
- `github.com/TixiaOTA/tixia-pkg/otel` - OpenTelemetry integration

## Contributing

1. Follow the existing code style and patterns
2. Add tests for new functionality
3. Update documentation for API changes
4. Ensure all tests pass before submitting PR

## License

This package is part of the TixiaOTA internal package collection.
