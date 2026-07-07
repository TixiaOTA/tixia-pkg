# Xendit Provider

A Go client library for the [Xendit Payments API](https://docs.xendit.co/apidocs/create-session). This package provides typed access to Xendit payment session creation and status retrieval for hosted checkout and Component SDK flows.

## Overview

The Xendit package enables you to:

- Create one-time payment sessions (`session_type=PAY`)
- Create direct channel payment requests (`type=PAY`, `currency=IDR`, `country=ID`)
- Retrieve payment session status and payment references via `GetSession`
- Cancel active payment sessions via `CancelSession`
- Retrieve payment status and details via `GetPayment`
- Cancel authorized payments via `CancelPayment`
- Capture authorized payments via `CapturePayment`
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

### Get Payment Session Status

Retrieves the current status of a payment session and related payment references (for example after checkout or for polling).

```go
session, xerr := client.GetSession(context.Background(), xendit.GetSessionRequest{
    PaymentSessionID: "ps-661f87c614802d6c402cd82d",
})
if xerr != nil {
    log.Fatalf("GetSession failed: %v", xerr)
}

switch session.Status {
case string(xendit.SessionStatusCompleted):
    log.Printf("Payment completed: %s", *session.PaymentID)
case string(xendit.SessionStatusActive):
    log.Printf("Session still active: %s", session.PaymentSessionID)
case string(xendit.SessionStatusExpired), string(xendit.SessionStatusCanceled):
    log.Printf("Session ended with status: %s", session.Status)
}
```

### Cancel Payment Session

Cancels an active payment session and invalidates its hosted checkout URL (for example when a customer abandons checkout).

```go
session, xerr := client.CancelSession(context.Background(), xendit.CancelSessionRequest{
    PaymentSessionID: "ps-661f87c614802d6c402cd82d",
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.InvalidSessionStatus:
        log.Printf("Session cannot be canceled: %s", xerr.Message)
    default:
        log.Fatalf("CancelSession failed: %v", xerr)
    }
    return
}

log.Printf("Session canceled: %s", session.PaymentSessionID)
```

### Create Payment Request (IDR channel payment)

Creates a direct channel payment request for Indonesia (`country=ID`, `currency=IDR`). Use when you already know the payment channel (for example `OVO`, `DANA`, `QRIS`) and want to collect via redirect or channel-specific flow.

```go
captureMethod := xendit.CaptureMethodAutomatic

paymentRequest, xerr := client.CreatePaymentRequest(context.Background(), xendit.CreatePaymentRequestRequest{
    ReferenceID:   "order_12345",
    Type:          xendit.PaymentTypePay,
    Country:       "ID",
    Currency:      "IDR",
    RequestAmount: xendit.NewDecimalAmountFromFloat(150000),
    ChannelCode:   xendit.ChannelCodeOVO,
    CaptureMethod: &captureMethod,
    ChannelProperties: xendit.NewOVOChannelProperties(xendit.OVOChannelProperties{
        AccountMobileNumber: "+628123456789",
    }),
    Description: strPtr("Hotel booking"),
    Metadata: map[string]string{
        "order_id": "12345",
    },
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.InvalidPaymentChannel:
        log.Printf("Unsupported channel: %s", xerr.Message)
    default:
        log.Fatalf("CreatePaymentRequest failed: %v", xerr)
    }
    return
}

switch paymentRequest.Status {
case string(xendit.PaymentRequestStatusRequiresAction):
    for _, action := range paymentRequest.Actions {
        if action.Type == xendit.PaymentRequestActionTypeRedirectCustomer {
            log.Printf("Redirect customer to: %s", action.Value)
        }
    }
case string(xendit.PaymentRequestStatusSucceeded):
    log.Printf("Payment request succeeded: %s", paymentRequest.PaymentRequestID)
default:
    log.Printf("Payment request status: %s", paymentRequest.Status)
}
```

**Validation notes:** `CreatePaymentRequestRequest` enforces `country=ID` and `currency=IDR` only. `channel_code` and typed `channel_properties` are required. The `channel_properties` struct must match the selected channel. For example, OVO requires `OVOChannelProperties`, DANA requires `SingleRedirectChannelProperties`, and GOPAY requires `TripleRedirectChannelProperties`.

**DANA example:**

```go
paymentRequest, xerr := client.CreatePaymentRequest(context.Background(), xendit.CreatePaymentRequestRequest{
    ReferenceID:   "order_12345",
    Type:          xendit.PaymentTypePay,
    Country:       "ID",
    Currency:      "IDR",
    RequestAmount: xendit.NewDecimalAmountFromFloat(150000),
    ChannelCode:   xendit.ChannelCodeDANA,
    ChannelProperties: xendit.NewSingleRedirectChannelProperties(xendit.SingleRedirectChannelProperties{
        SuccessReturnURL: "https://example.com/success",
    }),
})
```

### Get Payment Request Status

Retrieves the status and details of a payment request. Use to poll after `CreatePaymentRequest`, or to fetch `latest_payment_id` and `payment_token_id`.

```go
paymentRequest, xerr := client.GetPaymentRequest(context.Background(), xendit.GetPaymentRequestStatusRequest{
    PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.DataNotFound:
        log.Printf("Payment request not found: %s", xerr.Message)
    default:
        log.Fatalf("GetPaymentRequest failed: %v", xerr)
    }
    return
}

log.Printf("Payment request status: %s", paymentRequest.Status)
if paymentRequest.LatestPaymentID != nil {
    log.Printf("Latest payment id: %s", *paymentRequest.LatestPaymentID)
}
```

### Cancel Payment Request

Cancels a payment request before the end user completes payment.

```go
paymentRequest, xerr := client.CancelPaymentRequest(context.Background(), xendit.CancelPaymentRequestStatusRequest{
    PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.IneligibleTransactionStatus:
        log.Printf("Payment request cannot be canceled: %s", xerr.Message)
    default:
        log.Fatalf("CancelPaymentRequest failed: %v", xerr)
    }
    return
}

log.Printf("Payment request canceled: %s", paymentRequest.PaymentRequestID)
```

### Simulate Payment Request (Test Mode)

Simulates payment completion for a payment request. Test mode only — use a test API key.

```go
simulation, xerr := client.SimulatePaymentRequest(context.Background(), xendit.SimulatePaymentRequestRequest{
    PaymentRequestID: "pr-1102feb0-bb79-47ae-9d1e-e69394d3949c",
})
if xerr != nil {
    log.Fatalf("SimulatePaymentRequest failed: %v", xerr)
}

log.Printf("Simulation status: %s", simulation.Status)
log.Printf("Simulation message: %s", simulation.Message)
```

### Get Payment Status

Retrieves the status and details of a payment. Use after `GetSession` returns a `payment_id`, or to fetch card authorization data.

```go
session, xerr := client.GetSession(context.Background(), xendit.GetSessionRequest{
    PaymentSessionID: "ps-661f87c614802d6c402cd82d",
})
if xerr != nil {
    log.Fatalf("GetSession failed: %v", xerr)
}

if session.PaymentID == nil {
    log.Fatal("payment not yet created for session")
}

payment, xerr := client.GetPayment(context.Background(), xendit.GetPaymentRequest{
    PaymentID: *session.PaymentID,
})
if xerr != nil {
    log.Fatalf("GetPayment failed: %v", xerr)
}

switch payment.Status {
case string(xendit.PaymentStatusSucceeded):
    log.Printf("Payment succeeded: %s", payment.PaymentID)
case string(xendit.PaymentStatusAuthorized):
    log.Printf("Payment authorized, awaiting capture: %s", payment.PaymentID)
case string(xendit.PaymentStatusFailed), string(xendit.PaymentStatusExpired), string(xendit.PaymentStatusCanceled):
    log.Printf("Payment ended with status: %s", payment.Status)
}

if payment.PaymentDetails != nil && payment.PaymentDetails.AuthorizationData != nil {
    auth := payment.PaymentDetails.AuthorizationData
    if auth.AuthorizationCode != nil {
        log.Printf("Authorization code: %s", *auth.AuthorizationCode)
    }
}
```

### Cancel Payment

Cancels an authorized payment before capture. Only payments in `AUTHORIZED` status can be canceled.

```go
payment, xerr := client.CancelPayment(context.Background(), xendit.CancelPaymentRequest{
    PaymentID: "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.IneligibleTransactionStatus:
        log.Printf("Payment cannot be canceled: %s", xerr.Message)
    default:
        log.Fatalf("CancelPayment failed: %v", xerr)
    }
    return
}

log.Printf("Payment canceled: %s", payment.PaymentID)
```

### Capture Payment

Captures an authorized payment to complete collection. Use with `capture_method=MANUAL` flows.

```go
payment, xerr := client.CapturePayment(context.Background(), xendit.CapturePaymentRequest{
    PaymentID:     "py-1402feb0-bb79-47ae-9d1e-e69394d3949c",
    CaptureAmount: xendit.NewDecimalAmountFromFloat(1999.01),
})
if xerr != nil {
    switch xerr.ErrorCode {
    case xendit.CaptureAmountExceeded:
        log.Printf("Capture amount too high: %s", xerr.Message)
    case xendit.IneligibleTransactionStatus:
        log.Printf("Payment is not capturable: %s", xerr.Message)
    default:
        log.Fatalf("CapturePayment failed: %v", xerr)
    }
    return
}

log.Printf("Payment captured: %s", payment.PaymentID)
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

#### GetSession

Retrieves the status of a payment session and related payment references.

```go
func (c *Client) GetSession(ctx context.Context, req GetSessionRequest) (*GetSessionResponse, *XenditError)
```

**Request Parameters:**

- `PaymentSessionID` (string, required): Xendit payment session id (27 characters, for example `ps-661f87c614802d6c402cd82d`)

**Response:**

- `PaymentSessionID` (string): Xendit payment session id
- `Status` (string): Session status — `ACTIVE`, `COMPLETED`, `EXPIRED`, or `CANCELED`
- `PaymentID` (*string): Payment id when the session is completed
- `PaymentRequestID` (*string): Payment request id when available
- `PaymentLinkURL` (*string): Hosted checkout URL when `Mode` is `PAYMENT_LINK`
- `ReferenceID`, `Amount`, `Currency`, `Country`, `Mode`: Echoed session details
- Additional fields from Xendit session response (see `models.go`)

**Timeout:** 15 seconds

**Retries:** 3 retries with exponential backoff (configurable via `Config.RetryConfig`)

**Special Notes:** Use after `CreateSession` to poll session status or reconcile payment state. Prefer webhooks (`payment_session.completed`, `payment_session.expired`) for production flows.

#### CancelSession

Cancels an active payment session. Only sessions with status `ACTIVE` can be canceled.

```go
func (c *Client) CancelSession(ctx context.Context, req CancelSessionRequest) (*CancelSessionResponse, *XenditError)
```

**Request Parameters:**

- `PaymentSessionID` (string, required): Xendit payment session id (27 characters, for example `ps-661f87c614802d6c402cd82d`)

**Response:**

- `PaymentSessionID` (string): Xendit payment session id
- `Status` (string): `CANCELED` on success
- `PaymentLinkURL` (*string): Previously issued hosted checkout URL (no longer valid)
- `ReferenceID`, `Amount`, `Currency`, `Country`, `Mode`: Echoed session details
- Additional fields from Xendit session response (see `models.go`)

**Timeout:** 15 seconds

**Retries:** 3 retries with exponential backoff (configurable via `Config.RetryConfig`)

**Special Notes:** Canceled sessions cannot be reactivated. Use when invalidating checkout links for abandoned orders or explicit merchant-side cancellation.

#### GetPayment

Retrieves the status and details of a payment.

```go
func (c *Client) GetPayment(ctx context.Context, req GetPaymentRequest) (*GetPaymentResponse, *XenditError)
```

**Request Parameters:**

- `PaymentID` (string, required): Xendit payment id (39 characters, for example `py-1402feb0-bb79-47ae-9d1e-e69394d3949c`)

**Response:**

- `PaymentID` (string): Xendit payment id
- `Status` (string): Payment status — `AUTHORIZED`, `CANCELED`, `SUCCEEDED`, `FAILED`, `EXPIRED`, or `PENDING`
- `RequestAmount` (DecimalAmount): Intended payment amount (accepts string or number from API)
- `PaymentRequestID` (*string): Associated payment request id
- `ChannelCode` (string): Payment channel used (for example `CARDS`, `OVO`)
- `PaymentDetails` (*PaymentDetails): Provider-specific details including `authorization_data` for card payments
- `Captures` ([]PaymentCapture): Capture records when applicable
- `ReferenceID`, `Currency`, `Country`, `Type`: Echoed payment details
- Additional fields from Xendit payment response (see `payment_models.go`)

**Timeout:** 15 seconds

**Retries:** 3 retries with exponential backoff (configurable via `Config.RetryConfig`)

**Special Notes:** Use after `GetSession` when you need payment-level status or card authorization data. `GetSession` returns `payment_id`; pass it to `GetPayment` for full payment details.

#### CancelPayment

Cancels an authorized payment. Only applicable for payments in `AUTHORIZED` status.

```go
func (c *Client) CancelPayment(ctx context.Context, req CancelPaymentRequest) (*CancelPaymentResponse, *XenditError)
```

**Request Parameters:**

- `PaymentID` (string, required): Xendit payment id (39 characters)

**Response:** Same shape as `GetPaymentResponse` with `Status` set to `CANCELED`.

**Timeout:** 15 seconds

**Special Notes:** Use to void an authorization before capture. Returns `INELIGIBLE_TRANSACTION_STATUS` if the payment is not in `AUTHORIZED` status.

#### CapturePayment

Captures an authorized payment to complete collection.

```go
func (c *Client) CapturePayment(ctx context.Context, req CapturePaymentRequest) (*CapturePaymentResponse, *XenditError)
```

**Request Parameters:**

- `PaymentID` (string, required): Xendit payment id (39 characters)
- `CaptureAmount` (DecimalAmount, required): Amount to capture; must be greater than zero and not exceed the authorized amount

**Response:** Same shape as `GetPaymentResponse` with `Status` set to `SUCCEEDED` and `captures` populated.

**Timeout:** 15 seconds

**Special Notes:** Use for manual capture flows (`capture_method=MANUAL`). Returns `CAPTURE_AMOUNT_EXCEEDED` when the capture amount exceeds the authorized amount.

#### CreatePaymentRequest

Creates a direct channel payment request for Indonesia. Scoped to `type=PAY`, `country=ID`, and `currency=IDR`.

```go
func (c *Client) CreatePaymentRequest(ctx context.Context, req CreatePaymentRequestRequest) (*CreatePaymentRequestResponse, *XenditError)
```

**Request Parameters:**

- `ReferenceID` (string, required): Merchant reference for the payment request
- `Type` (PaymentType, required): Must be `PAY`
- `Country` (string, required): Must be `ID`
- `Currency` (string, required): Must be `IDR`
- `RequestAmount` (DecimalAmount, required): Payment amount, must be greater than zero
- `ChannelCode` (string, required): Payment channel code constant (for example `xendit.ChannelCodeOVO`, `xendit.ChannelCodeDANA`, `xendit.ChannelCodeQRIS`)
- `ChannelProperties` (PaymentRequestChannelProperties, required): Typed channel payload that must match `ChannelCode`
- `CaptureMethod` (*CaptureMethod, optional): `AUTOMATIC` or `MANUAL`
- `CustomerID` (*string, optional): Existing Xendit customer id
- `Customer` (*CustomerDetails, optional): Inline customer details
- `Description` (*string, optional): Payment description
- `Metadata` (map[string]string, optional): Merchant metadata

**Typed `channel_properties` by channel:**

| Channel group | Struct constructor | Channels |
|---------------|-------------------|----------|
| `OVOChannelProperties` | `NewOVOChannelProperties` | `OVO` |
| `SingleRedirectChannelProperties` | `NewSingleRedirectChannelProperties` | `DANA`, `LINKAJA`, `SHOPEEPAY` |
| `DualRedirectChannelProperties` | `NewDualRedirectChannelProperties` | `ASTRAPAY`, `NEXCASH` |
| `TripleRedirectChannelProperties` | `NewTripleRedirectChannelProperties` | `GOPAY` |
| `QRISChannelProperties` | `NewQRISChannelProperties` | `QRIS` |
| `CardsChannelProperties` | `NewCardsChannelProperties` | `CARDS` |
| `VirtualAccountChannelProperties` | `NewVirtualAccountChannelProperties` | `BCA_VIRTUAL_ACCOUNT`, `BNI_VIRTUAL_ACCOUNT`, `BRI_VIRTUAL_ACCOUNT`, and other VA channels |
| `RetailOutletChannelProperties` | `NewRetailOutletChannelProperties` | `ALFAMART`, `INDOMARET` |
| `JeniusPayChannelProperties` | `NewJeniusPayChannelProperties` | `JENIUSPAY` |

**Response:**

- `PaymentRequestID` (string): Xendit payment request id
- `Status` (string): Payment request status — `REQUIRES_ACTION`, `PENDING`, `SUCCEEDED`, `FAILED`, `CANCELED`, or `EXPIRED`
- `Actions` ([]PaymentRequestAction): Customer actions required to complete payment (for example redirect URL)
- `ChannelCode` (string): Payment channel used
- `ChannelProperties` (PaymentRequestChannelProperties): Echoed channel properties; inbound JSON decodes to `map[string]any` without channel validation
- `ReferenceID`, `RequestAmount`, `Currency`, `Country`: Echoed payment request details
- Additional fields from Xendit payment request response (see `payment_request_models.go`)

**Timeout:** 15 seconds

**Retries:** 3 retries with exponential backoff (configurable via `Config.RetryConfig`)

**Special Notes:** For redirect channels, inspect `actions` when `status` is `REQUIRES_ACTION` and send the customer to the redirect URL. Payment completion is delivered asynchronously via Xendit webhooks. Webhook handling is outside this package.

#### GetPaymentRequest

Retrieves the status and details of a payment request.

```go
func (c *Client) GetPaymentRequest(ctx context.Context, req GetPaymentRequestStatusRequest) (*CreatePaymentRequestResponse, *XenditError)
```

**Request Parameters:**

- `PaymentRequestID` (string, required): Xendit payment request id (39 characters, for example `pr-1102feb0-bb79-47ae-9d1e-e69394d3949c`)

**Response:** Same shape as `CreatePaymentRequestResponse`, including optional `latest_payment_id`, `payment_token_id`, `failure_code`, `items`, and `shipping_information` when returned by Xendit.

**Timeout:** 15 seconds

**Special Notes:** Use `latest_payment_id` with `GetPayment` to fetch card authorization data. Use `payment_token_id` for subsequent pay-with-token flows.

#### CancelPaymentRequest

Cancels a payment request and prevents the end user from completing payment.

```go
func (c *Client) CancelPaymentRequest(ctx context.Context, req CancelPaymentRequestStatusRequest) (*CreatePaymentRequestResponse, *XenditError)
```

**Request Parameters:**

- `PaymentRequestID` (string, required): Xendit payment request id (39 characters)

**Response:** Same shape as `CreatePaymentRequestResponse` with `Status` set to `CANCELED`.

**Timeout:** 15 seconds

**Special Notes:** Returns `INELIGIBLE_TRANSACTION_STATUS` when the payment request is not in a cancelable status.

#### SimulatePaymentRequest

Simulates payment completion for a payment request in test mode.

```go
func (c *Client) SimulatePaymentRequest(ctx context.Context, req SimulatePaymentRequestRequest) (*SimulatePaymentRequestResponse, *XenditError)
```

**Request Parameters:**

- `PaymentRequestID` (string, required): Xendit payment request id (39 characters)
- `Amount` (*DecimalAmount, optional): Amount to simulate

**Response:**

- `Status` (string): Simulation status — always `PENDING`
- `Message` (string): Processing message; final result is delivered via webhook

**Timeout:** 15 seconds

**Special Notes:** Test mode only. Use a test API key. Payment completion is delivered asynchronously via Xendit webhooks.

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

`CreateSession`, `GetSession`, `CancelSession`, `CreatePaymentRequest`, `GetPaymentRequest`, `CancelPaymentRequest`, `SimulatePaymentRequest`, `GetPayment`, `CancelPayment`, and `CapturePayment` return `*XenditError` for validation failures and Xendit API errors.

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
| `NOT_FOUND` | Payment session does not exist or does not belong to your business |
| `DATA_NOT_FOUND` | Payment id does not exist or does not belong to your business |
| `INELIGIBLE_TRANSACTION_STATUS` | Payment is not in the required status for cancel or capture |
| `CAPTURE_AMOUNT_EXCEEDED` | Capture amount exceeds the authorized amount |
| `INVALID_MERCHANT_SETTINGS` | Merchant credentials error with the payment provider |
| `INVALID_PAYMENT_CHANNEL` | Unsupported or invalid payment channel |
| `CHANNEL_UNAVAILABLE` | Payment channel is temporarily unavailable |
| `INVALID_SESSION_STATUS` | Session is not `ACTIVE` and cannot be canceled |
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
- `GetSessionRequest` / `GetSessionResponse`
- `CancelSessionRequest` / `CancelSessionResponse`
- `GetPaymentRequest` / `GetPaymentResponse`
- `CancelPaymentRequest` / `CancelPaymentResponse`
- `CapturePaymentRequest` / `CapturePaymentResponse`
- `CreatePaymentRequestRequest` / `CreatePaymentRequestResponse`
- `GetPaymentRequestStatusRequest`
- `CancelPaymentRequestStatusRequest`
- `SimulatePaymentRequestRequest` / `SimulatePaymentRequestResponse`
- `PaymentRequestItem`, `PaymentRequestShippingInformation`
- `PaymentRequestChannelProperties`, `OVOChannelProperties`, `SingleRedirectChannelProperties`, `DualRedirectChannelProperties`, `TripleRedirectChannelProperties`, `QRISChannelProperties`, `CardsChannelProperties`, `VirtualAccountChannelProperties`, `RetailOutletChannelProperties`, `JeniusPayChannelProperties`
- `PaymentRequestAction`, `PaymentRequestStatus`
- `PaymentCapture`, `PaymentDetails`, `PaymentAuthorizationData`, `DecimalAmount`
- `CustomerDetails`, `IndividualDetail`
- `SessionItem`, `ComponentsConfiguration`, `SessionSubscription`
- `SessionType`, `SessionMode`, `CaptureMethod`, `SessionStatus`
- `PaymentType`, `PaymentStatus`

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
