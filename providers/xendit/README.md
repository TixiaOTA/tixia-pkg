# Xendit Provider

A Go client library for the [Xendit Payments API](https://docs.xendit.co/apidocs/create-session). This package provides typed access to Xendit payment session creation for hosted checkout and Component SDK flows.

## Overview

The Xendit package enables you to:

- Create one-time payment sessions (`session_type=PAY`)
- Use hosted checkout (`mode=PAYMENT_LINK`) and redirect via `payment_link_url`
- Use Xendit Component SDK (`mode=COMPONENTS`) with `components_sdk_key`
- Validate requests before calling the Xendit API
- Handle Xendit API errors with typed error codes
- Trace and log requests using OpenTelemetry and structured logging

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/providers/xendit
```

## Configuration

### Client Configuration

```go
package main

import (
    "log"
    "os"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/xendit"
)

func main() {
    logInstance := logger.New(logger.Config{
        Level:      "info",
        JSONOutput: false,
    })

    client, err := xendit.New(xendit.Config{
        SecretAPIKey: os.Getenv("XENDIT_SECRET_KEY"),
        ForUserID:    "", // optional XenPlatform sub-account
        SplitRuleID:  "", // optional XenPlatform split rule
    }, logInstance)
    if err != nil {
        log.Fatal(err)
    }

    _ = client
}
```

### Configuration Fields

- **BaseURL** (string, optional): Xendit API base URL. Default: `https://api.xendit.co`
- **SecretAPIKey** (string, required): Secret API key used for Basic Auth (username with empty password)
- **ForUserID** (string, optional): XenPlatform sub-account user id sent as `for-user-id` header
- **SplitRuleID** (string, optional): XenPlatform split rule id sent as `with-split-rule` header
- **RetryConfig** (*client.RetryConfig, optional): HTTP retry settings. Default: 3 retries, 5s initial backoff, 15s max backoff, 2x multiplier

## Usage

### Create a Payment Session

Creates a one-time payment session and returns a hosted checkout URL.

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/xendit"
)

func strPtr(value string) *string {
    return &value
}

func main() {
    logInstance := logger.New(logger.Config{
        Level:      "info",
        JSONOutput: false,
    })

    client, err := xendit.New(xendit.Config{
        SecretAPIKey: os.Getenv("XENDIT_SECRET_KEY"),
    }, logInstance)
    if err != nil {
        log.Fatal(err)
    }

    resp, xerr := client.CreateSession(context.Background(), xendit.CreateSessionRequest{
        ReferenceID: "order_12345",
        SessionType: xendit.SessionTypePay,
        Mode:        xendit.SessionModePaymentLink,
        Amount:      150000,
        Currency:    "IDR",
        Country:     "ID",
        Customer: &xendit.CustomerDetails{
            ReferenceID: "cust_001",
            Type:        xendit.CustomerTypeIndividual,
            Email:       strPtr("customer@example.com"),
            IndividualDetail: xendit.IndividualDetail{
                GivenNames: "John",
                Surname:    strPtr("Doe"),
            },
        },
        SuccessReturnURL: strPtr("https://yourcompany.com/success"),
        CancelReturnURL:  strPtr("https://yourcompany.com/cancel"),
    })
    if xerr != nil {
        log.Fatalf("CreateSession failed: %v", xerr)
    }

    // Redirect the customer to resp.PaymentLinkURL
    log.Printf("Payment session created: %s", resp.PaymentSessionID)
    log.Printf("Checkout URL: %s", *resp.PaymentLinkURL)
}
```

### Create a Components Session

Creates a one-time payment session for embedding the Xendit Component SDK on your site.

```go
resp, xerr := client.CreateSession(context.Background(), xendit.CreateSessionRequest{
    ReferenceID: "order_12345",
    SessionType: xendit.SessionTypePay,
    Mode:        xendit.SessionModeComponents,
    Amount:      10000,
    Currency:    "PHP",
    Country:     "PH",
    Customer: &xendit.CustomerDetails{
        ReferenceID: "cust_001",
        Type:        xendit.CustomerTypeIndividual,
        Email:       strPtr("customer@example.com"),
        IndividualDetail: xendit.IndividualDetail{
            GivenNames: "John",
            Surname:    strPtr("Doe"),
        },
    },
    ComponentsConfiguration: &xendit.ComponentsConfiguration{
        Origins: []string{"https://yourcompany.com"},
    },
})
if xerr != nil {
    log.Fatalf("CreateSession failed: %v", xerr)
}

