# MG Package

A Go client library for interacting with the MG Hotel Booking API. This package provides methods for searching hotels, rechecking availability, booking hotels, and retrieving reservation details.

## Overview

The MG package enables you to:

- Search for hotels with availability and pricing
- Recheck hotel availability and pricing using session IDs
- Book hotels with guest details
- Retrieve reservation/booking details

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/providers/mg
```

## Configuration

### Client Configuration

Create a client instance with the required configuration:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/providers/mg"
    "github.com/TixiaOTA/tixia-pkg/http/client"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

config := mg.Config{
    BaseURL: "https://api.mg.com", // Required
    Credentials: mg.Credentials{
        AgencyCode: "your-agency-code", // Required
        Username:   "your-username",    // Required
        Password:   "your-password",    // Required
    },
    RetryConfig: &client.RetryConfig{ // Optional
        MaxRetries:        3,
        InitialBackoff:    5 * time.Second,
        MaxBackoff:        15 * time.Second,
        BackoffMultiplier: 2.0,
    },
}

logger := logger.New() // Create your logger instance

client, err := mg.New(config, logger)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Fields

- **BaseURL** (string, required): Base URL for the MG API
- **Credentials** (Credentials, required): Authentication credentials containing:
  - `AgencyCode` (string, required): Agency code for authentication
  - `Username` (string, required): Username for authentication
  - `Password` (string, required): Password for authentication
- **RetryConfig** (\*client.RetryConfig, optional): Retry configuration for HTTP requests. Defaults to 3 retries with exponential backoff (5s initial, 15s max, 2x multiplier)

## DTO Pattern

This package uses Data Transfer Objects (DTOs) to separate external API request structures from internal business models. DTOs are defined in `dto.go` and provide:

- **API-specific serialization**: DTOs handle JSON marshaling/unmarshaling for API communication
- **Conversion methods**: Each DTO includes a `ToInternalModel()` method to convert to internal types
- **Lightweight structures**: DTOs contain only the fields needed for API requests

When calling client methods, use the DTO types (e.g., `SearchHotelRequestDTO`) instead of the internal model types. The client automatically converts DTOs to internal models for processing.

## Usage

### Search Hotel

Search for hotels with availability and pricing information:

```go
ctx := context.Background()

req := mg.SearchHotelRequestDTO{
    Nationality: "US",
    Country:     "US",
    City:        "New York",
    CheckIn:     "2024-12-01", // Format: YYYY-MM-DD
    CheckOut:    "2024-12-05", // Format: YYYY-MM-DD
    Rooms: &struct {
        Room []mg.RoomRequest `json:"Room"`
    }{
        Room: []mg.RoomRequest{
            {
                RoomNo:     1,
                NoOfAdults: 2,
                NoOfChild:  0,
                ExtraBed:   false,
            },
        },
    },
    Currency:     "USD",
    Language:     "en",
    AvailFlag:    true,
    MaxNoOfHotel: 10,
    DetailLevel:  "Full",
}

response, mgErr := client.SearchHotel(ctx, req)
if mgErr != nil {
    log.Fatal(mgErr)
}

fmt.Printf("Found %d hotels\n", response.NoOfHotels)
fmt.Printf("Session ID: %s\n", response.SessionID)

for _, hotel := range response.Hotels.Hotel {
    fmt.Printf("Hotel: %s (Code: %s)\n", hotel.Name, hotel.Code)
    for _, roomDetail := range hotel.RoomDetails {
        fmt.Printf("  Room: %s - Price: %.2f %s\n", roomDetail.Name, roomDetail.NetPrice, response.Currency)
    }
}
```

### Recheck Hotel

Recheck hotel availability and pricing using a session ID from a previous search:

```go
req := mg.RecheckHotelRequestDTO{
    SessionID:   "session-id-from-search", // From SearchHotel response
    Nationality: "US",
    Country:     "US",
    City:        "New York",
    HotelCode:   "HOTEL123", // Optional
    CheckIn:     "2024-12-01",
    CheckOut:    "2024-12-05",
    RoomDetails: &mg.RoomDetailsRequest{
        Code:                   "ROOM123",
        MealPlan:               "BB",
        CancellationPolicyType: "FreeCancellation",
        PackageRate:            false,
    },
    Rooms: struct {
        Room []mg.RoomRequest `json:"Room"`
    }{
        Room: []mg.RoomRequest{
            {
                RoomNo:     1,
                NoOfAdults: 2,
                NoOfChild:  0,
                ExtraBed:   false,
                RateKey:    "rate-key-from-search", // From SearchHotel response
            },
        },
    },
    Currency:    "USD",
    Language:    "en",
    AvailFlag:   true,
    DetailLevel: "Full",
}

