# Vervotech Package

A Go client library for interacting with the Vervotech Hotel Mapping API. This package provides both low-level API client methods and high-level usecase methods for synchronizing hotel mappings and content.

## Overview

The Vervotech package enables you to:
- Retrieve hotel mappings (new, updated, deleted)
- Fetch all Vervotech hotel IDs
- Get curated hotel content
- Get provider-specific hotel content
- Perform initial and incremental synchronization operations

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/vervotech
```

## Configuration

### Client Configuration

Create a client instance with the required configuration:

```go
import (
    "time"
    "github.com/TixiaOTA/tixia-pkg/vervotech"
    "github.com/TixiaOTA/tixia-pkg/http/client"
)

config := vervotech.Config{
    BaseURL:    "https://hotelmapping.vervotech.com/api", // Optional, defaults to this value
    APIKey:     "your-api-key",                           // Required
    AccountID:  "your-account-id",                        // Required
    APIVersion: "3.0",                                    // Optional, defaults to "3.0"
    SyncID:     "optional-sync-id",                       // Optional
    RetryConfig: &client.RetryConfig{                     // Optional
        MaxRetries:        3,
        InitialBackoff:    1 * time.Second,
        MaxBackoff:        5 * time.Second,
        BackoffMultiplier: 2.0,
    },
}

client, err := vervotech.New(config)
if err != nil {
    log.Fatal(err)
}
```

### Configuration Fields

- **BaseURL** (string, optional): Base URL for the API. Defaults to `"https://hotelmapping.vervotech.com/api"`
- **APIKey** (string, required): API key for authentication
- **AccountID** (string, required): Account ID used in request headers
- **APIVersion** (string, optional): API version used in URL path. Defaults to `"3.0"`
- **SyncID** (string, optional): Sync ID used in some request headers
- **RetryConfig** (*client.RetryConfig, optional): Retry configuration for HTTP requests. Defaults to 3 retries with exponential backoff

## Usage

### Basic Client Usage

#### Get New Mappings

```go
ctx := context.Background()
limit := 100
providerFamily := "Booking.com"
lastUpdateDateTime := time.Now().Add(-24 * time.Hour)

req := vervotech.GetNewMappingsRequest{
    Limit:              &limit,
    ProviderFamily:     providerFamily,
    LastUpdateDateTime: &lastUpdateDateTime,
    ResumeKey:          "", // Empty for first page
}

