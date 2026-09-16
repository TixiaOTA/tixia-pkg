package xendit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// XenditErrorCode represents error codes returned by the Xendit API.
type XenditErrorCode string

const (
	// InvalidRequest is a client-side validation error, not returned by Xendit.
	InvalidRequest XenditErrorCode = "INVALID_REQUEST"

	// InternalServerError is used for network or unexpected client-side failures.
	InternalServerError XenditErrorCode = "INTERNAL_SERVER_ERROR"

	// HTTP 400 Bad Request error codes.
	InvalidValueError            XenditErrorCode = "INVALID_VALUE_ERROR"
	APIValidationError           XenditErrorCode = "API_VALIDATION_ERROR"
	CardExpired                  XenditErrorCode = "CARD_EXPIRED"
	InvalidPaymentDetails        XenditErrorCode = "INVALID_PAYMENT_DETAILS"
	PaymentRequestRateLimited    XenditErrorCode = "PAYMENT_REQUEST_RATE_LIMITED"
	InvalidToken                 XenditErrorCode = "INVALID_TOKEN"
	InvalidAmount                XenditErrorCode = "INVALID_AMOUNT"
	InvalidPaymentChannel        XenditErrorCode = "INVALID_PAYMENT_CHANNEL"
	InvalidPrimaryPaymentChannel XenditErrorCode = "INVALID_PRIMARY_PAYMENT_CHANNEL"
	InvalidMetadata              XenditErrorCode = "INVALID_METADATA"
	InvalidItemMetadata          XenditErrorCode = "INVALID_ITEM_METADATA"
	InvalidCurrency              XenditErrorCode = "INVALID_CURRENCY"
	InvalidCountry               XenditErrorCode = "INVALID_COUNTRY"
	InvalidCustomerID            XenditErrorCode = "INVALID_CUSTOMER_ID"
	InvalidURL                   XenditErrorCode = "INVALID_URL"
	MissingCustomer              XenditErrorCode = "MISSING_CUSTOMER"
	DisallowedCaptureMethod      XenditErrorCode = "DISALLOWED_CAPTURE_METHOD"
	InvalidExpiryDate            XenditErrorCode = "INVALID_EXPIRY_DATE"
	ItemsNotSupported            XenditErrorCode = "ITEMS_NOT_SUPPORTED"
	SessionTypeNotSupported      XenditErrorCode = "SESSION_TYPE_NOT_SUPPORTED"
	ModeNotSupported             XenditErrorCode = "MODE_NOT_SUPPORTED"
	IneligibleTransactionStatus  XenditErrorCode = "INELIGIBLE_TRANSACTION_STATUS"
	CaptureAmountExceeded        XenditErrorCode = "CAPTURE_AMOUNT_EXCEEDED"
	InactivePaymentRequest       XenditErrorCode = "INACTIVE_PAYMENT_REQUEST"
	TemporarilyUnavailable       XenditErrorCode = "TEMPORARILY_UNAVAILABLE"

	// HTTP 401 Unauthorized error codes.
	InvalidAPIKey XenditErrorCode = "INVALID_API_KEY"

	// HTTP 403 Forbidden error codes.
	InvalidMerchantSettings XenditErrorCode = "INVALID_MERCHANT_SETTINGS"

	// HTTP 404 Not Found error codes.
	NotFound     XenditErrorCode = "NOT_FOUND"
	DataNotFound XenditErrorCode = "DATA_NOT_FOUND"

	// HTTP 409 Conflict error codes.
	DuplicateError XenditErrorCode = "DUPLICATE_ERROR"

	// HTTP 422 Unprocessable Entity error codes.
	InvalidSessionStatus XenditErrorCode = "INVALID_SESSION_STATUS"

	// HTTP 500 Internal Server Error codes.
	ServerError XenditErrorCode = "SERVER_ERROR"

	// HTTP 503 Service Unavailable error codes.
	ChannelUnavailable XenditErrorCode = "CHANNEL_UNAVAILABLE"
	IssuerUnavailable  XenditErrorCode = "ISSUER_UNAVAILABLE"
)

// XenditError represents an error response from the Xendit API.
type XenditError struct {
	ErrorCode  XenditErrorCode `json:"error_code"`
	Message    string          `json:"message"`
	Errors     json.RawMessage `json:"errors,omitempty"`
	HTTPStatus int             `json:"-"`
}

func (e *XenditError) Error() string {
	if e.HTTPStatus > 0 {
		return fmt.Sprintf("Xendit error (HTTP %d): %s - %s", e.HTTPStatus, e.ErrorCode, e.Message)
	}

	return fmt.Sprintf("Xendit error: %s - %s", e.ErrorCode, e.Message)
}

// IsBadRequest reports whether err has HTTP status 400.
func IsBadRequest(err error) bool {
	return hasHTTPStatus(err, http.StatusBadRequest)
}

// IsUnauthorized reports whether err has HTTP status 401.
func IsUnauthorized(err error) bool {
	return hasHTTPStatus(err, http.StatusUnauthorized)
}

// IsForbidden reports whether err has HTTP status 403.
func IsForbidden(err error) bool {
	return hasHTTPStatus(err, http.StatusForbidden)
}

// IsNotFound reports whether err has HTTP status 404.
func IsNotFound(err error) bool {
	return hasHTTPStatus(err, http.StatusNotFound)
}

// IsConflict reports whether err has HTTP status 409.
func IsConflict(err error) bool {
	return hasHTTPStatus(err, http.StatusConflict)
}

// IsUnprocessableEntity reports whether err has HTTP status 422.
func IsUnprocessableEntity(err error) bool {
	return hasHTTPStatus(err, http.StatusUnprocessableEntity)
}

// IsInternalServerError reports whether err has HTTP status 500.
func IsInternalServerError(err error) bool {
	return hasHTTPStatus(err, http.StatusInternalServerError)
}

// IsServiceUnavailable reports whether err has HTTP status 503.
func IsServiceUnavailable(err error) bool {
	return hasHTTPStatus(err, http.StatusServiceUnavailable)
}

func hasHTTPStatus(err error, status int) bool {
	var xenditErr *XenditError
	return errors.As(err, &xenditErr) && xenditErr.HTTPStatus == status
}

// APIError represents transport-level failures such as network errors or timeouts.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (code %d): %s", e.Code, e.Message)
}