response, mgErr := client.RecheckHotel(ctx, req)
if mgErr != nil {
    log.Fatal(mgErr)
}

fmt.Printf("Recheck Status: %t\n", response.Status)
fmt.Printf("Hotel: %s\n", response.Hotels.Hotel.Name)
```

### Book Hotel

Book a hotel with guest details:

```go
req := mg.BookHotelRequestDTO{
    SessionID:       "session-id-from-search",
    AgencyBookingID: "TIXIA-BOOKING-12345", // Your generated booking ID
    Nationality:     "US",
    CheckIn:         "2024-12-01",
    CheckOut:        "2024-12-05",
    HotelCode:       "HOTEL123",
    RoomDetails: mg.RoomDetailsRequest{
        Code:                   "ROOM123",
        MealPlan:               "BB",
        CancellationPolicyType: "FreeCancellation",
        PackageRate:            false,
    },
    Rooms: struct {
        Room []mg.RoomBookingRequest `json:"Room"`
    }{
        Room: []mg.RoomBookingRequest{
            {
                RoomNo:     1,
                NoOfAdults: 2,
                NoOfChild:  0,
                ExtraBed:   false,
                RateKey:    "rate-key-from-search",
                PaxDetails: struct {
                    Pax []mg.Pax `json:"Pax"`
                }{
                    Pax: []mg.Pax{
                        {
                            Salutation: "Mr",
                            FirstName:  "John",
                            LastName:   "Doe",
                            Age:        30,
                        },
                        {
                            Salutation: "Mrs",
                            FirstName:  "Jane",
                            LastName:   "Doe",
                            Age:        28,
                        },
                    },
                },
            },
        },
    },
    Currency:    "USD",
    Language:    "en",
    AvailFlag:   true,
    OnHold:      false,
    SpecialReq:  "Late check-in requested",
    DetailLevel: "Full",
}

response, mgErr := client.BookHotel(ctx, req)
if mgErr != nil {
    log.Fatal(mgErr)
}

fmt.Printf("Booking Status: %t\n", response.Status)
fmt.Printf("MG Booking ID: %s\n", response.BookingDetail.MgBookingID)
fmt.Printf("Agency Booking ID: %s\n", response.BookingDetail.AgencyBookingID)
```

**Note:** BookHotel has a 120-second timeout and retries are disabled (no automatic retries) to prevent duplicate bookings.

### Get Reservation Details

Retrieve reservation/booking details:

```go
req := mg.GetRSVNDetailsRequestDTO{
    MGBookingID:     "MG-BOOKING-12345",     // MG's booking ID
    AgencyBookingID: "TIXIA-BOOKING-12345",  // Your booking ID
}

response, mgErr := client.GetRSVNDetails(ctx, req)
if mgErr != nil {
    log.Fatal(mgErr)
}