response, err := client.GetNewMappings(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d new mappings\n", len(response.Mappings))
fmt.Printf("Resume key: %s\n", response.ResumeKey)
```

#### Get Updated Mappings

```go
req := vervotech.GetUpdatedMappingsRequest{
    Limit:              &limit,
    ProviderFamily:     providerFamily,
    LastUpdateDateTime: &lastUpdateDateTime,
    ResumeKey:          "",
}

response, err := client.GetUpdatedMappings(ctx, req)
if err != nil {
    log.Fatal(err)
}
```

#### Get Deleted Mappings

```go
req := vervotech.GetDeletedMappingsRequest{
    Limit:              &limit,
    ProviderFamily:     providerFamily,
    LastUpdateDateTime: &lastUpdateDateTime,
    ResumeKey:          "",
}

response, err := client.GetDeletedMappings(ctx, req)
if err != nil {
    log.Fatal(err)
}
```

#### Get All Vervotech IDs

```go
req := vervotech.GetAllVervotechIdsRequest{
    LastUpdateDateTime: &lastUpdateDateTime,
    ResumeKey:          "",
}

response, err := client.GetAllVervotechIds(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d Vervotech IDs\n", len(response.VervotechIds))
```

#### Get Curated Content

```go
vervotechIds := []int64{12345, 67890, 11111}

req := vervotech.GetCuratedContentByVervotechIdsRequest{
    VervotechIds: vervotechIds,
}

response, err := client.GetCuratedContentByVervotechIds(ctx, req)
if err != nil {
    log.Fatal(err)
}

for _, hotel := range response.CuratedHotels {
    fmt.Printf("Hotel: %s (ID: %d)\n", hotel.Name, hotel.VervotechID)
}
```

#### Get Provider Content

```go
req := vervotech.GetProviderContentByVervotechIdsRequest{
    NumberOfProviderContentCopies: 1,
    ProviderPreferences:           []string{"Booking.com"},
    VervotechIds:                  vervotechIds,
}

response, err := client.GetProviderContentByVervotechIds(ctx, req)
if err != nil {
    log.Fatal(err)
}

for _, group := range response.Hotels {
    fmt.Printf("Vervotech ID: %d\n", group.VervotechID)
    for _, hotel := range group.ProviderHotels {
        fmt.Printf("  Provider Hotel: %s\n", hotel.Name)
    }
}
```

### Usecase Usage

The usecase provides high-level methods that combine multiple API calls for common synchronization workflows.

#### Initialize Sync

Performs an initial synchronization, fetching:
- All Vervotech IDs
- New mappings
- Curated content for the Vervotech IDs

```go
usecase := vervotech.NewUsecase(client)

ctx := context.Background()
lastUpdate := time.Now().Add(-7 * 24 * time.Hour)
limit := 100

req := vervotech.InitializeSyncRequest{
    ProviderFamily:     stringPtr("Booking.com"),
    LastUpdateDateTime: &lastUpdate,
    Limit:              &limit,
    // Resume keys are optional for pagination
    AllVervotechIdsResumeKey: nil,
    NewMappingsResumeKey:      nil,
}

response, err := usecase.InitializeSync(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Vervotech IDs: %d\n", len(response.VervotechIds.VervotechIds))
fmt.Printf("New Mappings: %d\n", len(response.NewMappings.Mappings))
fmt.Printf("Curated Hotels: %d\n", len(response.CuratedContent.CuratedHotels))
```

#### Incremental Sync

Performs an incremental synchronization, fetching:
- Deleted mappings
- New mappings
- Updated mappings
- Curated content
- Provider content

```go
req := vervotech.IncrementalSyncRequest{
    ProviderFamily:     stringPtr("Booking.com"),
    LastUpdateDateTime: &lastUpdate,
    Limit:              &limit,
    // All resume keys are optional
    DeletedMappingsResumeKey: nil,
    NewMappingsResumeKey:      nil,
    UpdatedMappingsResumeKey:   nil,
    AllVervotechIdsResumeKey:  nil,
}

response, err := usecase.IncrementalSync(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deleted Mappings: %d\n", len(response.DeletedMappings.Mappings))
fmt.Printf("New Mappings: %d\n", len(response.NewMappings.Mappings))
fmt.Printf("Updated Mappings: %d\n", len(response.UpdatedMappings.Mappings))
fmt.Printf("Curated Hotels: %d\n", len(response.CuratedContent.CuratedHotels))
fmt.Printf("Provider Hotel Groups: %d\n", len(response.ProviderContent.Hotels))
```

### Helper Function

```go
func stringPtr(s string) *string {
    return &s
}
```

## API Reference

### Client Methods

#### GetNewMappings

Retrieves new hotel mappings from the Vervotech API.

```go
func (c *Client) GetNewMappings(ctx context.Context, req GetNewMappingsRequest) (*GetNewMappingsResponse, error)
```

**Request Parameters:**
- `Limit` (*int): Maximum number of results to return
- `ProviderFamily` (string): Provider family name filter
- `LastUpdateDateTime` (*time.Time): Filter by last update date/time
- `ResumeKey` (string): Resume key for pagination

**Response:**
- `ResumeKey` (string): Resume key for next page
- `Mappings` ([]Mapping): List of new mappings
- `Status` (bool): Request status
- `Message` (string): Status message
- `StatusCode` (int): HTTP status code

#### GetUpdatedMappings

Retrieves updated hotel mappings from the Vervotech API.

```go
func (c *Client) GetUpdatedMappings(ctx context.Context, req GetUpdatedMappingsRequest) (*GetUpdatedMappingsResponse, error)
```

**Request Parameters:** Same as `GetNewMappings`

#### GetDeletedMappings

Retrieves deleted hotel mappings from the Vervotech API.

```go
func (c *Client) GetDeletedMappings(ctx context.Context, req GetDeletedMappingsRequest) (*GetDeletedMappingsResponse, error)
```

**Request Parameters:** Same as `GetNewMappings`

#### GetAllVervotechIds

Retrieves all Vervotech hotel IDs from the API.

```go
func (c *Client) GetAllVervotechIds(ctx context.Context, req GetAllVervotechIdsRequest) (*GetAllVervotechIdsResponse, error)
```

**Request Parameters:**
- `LastUpdateDateTime` (*time.Time): Filter by last update date/time
- `ResumeKey` (string): Resume key for pagination

**Response:**
- `ResumeKey` (string): Resume key for next page
- `VervotechIds` ([]int64): List of Vervotech IDs
- `Status` (bool): Request status
- `Message` (string): Status message
- `StatusCode` (int): HTTP status code

#### GetCuratedContentByVervotechIds

Retrieves curated hotel content for the given Vervotech IDs.

```go
func (c *Client) GetCuratedContentByVervotechIds(ctx context.Context, req GetCuratedContentByVervotechIdsRequest) (*GetCuratedContentByVervotechIdsResponse, error)
```

**Request Parameters:**
- `VervotechIds` ([]int64): List of Vervotech IDs

**Response:**
- `CuratedHotels` ([]CuratedHotel): List of curated hotel content
- `Status` (bool): Request status
- `Message` (string): Status message
- `StatusCode` (int): HTTP status code

#### GetProviderContentByVervotechIds

Retrieves provider-specific hotel content for the given Vervotech IDs.

```go
func (c *Client) GetProviderContentByVervotechIds(ctx context.Context, req GetProviderContentByVervotechIdsRequest) (*GetProviderContentByVervotechIdsResponse, error)
```

**Request Parameters:**
- `NumberOfProviderContentCopies` (int): Number of provider content copies to retrieve
- `ProviderPreferences` ([]string): List of provider preferences
- `VervotechIds` ([]int64): List of Vervotech IDs

**Response:**
- `Hotels` ([]ProviderHotelGroup): List of provider hotel groups
- `Status` (bool): Request status
- `Message` (string): Status message
- `StatusCode` (int): HTTP status code

### Usecase Methods

#### InitializeSync

Performs an initial synchronization by fetching Vervotech IDs, new mappings, and curated content.

```go
func (u *Usecase) InitializeSync(ctx context.Context, req InitializeSyncRequest) (*InitializeSyncResponse, error)
```

**Request Parameters:**
- `AllVervotechIdsResumeKey` (*string): Resume key for pagination of all Vervotech IDs (optional, default: "")
- `NewMappingsResumeKey` (*string): Resume key for pagination of new mappings (optional, default: "")
- `ProviderFamily` (*string): Provider family name filter (optional, default: "")
- `LastUpdateDateTime` (*time.Time): Last update date/time filter (optional, default: nil)
- `Limit` (*int): Maximum number of results to return (optional, default: 100)

**Response:**
- `VervotechIds` (GetAllVervotechIdsResponse): All Vervotech IDs response
- `NewMappings` (GetNewMappingsResponse): New mappings response
- `CuratedContent` (GetCuratedContentByVervotechIdsResponse): Curated content response

#### IncrementalSync

Performs an incremental synchronization by fetching deleted, new, and updated mappings, along with curated and provider content.

```go
func (u *Usecase) IncrementalSync(ctx context.Context, req IncrementalSyncRequest) (*IncrementalSyncResponse, error)
```

**Request Parameters:**
- `DeletedMappingsResumeKey` (*string): Resume key for pagination of deleted mappings (optional, default: "")
- `NewMappingsResumeKey` (*string): Resume key for pagination of new mappings (optional, default: "")
- `UpdatedMappingsResumeKey` (*string): Resume key for pagination of updated mappings (optional, default: "")
- `AllVervotechIdsResumeKey` (*string): Resume key for pagination of all Vervotech IDs (optional, default: "")
- `ProviderFamily` (*string): Provider family name filter (optional, default: "")
- `LastUpdateDateTime` (*time.Time): Last update date/time filter (optional, default: nil)
- `Limit` (*int): Maximum number of results to return (optional, default: 100)

**Response:**
- `DeletedMappings` (GetDeletedMappingsResponse): Deleted mappings response
- `NewMappings` (GetNewMappingsResponse): New mappings response
- `UpdatedMappings` (GetUpdatedMappingsResponse): Updated mappings response
- `CuratedContent` (GetCuratedContentByVervotechIdsResponse): Curated content response
- `ProviderContent` (GetProviderContentByVervotechIdsResponse): Provider content response

### Client Utility Methods

#### SetSyncID

Sets or updates the SyncID for subsequent requests.

```go
func (c *Client) SetSyncID(syncID string)
```

## Data Models

### Mapping

Represents a single hotel mapping.

```go
type Mapping struct {
    VervotechId          int64    `json:"VervotechId"`
    ProviderHotelId      string   `json:"ProviderHotelId"`
    ProviderName         string   `json:"ProviderName"`
    ChannelIds           []string `json:"ChannelIds"`
    ProviderLocationCode *string  `json:"ProviderLocationCode"`
    Type                 string   `json:"Type"`
}
```

### CuratedHotel

Represents curated hotel content with comprehensive hotel information.

```go
type CuratedHotel struct {
    VervotechID            int
    Name                   string
    Rating                 string
    AddressLine1           string
    AddressLine2           string
    CityName               string
    CityCode               string
    StateName              string
    StateCode              string
    CountryName            string
    CountryCode            string
    PostalCode             string
    Latitude               float64
    Longitude              float64
    Website                string
    Fax                    string
    PropertyType           string
    ChainName              string
    BrandName              string
    CityLocationID         string
    MasterCityLocationName string
    GooglePlusCode         string
    LocationIds            string
    Phones                 []string
    Emails                 []string
    Facilities             []Facility
    Airports               []Location
    TrainStations          []Location
    PointOfInterests       []Location
    ConnectedLocations     []Location
    Neighbourhoods         []Location
    Stadiums               []Location
    LocationMappings       []LocationMapping
}
```

### ProviderHotel

Represents provider-specific hotel content with detailed room, facility, and policy information.

```go
type ProviderHotel struct {
    ProviderHotelID       string
    VervotechID           int
    LanguageCode          string
    Name                  string
    ProviderFamily        string
    Rating                float64
    BrandName             string
    BrandCode             string
    ChainCode             string
    ChainName             string
    MasterPropertyType    string
    PropertyType          string
    Facilities            []ProviderFacility
    HeroImage             string
    ProviderHeroImageHref string
    Rooms                 []ProviderRoom
    MasterChainName       string
    PopularityScore       int
    ContentScore          int
    LicenseNumber         interface{}
    GeoCode               ProviderGeoCode
    ProviderGeocodes      ProviderGeoCode
    Checkout              ProviderCheckout
    Descriptions          []ProviderDescription
    Policies              []ProviderPolicy
    Checkin               ProviderCheckin
    Attributes            []ProviderAttribute
    Statistics            []ProviderStatistic
    Reviews               []ProviderReview
    PointOfInterests      []PointOfInterest
    Fees                  []ProviderFee
    Contact               ProviderContact
    Images                []ProviderImage
}
```

See `models.go` for complete model definitions.

## Error Handling

The package returns `APIError` for API-related errors:

```go
type APIError struct {
    Code    int
    Message string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (status %d): %s", e.Code, e.Message)
}
```

**Example:**

```go
response, err := client.GetNewMappings(ctx, req)
if err != nil {
    if apiErr, ok := err.(*vervotech.APIError); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## Pagination

Most endpoints support pagination using `ResumeKey`. To paginate through results:

```go
resumeKey := ""
for {
    req := vervotech.GetNewMappingsRequest{
        Limit:     &limit,
        ResumeKey: resumeKey,
    }
    
    response, err := client.GetNewMappings(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process mappings
    for _, mapping := range response.Mappings {
        // Process each mapping
    }
    
    // Check if there are more pages
    if response.ResumeKey == "" {
        break
    }
    
    resumeKey = response.ResumeKey
}
```

## Examples

### Complete Initial Sync Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/TixiaOTA/tixia-pkg/vervotech"
)

func main() {
    // Create client
    config := vervotech.Config{
        APIKey:    "your-api-key",
        AccountID: "your-account-id",
    }
    
    client, err := vervotech.New(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create usecase
    usecase := vervotech.NewUsecase(client)
    
    // Perform initial sync
    ctx := context.Background()
    lastUpdate := time.Now().Add(-7 * 24 * time.Hour)
    limit := 100
    
    req := vervotech.InitializeSyncRequest{
        ProviderFamily:     stringPtr("Booking.com"),
        LastUpdateDateTime: &lastUpdate,
        Limit:              &limit,
    }
    
    response, err := usecase.InitializeSync(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Initial Sync Complete:\n")
    fmt.Printf("  Vervotech IDs: %d\n", len(response.VervotechIds.VervotechIds))
    fmt.Printf("  New Mappings: %d\n", len(response.NewMappings.Mappings))
    fmt.Printf("  Curated Hotels: %d\n", len(response.CuratedContent.CuratedHotels))
}

func stringPtr(s string) *string {
    return &s
}
```

## License

This package is part of the TixiaOTA project.