// Initialize Xendit Component SDK on frontend with resp.ComponentsSDKKey
log.Printf("Components SDK key: %s", *resp.ComponentsSDKKey)
```

## API Reference

### Client Methods

#### CreateSession

Creates a Xendit payment session for one-time payments using hosted checkout or Component SDK.

```go
func (c *Client) CreateSession(ctx context.Context, req CreateSessionRequest) (*CreateSessionResponse, *XenditError)
```

**Request Parameters:**

- `ReferenceID` (string, required): Merchant reference for the payment session
- `Currency` (string, required): ISO 4217 currency code (`IDR`, `PHP`, `VND`, `THB`, `SGD`, `MYR`, `USD`)
- `Amount` (float64, required): Payment amount, must be greater than zero
- `Country` (string, required): ISO 3166-1 alpha-2 country code (`ID`, `PH`, `VN`, `TH`, `SG`, `MY`)
- `SessionType` (SessionType, required): Must be `PAY`
- `Mode` (SessionMode, required): `PAYMENT_LINK` or `COMPONENTS`
- `CustomerID` (*string, optional): Existing Xendit customer id
- `Customer` (*CustomerDetails, optional): Inline customer details
- `CaptureMethod` (*CaptureMethod, optional): `AUTOMATIC` or `MANUAL`
- `Description` (*string, optional): Checkout page description
- `SuccessReturnURL` (*string, optional): HTTPS redirect URL after success (typically used with `PAYMENT_LINK`)
- `CancelReturnURL` (*string, optional): HTTPS redirect URL after cancel (typically used with `PAYMENT_LINK`)
- `Locale` (*string, optional): Hosted checkout locale
- `Metadata` (map[string]string, optional): Merchant metadata
- `Items` ([]SessionItem, optional): Line items for the session
- `AllowedPaymentChannels` ([]string, optional): Restrict available payment channels
- `ExpiresAt` (*time.Time, optional): Session expiry time
- `ComponentsConfiguration` (*ComponentsConfiguration, optional): Required when `Mode` is `COMPONENTS`; must include HTTPS `origins` for CORS validation

**Response:**

- `PaymentSessionID` (string): Xendit payment session id
- `Status` (string): Session status (for example `ACTIVE`)
- `PaymentLinkURL` (*string): Hosted checkout URL when `Mode` is `PAYMENT_LINK`
- `ComponentsSDKKey` (*string): SDK key for initializing Xendit Component SDK when `Mode` is `COMPONENTS`
- `ComponentsConfiguration` (*ComponentsConfiguration): Echoed components configuration from the response
- `ReferenceID`, `Amount`, `Currency`, `Country`: Echoed session details
- Additional fields from Xendit session response (see `models.go`)

**Timeout:** 15 seconds

**Retries:** 3 retries with exponential backoff (configurable via `Config.RetryConfig`)

**Special Notes:** For `PAYMENT_LINK`, redirect the customer to `payment_link_url`. For `COMPONENTS`, use `components_sdk_key` with the Xendit Component SDK on your frontend. Payment completion is delivered asynchronously via Xendit webhooks (`payment_session.completed`). Webhook handling is outside this package.

## Error Handling

### Error Types

```go
type XenditError struct {
    ErrorCode  XenditErrorCode
    Message    string
    Errors     json.RawMessage
    HTTPStatus int
}

type APIError struct {
    Code    int
    Message string
}
```

`CreateSession` returns `*XenditError` for validation failures and Xendit API errors.

### Common Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Client-side validation failure |
| `API_VALIDATION_ERROR` | Xendit request validation failure |
| `INVALID_AMOUNT` | Invalid payment amount |
| `INVALID_CURRENCY` | Unsupported currency |
| `INVALID_COUNTRY` | Unsupported country |
| `INVALID_URL` | Return URL is not valid HTTPS |
| `MISSING_CUSTOMER` | Customer object or customer id required |
| `INVALID_API_KEY` | Invalid or unauthorized API key |
| `DUPLICATE_ERROR` | Duplicate customer reference id |
| `SERVER_ERROR` | Xendit internal server error |
| `INTERNAL_SERVER_ERROR` | Network or unexpected client-side failure |

### Error Handling Example

```go
resp, xerr := client.CreateSession(ctx, req)
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.InvalidRequest:
        // Fix request payload before retrying
    case xendit.InvalidAPIKey:
        // Check API key configuration
    case xendit.APIValidationError:
        // Inspect xerr.Message and xerr.Errors
    default:
        // Log and handle other Xendit errors
    }
    return
}
```

## Data Models

Key request and response types are defined in `models.go`:

- `CreateSessionRequest` / `CreateSessionResponse`
- `CustomerDetails`, `IndividualDetail`
- `SessionItem`, `ComponentsConfiguration`
- `SessionType`, `SessionMode`, `CaptureMethod`

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/providers/xendit"
)

func strPtr(value string) *string {
    return &value
}

func main() {
    logInstance := logger.New(logger.Config{
        Level:      "info",
        JSONOutput: true,
    })

    client, err := xendit.New(xendit.Config{
        SecretAPIKey: os.Getenv("XENDIT_SECRET_KEY"),
    }, logInstance)
    if err != nil {
        log.Fatal(err)
    }

    resp, xerr := client.CreateSession(context.Background(), xendit.CreateSessionRequest{
        ReferenceID: "order_12345",
        SessionType: xendit.SessionTypePay,
        Mode:        xendit.SessionModePaymentLink,
        Amount:      150000,
        Currency:    "IDR",
        Country:     "ID",
        Description: strPtr("Payment for order #12345"),
        Customer: &xendit.CustomerDetails{
            ReferenceID: "cust_001",
            Type:        xendit.CustomerTypeIndividual,
            Email:       strPtr("customer@example.com"),
            IndividualDetail: xendit.IndividualDetail{
                GivenNames: "John",
            },
        },
        SuccessReturnURL: strPtr("https://yourcompany.com/order/complete"),
        CancelReturnURL:  strPtr("https://yourcompany.com/order/cancel"),
    })
    if xerr != nil {
        log.Fatalf("payment session creation failed: %v", xerr)
    }

    fmt.Printf("Redirect customer to: %s\n", *resp.PaymentLinkURL)
}
```

## License

This package is part of the TixiaOTA project.
