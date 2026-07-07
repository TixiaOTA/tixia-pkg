package xendit

import (
	"encoding/json"
	"fmt"
)

// XenditErrorCode represents error codes returned by the Xendit API.
type XenditErrorCode string

const (
	// InvalidRequest is a client-side validation error, not returned by Xendit.
	InvalidRequest XenditErrorCode = "INVALID_REQUEST"

	// InternalServerError is used for network or unexpected client-side failures.
	InternalServerError XenditErrorCode = "INTERNAL_SERVER_ERROR"

	APIValidationError           XenditErrorCode = "API_VALIDATION_ERROR"
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
	InvalidAPIKey                XenditErrorCode = "INVALID_API_KEY"
	NotFound                     XenditErrorCode = "NOT_FOUND"
	InvalidSessionStatus         XenditErrorCode = "INVALID_SESSION_STATUS"
	DataNotFound                 XenditErrorCode = "DATA_NOT_FOUND"
	IneligibleTransactionStatus  XenditErrorCode = "INELIGIBLE_TRANSACTION_STATUS"
	InvalidMerchantSettings      XenditErrorCode = "INVALID_MERCHANT_SETTINGS"
	ChannelUnavailable           XenditErrorCode = "CHANNEL_UNAVAILABLE"
	CaptureAmountExceeded        XenditErrorCode = "CAPTURE_AMOUNT_EXCEEDED"
	InvalidValueError            XenditErrorCode = "INVALID_VALUE_ERROR"
	InactivePaymentRequest       XenditErrorCode = "INACTIVE_PAYMENT_REQUEST"
	TemporarilyUnavailable       XenditErrorCode = "TEMPORARILY_UNAVAILABLE"
	IssuerUnavailable            XenditErrorCode = "ISSUER_UNAVAILABLE"
	DuplicateError               XenditErrorCode = "DUPLICATE_ERROR"
	ServerError                  XenditErrorCode = "SERVER_ERROR"
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

// APIError represents transport-level failures such as network errors or timeouts.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (code %d): %s", e.Code, e.Message)
}
