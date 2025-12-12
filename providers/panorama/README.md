# Panorama Provider

A Go client library for interacting with the Panorama API, providing airline booking management functionality including authentication, booking retrieval, and booking cancellation.

## Features

- **Authentication**: Sign in to the Panorama API using credentials
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
    InvalidRequest      PanoramaErrorCode = 400   // Invalid request data
    InternalServerError PanoramaErrorCode = 500   // Server or network error
    Success             PanoramaErrorCode = 0     // Operation successful
    AccessDenied        PanoramaErrorCode = 4002  // Authentication failed
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