fmt.Printf("Status: %s\n", response.BookingDetail.Status)
fmt.Printf("Check-in: %s\n", response.BookingDetail.CheckIn)
fmt.Printf("Check-out: %s\n", response.BookingDetail.CheckOut)
fmt.Printf("Hotel: %s\n", response.BookingDetail.Hotels.Hotel.Name)
```

## API Reference

### Client Methods

#### SearchHotel

Searches for hotels with availability and pricing information.

```go
func (c *Client) SearchHotel(ctx context.Context, req SearchHotelRequestDTO) (*SearchHotelResponse, *MGError)
```

**Request Parameters:**

- `Login` (Credentials, required): Authentication credentials
- `Nationality` (string, required): Guest nationality code
- `Country` (string, required): Country code
- `City` (string, required): City name
- `Hotels` (\*struct, optional): Specific hotel codes to search
- `CheckIn` (string, required): Check-in date (YYYY-MM-DD)
- `CheckOut` (string, required): Check-out date (YYYY-MM-DD)
- `Rooms` (\*struct, required): Room configuration with guest details
- `Currency` (string, required): Currency code
- `Language` (string, required): Language code
- `AvailFlag` (bool, required): Availability flag
- `MaxNoOfHotel` (int, required): Maximum number of hotels to return
- `DetailLevel` (string, required): Detail level (e.g., "Full")

**Response:**

- `Status` (bool): Request status
- `SessionID` (string): Session ID for subsequent operations
- `NoOfHotels` (int): Number of hotels found
- `Hotels` (struct): List of hotels with room details
- `CheckIn` (string): Check-in date
- `CheckOut` (string): Check-out date
- `Currency` (string): Currency used
- `ErrorCode` (\*string): Error code if status is false
- `ErrorMessage` (\*string): Error message if status is false

**Timeout:** 10 seconds

#### RecheckHotel

Rechecks hotel availability and pricing using a session ID.

```go
func (c *Client) RecheckHotel(ctx context.Context, req RecheckHotelRequestDTO) (*RecheckHotelResponse, *MGError)
```

**Request Parameters:**

- `Login` (Credentials, required): Authentication credentials
- `SessionID` (string, required): Session ID from SearchHotel response
- `Nationality` (string, required): Guest nationality code
- `Country` (string, required): Country code
- `City` (string, required): City name
- `HotelCode` (string, optional): Specific hotel code
- `CheckIn` (string, required): Check-in date (YYYY-MM-DD)
- `CheckOut` (string, required): Check-out date (YYYY-MM-DD)
- `RoomDetails` (\*RoomDetailsRequest, required): Room details configuration
- `Rooms` (struct, required): Room configuration with rate keys
- `Currency` (string, required): Currency code
- `Language` (string, optional): Language code
- `AvailFlag` (bool, optional): Availability flag
- `DetailLevel` (string, optional): Detail level

**Response:**

- `Status` (bool): Request status
- `SessionID` (string): Session ID
- `Hotels` (struct): Hotel information with updated pricing
- `CheckIn` (string): Check-in date
- `CheckOut` (string): Check-out date
- `Currency` (string): Currency used
- `ErrorCode` (\*string): Error code if status is false
- `ErrorMessage` (\*string): Error message if status is false

**Timeout:** 10 seconds

#### BookHotel

Books a hotel with guest details.

```go
func (c *Client) BookHotel(ctx context.Context, req BookHotelRequestDTO) (*BookHotelResponse, *MGError)
```

**Request Parameters:**

- `Login` (Credentials, required): Authentication credentials
- `SessionID` (string, required): Session ID from SearchHotel response
- `AgencyBookingID` (string, required): Your generated booking ID
- `Nationality` (string, required): Guest nationality code
- `CheckIn` (string, required): Check-in date (YYYY-MM-DD)
- `CheckOut` (string, required): Check-out date (YYYY-MM-DD)
- `HotelCode` (string, required): Hotel code
- `RoomDetails` (RoomDetailsRequest, required): Room details configuration
- `Rooms` (struct, required): Room configuration with guest details (Pax)
- `Currency` (string, required): Currency code
- `Language` (string, optional): Language code
- `AvailFlag` (bool, optional): Availability flag
- `OnHold` (bool, optional): Whether to hold the booking
- `SpecialReq` (string, optional): Special requests
- `DetailLevel` (string, optional): Detail level

**Response:**

- `Status` (bool): Booking status
- `SessionID` (string): Session ID
- `BookingDetail` (BookingDetails): Complete booking details including MG booking ID

**Timeout:** 120 seconds  
**Retries:** Disabled (0 retries) to prevent duplicate bookings

#### GetRSVNDetails

Retrieves reservation/booking details.

```go
func (c *Client) GetRSVNDetails(ctx context.Context, req GetRSVNDetailsRequestDTO) (*GetRSVNDetailsResponse, *MGError)
```

**Request Parameters:**

- `Login` (Credentials, required): Authentication credentials
- `MGBookingID` (string, required): MG's booking ID
- `AgencyBookingID` (string, required): Your booking ID

**Response:**

- `Status` (bool): Request status
- `SessionID` (string): Session ID
- `BookingDetail` (BookingDetails): Complete booking/reservation details

**Timeout:** 10 seconds

## Error Handling

The package returns `*MGError` for API-related errors:

```go
type MGError struct {
    Status       bool
    ErrorCode    MGErrorCode
    ErrorMessage string
}

func (e *MGError) Error() string {
    return fmt.Sprintf("MG error (status %t): %s - %s", e.Status, e.ErrorCode, e.ErrorMessage)
}
```

### Error Codes

- `InternalServerError`: Internal server error (network or server issues)
- `JRVXML161`: Intermittent server error (must retry)
- `JRVXML013`: Invalid request (request body is invalid)
- `JRVXML112`: AgencyBookingID already exists

**Example:**

```go
response, mgErr := client.SearchHotel(ctx, req)
if mgErr != nil {
    switch mgErr.ErrorCode {
    case mg.JRVXML161:
        // Retry the request
        fmt.Println("Intermittent error, retrying...")
    case mg.JRVXML013:
        // Fix the request
        fmt.Printf("Invalid request: %s\n", mgErr.ErrorMessage)
    default:
        fmt.Printf("Error: %v\n", mgErr)
    }
    return
}
```

For HTTP errors (non-200 status codes), the package returns `*APIError`:

```go
type APIError struct {
    Code    int
    Message string
}
```

## Data Models

### Credentials

Authentication credentials used in requests:

```go
type Credentials struct {
    AgencyCode string `json:"AgencyCode"`
    Username   string `json:"Username"`
    Password   string `json:"Password"`
}
```

### RoomRequest

Room configuration for search and recheck operations:

```go
type RoomRequest struct {
    RoomNo     uint   `json:"RoomNo"`
    NoOfAdults uint   `json:"NoOfAdults"`
    NoOfChild  uint   `json:"NoOfChild"`
    Child1Age  uint   `json:"Child1Age"`
    Child2Age  uint   `json:"Child2Age"`
    ExtraBed   bool   `json:"ExtraBed"`
    RateKey    string `json:"RateKey,omitempty"`
}
```

### RoomBookingRequest

Room configuration with guest details for booking:

```go
type RoomBookingRequest struct {
    PaxDetails struct {
        Pax []Pax `json:"Pax"`
    } `json:"PaxDetails"`
    RoomNo     uint   `json:"RoomNo"`
    NoOfAdults uint   `json:"NoOfAdults"`
    NoOfChild  uint   `json:"NoOfChild"`
    ExtraBed   bool   `json:"ExtraBed"`
    RateKey    string `json:"RateKey"`
}
```

### Pax

Guest/passenger information:

```go
type Pax struct {
    Salutation string `json:"Salutation"`
    FirstName  string `json:"FirstName"`
    LastName   string `json:"LastName"`
    Age        uint   `json:"Age"`
}
```

### BookingDetails

Complete booking information:

```go
type BookingDetails struct {
    Nationality       string
    CheckIn           string
    CheckOut          string
    Currency          string
    MgBookingID       string
    AgencyBookingID   string
    MgBookingVerID    string
    AgencyVoucherNo   string
    AgencyVoucherDate string
    Status            string
    SpecialRequest    CDataWrapper
    Hotels            struct {
        Hotel HotelRecheckResponse `json:"hotel"`
    }
}
```

See `models.go` for complete model definitions.

## Examples

### Complete Booking Flow Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/TixiaOTA/tixia-pkg/providers/mg"
    "github.com/TixiaOTA/tixia-pkg/logger"
)

func main() {
    // Create client
    config := mg.Config{
        BaseURL: "https://api.mg.com",
        Credentials: mg.Credentials{
            AgencyCode: "your-agency-code",
            Username:   "your-username",
            Password:   "your-password",
        },
    }

    logger := logger.New()
    client, err := mg.New(config, logger)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Step 1: Search for hotels
    searchReq := mg.SearchHotelRequest{
        Login:       config.Credentials,
        Nationality: "US",
        Country:     "US",
        City:        "New York",
        CheckIn:     "2024-12-01",
        CheckOut:    "2024-12-05",
        Rooms: &struct {
            Room []mg.RoomRequest `json:"Room"`
        }{
            Room: []mg.RoomRequest{
                {
                    RoomNo:     1,
                    NoOfAdults: 2,
                    NoOfChild:  0,
                    ExtraBed:   false,
                },
            },
        },
        Currency:     "USD",
        Language:     "en",
        AvailFlag:    true,
        MaxNoOfHotel: 10,
        DetailLevel:  "Full",
    }

    searchResp, mgErr := client.SearchHotel(ctx, searchReq)
    if mgErr != nil {
        log.Fatal(mgErr)
    }

    fmt.Printf("Found %d hotels\n", searchResp.NoOfHotels)

    if len(searchResp.Hotels.Hotel) == 0 {
        fmt.Println("No hotels found")
        return
    }

    // Step 2: Select a hotel and room
    selectedHotel := searchResp.Hotels.Hotel[0]
    selectedRoom := selectedHotel.RoomDetails[0]
    rateKey := selectedRoom.Rooms.Room[0].RateKey

    // Step 3: Recheck the hotel
    recheckReq := mg.RecheckHotelRequest{
        Login:       config.Credentials,
        SessionID:   searchResp.SessionID,
        Nationality: "US",
        Country:     "US",
        City:        "New York",
        HotelCode:   selectedHotel.Code,
        CheckIn:     "2024-12-01",
        CheckOut:    "2024-12-05",
        RoomDetails: &mg.RoomDetailsRequest{
            Code:                   selectedRoom.Code,
            MealPlan:               selectedRoom.MealPlan,
            CancellationPolicyType: selectedRoom.CancellationPolicyType,
            PackageRate:            selectedRoom.PackageRate,
        },
        Rooms: struct {
            Room []mg.RoomRequest `json:"Room"`
        }{
            Room: []mg.RoomRequest{
                {
                    RoomNo:     1,
                    NoOfAdults: 2,
                    NoOfChild:  0,
                    ExtraBed:   false,
                    RateKey:    rateKey,
                },
            },
        },
        Currency:    "USD",
        Language:    "en",
        AvailFlag:   true,
        DetailLevel: "Full",
    }

    recheckResp, mgErr := client.RecheckHotel(ctx, recheckReq)
    if mgErr != nil {
        log.Fatal(mgErr)
    }

    fmt.Printf("Recheck successful for hotel: %s\n", recheckResp.Hotels.Hotel.Name)

    // Step 4: Book the hotel
    bookReq := mg.BookHotelRequest{
        Login:           config.Credentials,
        SessionID:       searchResp.SessionID,
        AgencyBookingID: "TIXIA-BOOKING-12345",
        Nationality:     "US",
        CheckIn:         "2024-12-01",
        CheckOut:        "2024-12-05",
        HotelCode:       selectedHotel.Code,
        RoomDetails: mg.RoomDetailsRequest{
            Code:                   selectedRoom.Code,
            MealPlan:               selectedRoom.MealPlan,
            CancellationPolicyType: selectedRoom.CancellationPolicyType,
            PackageRate:            selectedRoom.PackageRate,
        },
        Rooms: struct {
            Room []mg.RoomBookingRequest `json:"Room"`
        }{
            Room: []mg.RoomBookingRequest{
                {
                    RoomNo:     1,
                    NoOfAdults: 2,
                    NoOfChild:  0,
                    ExtraBed:   false,
                    RateKey:    rateKey,
                    PaxDetails: struct {
                        Pax []mg.Pax `json:"Pax"`
                    }{
                        Pax: []mg.Pax{
                            {
                                Salutation: "Mr",
                                FirstName:  "John",
                                LastName:   "Doe",
                                Age:        30,
                            },
                            {
                                Salutation: "Mrs",
                                FirstName:  "Jane",
                                LastName:   "Doe",
                                Age:        28,
                            },
                        },
                    },
                },
            },
        },
        Currency:    "USD",
        Language:    "en",
        AvailFlag:   true,
        OnHold:      false,
        DetailLevel: "Full",
    }

    bookResp, mgErr := client.BookHotel(ctx, bookReq)
    if mgErr != nil {
        log.Fatal(mgErr)
    }

    fmt.Printf("Booking successful!\n")
    fmt.Printf("MG Booking ID: %s\n", bookResp.BookingDetail.MgBookingID)
    fmt.Printf("Agency Booking ID: %s\n", bookResp.BookingDetail.AgencyBookingID)

    // Step 5: Retrieve booking details
    detailsReq := mg.GetRSVNDetailsRequest{
        Login:           config.Credentials,
        MGBookingID:     bookResp.BookingDetail.MgBookingID,
        AgencyBookingID: bookResp.BookingDetail.AgencyBookingID,
    }

    detailsResp, mgErr := client.GetRSVNDetails(ctx, detailsReq)
    if mgErr != nil {
        log.Fatal(mgErr)
    }

    fmt.Printf("Booking Status: %s\n", detailsResp.BookingDetail.Status)
    fmt.Printf("Check-in: %s\n", detailsResp.BookingDetail.CheckIn)
    fmt.Printf("Check-out: %s\n", detailsResp.BookingDetail.CheckOut)
}
```

## License

This package is part of the TixiaOTA project.
